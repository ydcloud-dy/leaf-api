package oss

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/ydcloud-dy/leaf-api/config"
)

var client *oss.Client
var bucket *oss.Bucket
var useLocalStorage = false

// Init 初始化OSS客户端
func Init() error {
	// 检查 OSS 配置是否完整
	if config.AppConfig.OSS.Endpoint == "" ||
		config.AppConfig.OSS.AccessKeyID == "" ||
		config.AppConfig.OSS.AccessKeySecret == "" ||
		config.AppConfig.OSS.BucketName == "" {
		useLocalStorage = true
		// 创建本地上传目录
		os.MkdirAll("uploads", 0755)
		return fmt.Errorf("OSS configuration incomplete, using local storage")
	}

	var err error
	client, err = oss.New(
		config.AppConfig.OSS.Endpoint,
		config.AppConfig.OSS.AccessKeyID,
		config.AppConfig.OSS.AccessKeySecret,
	)
	if err != nil {
		useLocalStorage = true
		os.MkdirAll("uploads", 0755)
		return fmt.Errorf("failed to create OSS client: %w, using local storage", err)
	}

	bucket, err = client.Bucket(config.AppConfig.OSS.BucketName)
	if err != nil {
		useLocalStorage = true
		os.MkdirAll("uploads", 0755)
		return fmt.Errorf("failed to get bucket: %w, using local storage", err)
	}

	return nil
}

// UploadFile 上传文件到OSS或本地存储
func UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%s/%s%s",
		folder,
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	// 如果使用本地存储或 bucket 未初始化
	if useLocalStorage || bucket == nil {
		return uploadToLocal(src, filename)
	}

	// 尝试使用 OSS 上传，如果失败则使用本地存储
	url, err := uploadToOSS(src, filename)
	if err != nil {
		// OSS 上传失败，切换到本地存储
		useLocalStorage = true
		// 重新打开文件
		src, _ = file.Open()
		return uploadToLocal(src, filename)
	}

	return url, nil
}

// uploadToOSS 上传到 OSS（带超时）
func uploadToOSS(src multipart.File, filename string) (string, error) {
	// 创建一个带超时的通道
	done := make(chan error, 1)
	var uploadErr error

	go func() {
		uploadErr = bucket.PutObject(filename, src)
		done <- uploadErr
	}()

	// 等待上传完成或超时
	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("failed to upload file: %w", err)
		}
		url := fmt.Sprintf("%s/%s", config.AppConfig.OSS.BaseURL, filename)
		return url, nil
	case <-time.After(5 * time.Second):
		return "", fmt.Errorf("upload timeout after 5 seconds")
	}
}

// uploadToLocal 上传文件到本地存储
func uploadToLocal(src multipart.File, filename string) (string, error) {
	// 创建目标目录
	destPath := filepath.Join("uploads", filename)
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建目标文件
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// 返回完整 URL（包含服务器地址和端口）
	port := config.AppConfig.Server.Port
	url := fmt.Sprintf("http://localhost:%d/uploads/%s", port, filename)
	return url, nil
}

// DeleteFile 从OSS删除文件
func DeleteFile(objectKey string) error {
	err := bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetObjectKeyFromURL 从URL中提取对象键
func GetObjectKeyFromURL(url string) string {
	baseURL := config.AppConfig.OSS.BaseURL
	if len(url) > len(baseURL) {
		return url[len(baseURL)+1:]
	}
	return ""
}
