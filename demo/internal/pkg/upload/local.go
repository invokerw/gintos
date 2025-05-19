package upload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type local struct {
	savePath string
	baseUrl  string
}

var _ OSS = &local{}

func NewLocalOSS(savePath, baseUrl string) (OSS, error) {
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}
	return &local{
		savePath: savePath,
		baseUrl:  baseUrl,
	}, nil
}

func isValidString(s string) bool {
	// 定义正则表达式，匹配只包含大小写字母、数字和下划线的字符串
	pattern := `^[a-zA-Z0-9_]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(s)
}

// UploadFile saves the file locally and returns the file path and URL
func (l *local) UploadFile(ctx context.Context, category string, reader io.Reader, fileName string) (string, string, error) {
	ext := filepath.Ext(fileName)
	if len(ext) > 0 && !isValidString(ext[1:]) {
		return "", "", errors.New("file extension contains invalid characters")
	}
	name := strings.TrimSuffix(fileName, ext)
	if len(name) > 64 {
		name = name[:64]
	}
	if !isValidString(name) {
		return "", "", errors.New("file name contains invalid characters " + name)
	}
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	filePath := filepath.Join(l.savePath, category, filename)
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return "", "", fmt.Errorf("failed to create %s directory: %w", filePath, err)
	}
	url := l.baseUrl + "/"
	if len(category) > 0 {
		url += category + "/"
	}
	url += filename

	out, createErr := os.Create(filePath)
	if createErr != nil {
		return "", "", errors.New("function os.Create() failed, err:" + createErr.Error())
	}
	defer func() {
		_ = out.Close()
	}()

	_, copyErr := io.Copy(out, reader)
	if copyErr != nil {
		return "", "", errors.New("function io.Copy() failed, err:" + copyErr.Error())
	}
	return url, filename, nil
}

// DeleteFile removes the file from local storage
func (l *local) DeleteFile(ctx context.Context, category string, key string) error {
	ext := filepath.Ext(key)
	if len(ext) > 0 && !isValidString(ext[1:]) {
		return errors.New("file extension contains invalid characters")
	}
	name := strings.TrimSuffix(key, ext)
	if len(name) > 64 {
		name = name[:64]
	}
	if !isValidString(name) {
		return errors.New("file name contains invalid characters")
	}
	filePath := filepath.Join(l.savePath, category, key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
