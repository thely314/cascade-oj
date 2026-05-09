package data

import (
	"archive/zip"
	"bytes"
	"cascade-oj/app/services/admin/internal/biz"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

const logDir = "/data/logs"

type LogRepo struct {
	data *Data
	log  *log.Helper
}

func NewLogRepo(data *Data, logger log.Logger) biz.LogRepo {
	return &LogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *LogRepo) ListFiles(ctx context.Context) ([]os.FileInfo, error) {
	dir, err := os.Open(logDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *LogRepo) GetFileContent(ctx context.Context, filename string) ([]byte, error) {
	cleanName := filepath.Clean(filename)
	filePath := filepath.Join(logDir, cleanName)
	// Security check to prevent directory traversal
	if !strings.HasPrefix(filePath, logDir) || strings.Contains(filePath, "..") {
		return nil, os.ErrPermission
	}
	return os.ReadFile(filePath)
}

func (r *LogRepo) CreateZip(ctx context.Context, filenames []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, filename := range filenames {
		cleanName := filepath.Clean(filename)
		filePath := filepath.Join(logDir, cleanName)
		if !strings.HasPrefix(filePath, logDir) || strings.Contains(filePath, "..") {
			r.log.Warnf("Skipping file due to security policy: %s", cleanName)
			continue
		}

		// 每轮循环关闭文件，避免资源耗尽
		fatal, err := func() (bool, error) {
			fileToZip, err := os.Open(filePath)
			if err != nil {
				return false, fmt.Errorf("Failed to open file %s for zipping: %v", filename, err)
			}
			defer fileToZip.Close()

			info, err := fileToZip.Stat()
			if err != nil {
				return false, fmt.Errorf("Failed to get file info for %s: %v", filename, err)
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return false, fmt.Errorf("Failed to create zip header for %s: %v", filename, err)
			}
			header.Name = cleanName
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return true, fmt.Errorf("Failed to create zip writer for %s: %v", filename, err)
			}
			_, err = io.Copy(writer, fileToZip)
			if err != nil {
				return true, fmt.Errorf("Failed to copy file content for %s: %v", filename, err)
			}

			return false, nil
		}()
		if err != nil {
			if fatal {
				// zip 损坏性错误
				r.log.Errorf("Aborting zip creation: %v", filename, err)
				zipWriter.Close()
				return nil, err
			}
			// 跳过单个文件的错误
			r.log.Errorf("Skipping file %s: %v", filename, err)
			continue
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
