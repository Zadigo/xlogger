package tickerapp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/redis/go-redis/v9"
)

type FileRedis struct {
	Key         string `json:"key"`
	Files       []models.File `json:"files"`
	ctx         context.Context
	rootDir     string
	redisClient *redis.Client
}

// GetFile retrieves a file from Redis by its name and returns it as a File struct
func (f *FileRedis) GetFile(name string) (models.File, error) {
	cmd := f.redisClient.HGet(f.ctx, f.Key, name)
	if cmd.Err() != nil {
		return models.File{}, cmd.Err()
	}

	value := cmd.Val()
	if value == "" {
		return models.File{}, fmt.Errorf("file not found")
	}

	var file models.File
	err := json.Unmarshal([]byte(value), &file)
	if err != nil {
		return models.File{}, err
	}

	return file, nil
}

// ReadFile reads the content of a log file and returns it as a slice of strings
func (f *FileRedis) ReadFile(path string, serverConfig *utils.ServerConfig) ([]string, error) {
	file, err := os.Open(path)

	var logs []string = make([]string, 0)
	if err != nil {
		log.Printf("❌ Could not open file: %s", err)
		return logs, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		logs = append(logs, line)
	}
	
	if err := scanner.Err(); err != nil {
		return logs, err
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

// HasCached checks if the cached data for a specific file exists in Redis
func (f *FileRedis) HasCachedData(name string) bool {
	cmd := f.redisClient.Exists(f.ctx, fmt.Sprintf("go-xlogger:%s", name))
	if cmd.Err() != nil {
		return false
	}

	return cmd.Val() > 0
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
func (f *FileRedis) SaveFiles(files []models.File) error {
	for _, file := range files {
		jsonData, err := json.Marshal(file)
		if err != nil {
			return err
		}

		cmd := f.redisClient.HSet(f.ctx, f.Key, file.Name, jsonData)
		if err := cmd.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileRedis) GetCachedFiles() ([]models.File, error) {
	cmd := f.redisClient.HGetAll(f.ctx, f.Key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var files []models.File
	for _, data := range cmd.Val() {
		var file models.File
		err := json.Unmarshal([]byte(data), &file)
		if err != nil {
			continue
		}
		files = append(files, file)
	}

	return files, nil
}

func (f *FileRedis) NumberOfFilesInFolder() int {
	return 0
}

func NewFileRedis(ctx context.Context, redisClient *redis.Client) *FileRedis {
	if ctx == nil {
		log.Fatal("❌ Context is nil")
	}

	if redisClient == nil {
		log.Fatal("❌ Redis client is nil")
	}

	rootDirValue := ctx.Value("rootDir")
	if rootDirValue == nil {
		log.Fatal("❌ rootDir is not set in context")
	}

	rootDir := rootDirValue.(string)

	return &FileRedis{
		ctx:         ctx,
		rootDir:     rootDir,
		redisClient: redisClient,
		Files:       []models.File{},
		Key:         "go-xlogger:files",
	}
}
