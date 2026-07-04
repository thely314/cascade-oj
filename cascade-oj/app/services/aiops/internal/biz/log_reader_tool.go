package biz

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// LogReaderTool reads service log files and provides filtered views.
type LogReaderTool struct {
	logDir string
	log    *log.Helper
}

// LogEntry represents a single parsed log line.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Raw       string    `json:"raw"`
}

// LogSummary is a summary of log scanning results.
type LogSummary struct {
	Service    string     `json:"service"`
	TotalLines int        `json:"total_lines"`
	ErrorCount int        `json:"error_count"`
	WarnCount  int        `json:"warn_count"`
	TopErrors  []LogEntry `json:"top_errors"`
	ScanPeriod string     `json:"scan_period"`
}

// NewLogReaderTool creates a new LogReaderTool.
func NewLogReaderTool(logger log.Logger) *LogReaderTool {
	return &LogReaderTool{
		logDir: "/data/logs",
		log:    log.NewHelper(logger),
	}
}

// logFiles returns the list of log files in the log directory.
func (l *LogReaderTool) logFiles() ([]string, error) {
	entries, err := os.ReadDir(l.logDir)
	if err != nil {
		return nil, fmt.Errorf("reading log dir %s: %w", l.logDir, err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".log") {
			files = append(files, filepath.Join(l.logDir, e.Name()))
		}
	}
	return files, nil
}

// serviceNameFromPath extracts the service name from a log file path.
func serviceNameFromPath(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, ".log")
}

// ScanRecent scans recent log entries (within timeRange minutes) for errors and warnings.
// Returns at most maxLines per file to limit token consumption.
func (l *LogReaderTool) ScanRecent(ctx context.Context, timeRangeMinutes int32, services []string, maxLines int) ([]LogSummary, error) {
	if maxLines <= 0 {
		maxLines = 500
	}

	files, err := l.logFiles()
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-time.Duration(timeRangeMinutes) * time.Minute)
	serviceFilter := make(map[string]bool)
	for _, s := range services {
		serviceFilter[s] = true
	}

	var summaries []LogSummary

	for _, file := range files {
		svcName := serviceNameFromPath(file)
		if len(services) > 0 && !serviceFilter[svcName] {
			continue
		}

		summary, err := l.scanFile(ctx, file, svcName, cutoff, maxLines)
		if err != nil {
			l.log.Warnf("scanning %s: %v", file, err)
			continue
		}
		summaries = append(summaries, *summary)
	}

	return summaries, nil
}

// scanFile reads a log file and extracts error/warning entries.
func (l *LogReaderTool) scanFile(ctx context.Context, path, svcName string, cutoff time.Time, maxLines int) (*LogSummary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	summary := &LogSummary{
		Service:    svcName,
		ScanPeriod: fmt.Sprintf("last %s", time.Since(cutoff).Round(time.Minute)),
	}

	scanner := bufio.NewScanner(f)
	// Increase buffer for long log lines
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	linesRead := 0
	for scanner.Scan() && linesRead < maxLines*10 { // read more than max to find errors
		linesRead++
		line := scanner.Text()

		// Parse log level: look for ERROR, WARN, INFO etc.
		entry := parseLogLine(line, svcName)
		if entry.Timestamp.Before(cutoff) {
			continue
		}

		summary.TotalLines++

		switch strings.ToUpper(entry.Level) {
		case "ERROR", "FATAL", "PANIC":
			summary.ErrorCount++
			if len(summary.TopErrors) < maxLines {
				summary.TopErrors = append(summary.TopErrors, entry)
			}
		case "WARN", "WARNING":
			summary.WarnCount++
		}
	}

	if err := scanner.Err(); err != nil {
		l.log.Warnf("Scanner error for %s: %v", path, err)
	}

	// Sort errors by recency
	sort.Slice(summary.TopErrors, func(i, j int) bool {
		return summary.TopErrors[i].Timestamp.After(summary.TopErrors[j].Timestamp)
	})

	return summary, nil
}

// parseLogLine attempts to extract timestamp, level, and message from a log line.
func parseLogLine(line, svcName string) LogEntry {
	entry := LogEntry{
		Service: svcName,
		Raw:     line,
		Level:   "INFO",
	}

	// Try common log formats
	// Time format is ts=YYYY-MM-DDTHH:MM:SSZ
	if strings.Contains(line, "ts=") {
		tsStrStart := strings.Index(line, "ts=")
		if tsStrStart != -1 {
			tsStr := line[tsStrStart+3:]
			// The timestamp have no quoted. It ends at Z.
			tsStrEnd := strings.Index(tsStr, "Z")
			if tsStrEnd != -1 {
				tsStr = tsStr[:tsStrEnd+1]
				ts, err := time.Parse("2006-01-02T15:04:05Z", tsStr)
				if err == nil {
					entry.Timestamp = ts
					// Try to extract level
					// parts := strings.Fields(line[30:])
					// if len(parts) > 0 {
					// 	entry.Level = strings.ToUpper(parts[0])
					// 	if len(parts) > 1 {
					// 		entry.Message = strings.Join(parts[1:], " ")
					// 	}
					// }
					return entry
				}
			}
		}
	}

	// Simple scan for level keywords
	upper := strings.ToUpper(line)
	for _, level := range []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"} {
		if idx := strings.Index(upper, level); idx >= 0 {
			entry.Level = strings.ToUpper(strings.TrimRight(level, ":"))
			if len(line) > idx+len(level) {
				entry.Message = strings.TrimSpace(line[idx+len(level):])
				entry.Message = strings.TrimLeft(entry.Message, ": ")
			}
			break
		}
	}

	entry.Timestamp = time.Now()
	entry.Message = limitString(entry.Message, 500)

	return entry
}

// limitString truncates a string to maxLen characters.
func limitString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
