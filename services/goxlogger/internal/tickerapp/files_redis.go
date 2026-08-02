package tickerapp

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/redis/go-redis/v9"
)

// File represents a log file with its name and path
type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type FileRedis struct {
	Key         string `json:"key"`
	Files       []File `json:"files"`
	ctx         context.Context
	rootDir     string
	redisClient *redis.Client
}

// fileFromString creates a File struct from a given path and adds it to the Files slice
func (f *FileRedis) fileFromString(path string) File {
	baseName := filepath.Base(path)

	file := File{Name: baseName, Path: path}
	f.Files = append(f.Files, file)

	return file
}

// GetFile retrieves a file from Redis by its name and returns it as a File struct
func (f *FileRedis) GetFile(name string) (File, error) {
	cmd := f.redisClient.HGet(f.ctx, f.Key, name)
	if cmd.Err() != nil {
		return File{}, cmd.Err()
	}
	return File{Name: name, Path: cmd.Val()}, nil
}

// ReadFile reads the content of a log file and returns it as a slice of strings
func (f *FileRedis) ReadFile(path string, serverConfig *utils.ServerConfig) ([]string, error) {
	file, err := os.Open(path)

	var logs []string = make([]string, 0)
	if err != nil {
		log.Fatal("❌ Could not open file")
		return logs, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		logs = append(logs, line)
	}

	return logs, nil
}

func (f *FileRedis) DeleteFile() error {
	return nil
}

// GetLogs retrieves the cached logs for a specific file from Redis
// and returns them as a slice of LogLine structs
func (f *FileRedis) GetLogs(name string) (lines []LogLine, err error) {
	cmd := f.redisClient.LRange(f.ctx, fmt.Sprintf("go-xlogger:%s", name), 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var logs []LogLine
	for _, log := range cmd.Val() {
		line := LogLine{RawLine: log}
		// When a line cannot be parsed,
		// we skip it and continue with the next line
		_, err := line.ParseLine()

		if err != nil {
			continue
		}

		logs = append(logs, line)
	}

	return logs, nil
}

// CacheLogs caches the content of a log file in Redis using a list with the file name as the key
func (f *FileRedis) CacheLogs(fileName string, content []string) error {
	values := make([]any, len(content))
	for i, l := range content {
		values[i] = l
	}

	name := fmt.Sprintf("go-xlogger:%s", fileName)
	err := f.redisClient.RPush(f.ctx, name, values...).Err()
	if err != nil {
		return err
	}

	err = f.redisClient.Expire(f.ctx, name, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

// SaveFiles saves the list of log files in Redis using a
// hash with the file name as the key and the file path as the value
func (f *FileRedis) SaveFiles(files []File) error {
	for _, file := range files {
		cmd := f.redisClient.HSet(f.ctx, f.Key, file.Name, file.Path)
		if err := cmd.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileRedis) GetFiles() ([]File, error) {
	cmd := f.redisClient.HGetAll(f.ctx, f.Key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var files []File
	for name, path := range cmd.Val() {
		files = append(files, File{Name: name, Path: path})
	}

	return files, nil
}

// GetLocalLogs retrieves all the log files in the root directory
// and returns them as a slice of File structs
func (f *FileRedis) GetLocalLogs(path string) ([]File, error) {
	var files []File
	trimmedPath := strings.TrimSuffix(path, "/")

	if trimmedPath == "" {
		trimmedPath = "data"
	}

	fullpath, err := filepath.Abs(f.rootDir + fmt.Sprintf("/%s", trimmedPath))
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".log" {
			log.Printf("⚠️ Skipping file %s:", path)
			return nil
		}

		if !info.IsDir() {
			files = append(files, f.fileFromString(path))
		}

		return nil
	})
	return files, err
}

func NewFileRedis(ctx context.Context, redisClient *redis.Client) *FileRedis {
	rootDir := ctx.Value("rootDir").(string)

	return &FileRedis{
		ctx:         ctx,
		rootDir:     rootDir,
		redisClient: redisClient,
		Files:       []File{},
		Key:         "go-xlogger:files",
	}
}
