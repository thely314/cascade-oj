package biz

import (
	"bytes"
	"context"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type LogFile struct {
	Name    string
	Size    int64
	ModTime time.Time
}

type LogQuery struct {
	Filename  string
	Level     string
	StartTime *time.Time
	EndTime   *time.Time
	PageSize  int32
	Offset    int32
}

type LogRepo interface {
	ListFiles(ctx context.Context) ([]os.FileInfo, error)
	GetFileContent(ctx context.Context, filename string) ([]byte, error)
	CreateZip(ctx context.Context, filenames []string) (*bytes.Buffer, error)
}

type LogUseCase struct {
	repo LogRepo
	log  *log.Helper
}

func NewLogUseCase(repo LogRepo, logger log.Logger) *LogUseCase {
	return &LogUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *LogUseCase) ListLogFiles(ctx context.Context, pageSize, offset int32) ([]*LogFile, int, error) {
	files, err := uc.repo.ListFiles(ctx)
	if err != nil {
		return nil, 0, err
	}

	logFiles := make([]*LogFile, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			logFiles = append(logFiles, &LogFile{
				Name:    f.Name(),
				Size:    f.Size(),
				ModTime: f.ModTime(),
			})
		}
	}

	total := len(logFiles)
	start := int(offset)
	end := start + int(pageSize)
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return logFiles[start:end], total, nil
}

func (uc *LogUseCase) QueryLogContent(ctx context.Context, query *LogQuery) ([]string, int, error) {
	content, err := uc.repo.GetFileContent(ctx, query.Filename)
	if err != nil {
		return nil, 0, err
	}

	lines := strings.Split(string(content), "\n")
	var filteredLines []string

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Filtering logic
		matchLevel := query.Level == "" || strings.Contains(line, query.Level)

		var matchTime bool
		if query.StartTime == nil && query.EndTime == nil {
			matchTime = true
		} else {
			matchTime = false
			// The format is ts=YYYY-MM-DDTHH:MM:SSZ
			if strings.Contains(line, "ts=") {
				tsStrStart := strings.Index(line, "ts=")
				if tsStrStart != -1 {
					tsStr := line[tsStrStart+3:]
					// The timestamp have no quoted. It ends at Z.
					tsStrEnd := strings.Index(tsStr, "Z")
					if tsStrEnd != -1 {
						tsStr = tsStr[:tsStrEnd+1]
						logTime, err := time.Parse("2006-01-02T15:04:05Z", tsStr)
						if err == nil {
							isAfterStart := query.StartTime == nil || logTime.After(*query.StartTime) || logTime.Equal(*query.StartTime)
							isBeforeEnd := query.EndTime == nil || logTime.Before(*query.EndTime) || logTime.Equal(*query.EndTime)
							matchTime = isAfterStart && isBeforeEnd
						} else {
							uc.log.Warnf("Failed to parse log timestamp: %v, error: %v", tsStr, err)
						}
					}
				}
			}
		}

		if matchLevel && matchTime {
			filteredLines = append(filteredLines, line)
		}
	}

	total := len(filteredLines)
	start := int(query.Offset)
	end := start + int(query.PageSize)
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return filteredLines[start:end], total, nil
}

func (uc *LogUseCase) DownloadLogs(ctx context.Context, filenames []string) (*bytes.Buffer, error) {
	return uc.repo.CreateZip(ctx, filenames)
}
