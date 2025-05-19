package upload

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLocal_UploadFile(t *testing.T) {
	savePath := "./testdata/"
	baseUrl := "http://localhost/uploads"
	_ = os.MkdirAll(savePath, os.ModePerm)
	defer os.RemoveAll(savePath) // 清理测试数据

	localOSS, err := NewLocalOSS(savePath, baseUrl)
	if err != nil {
		t.Fatalf("failed to create local OSS: %v", err)
	}

	// 模拟文件内容
	fileContent := []byte("this is a test file")
	fileName := "test.txt"
	reader := bytes.NewReader(fileContent)

	// 测试上传文件
	category := "test_category"
	url, key, err := localOSS.UploadFile(context.Background(), category, reader, fileName)
	if err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}

	// 验证文件是否正确保存
	expectedFilePath := filepath.Join(savePath, category, key)
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Fatalf("uploaded file does not exist: %v", err)
	}

	// 验证返回的 URL 是否正确
	expectedUrl := baseUrl + "/" + category + "/" + key
	if url != expectedUrl {
		t.Errorf("expected URL %s, got %s", expectedUrl, url)
	}
}

func TestLocal_DeleteFile(t *testing.T) {
	savePath := "./testdata"
	baseUrl := "http://localhost/uploads"
	_ = os.MkdirAll(savePath, os.ModePerm)
	defer os.RemoveAll(savePath) // 清理测试数据

	localOSS, err := NewLocalOSS(savePath, baseUrl)
	if err != nil {
		t.Fatalf("failed to create local OSS: %v", err)
	}

	// 创建一个测试文件
	fileName := "test_delete.txt"
	filePath := filepath.Join(savePath, fileName)
	if err := os.WriteFile(filePath, []byte("test content"), os.ModePerm); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// 测试删除文件
	err = localOSS.DeleteFile(context.Background(), "", fileName)
	if err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}

	// 验证文件是否被删除
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("file was not deleted: %v", err)
	}
}
