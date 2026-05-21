package logic

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/redis/go-redis/v9"
)

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

func (f *FileRedis) FileFromString(path string) File {
	baseName := filepath.Base(path)
	file := File{Name: baseName, Path: path}
	f.Files = append(f.Files, file)
	return file
}

func (f *FileRedis) GetFile(name string) (File, error) {
	cmd := f.redisClient.HGet(f.ctx, f.Key, name)
	if cmd.Err() != nil {
		return File{}, cmd.Err()
	}
	return File{Name: name, Path: cmd.Val()}, nil
}

func (f *FileRedis) DeleteFile() error {
	return nil
}

// SaveFiles saves the list of log files in Redis using a
// hash with the file name as the key and the file path as the value
func (f *FileRedis) SaveFiles(path string) error {
	files, err := f.GetLocalLogs(path)

	if err != nil {
		return err
	}

	for _, file := range files {
		cmd := f.redisClient.HSet(f.ctx, f.Key, file.Name, file.Path)
		if err := cmd.Err(); err != nil {
			return err
		}
	}
	return nil
}

// GetLocalLogs retrieves all the log files in the root directory
// and returns them as a slice of File structs
func (f *FileRedis) GetLocalLogs(path string) ([]File, error) {
	var files []File
	_path := strings.TrimSuffix(path, "/")
	
	if _path == "" {
		_path = "data"
	}

	fullpath, err := filepath.Abs(f.rootDir + fmt.Sprintf("/%s", _path))
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".log" {
			log.Printf("⚠️ Skipping file %s:", path)
			return nil
		}

		if !info.IsDir() {
			files = append(files, f.FileFromString(path))
		}

		return nil
	})
	return files, err
}

// ReadFile reads the content of a log file and returns it as a slice of strings
func (f *FileRedis) ReadFile(path string, serverConfig *models.ServerConfig) ([]string, error) {
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

// CacheContent reads the content of a log file and caches it in Redis
func (f *FileRedis) CacheContent(file File) error {
	logs, err := f.ReadFile(file.Path, nil)
	if err != nil {
		return err
	}

	values := make([]any, len(logs))
	for i, l := range logs {
		values[i] = l
	}

	name := fmt.Sprintf("go-xlogger:%s", file.Name)
	err = f.redisClient.RPush(f.ctx, name, values...).Err()
	if err != nil {
		return err
	}

	err = f.redisClient.Expire(f.ctx, name, 15*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewFileRedis(ctx context.Context, rootDir string, redisClient *redis.Client) *FileRedis {
	return &FileRedis{
		ctx:         ctx,
		rootDir:     rootDir,
		redisClient: redisClient,
		Files:       []File{},
		Key:         "go-xlogger:files",
	}
}
