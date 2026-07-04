package service

import (
	"context"

	pb "cascade-oj/api/cascade/aiops/v1"
	"cascade-oj/app/services/aiops/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAIOpsService)

// AIOpsService implements the AIOps gRPC/HTTP API.
type AIOpsService struct {
	pb.UnimplementedAIOpsServer

	analyzer *biz.Analyzer
	reporter *biz.Reporter
	log      *log.Helper
}

// NewAIOpsService creates a new AIOpsService.
func NewAIOpsService(
	analyzer *biz.Analyzer,
	reporter *biz.Reporter,
	logger log.Logger,
) *AIOpsService {
	return &AIOpsService{
		analyzer: analyzer,
		reporter: reporter,
		log:      log.NewHelper(logger),
	}
}

// TriggerAnalysis manually triggers an alert analysis run.
func (s *AIOpsService) TriggerAnalysis(ctx context.Context, req *pb.TriggerAnalysisRequest) (*pb.TriggerAnalysisReply, error) {
	timeRange := int32(15)
	if req.TimeRangeMinutes > 0 {
		timeRange = req.TimeRangeMinutes
	}
	s.log.Infof("Manual analysis triggered: time_range=%dm, services=%v", timeRange, req.Services)

	// Detach from HTTP request context — the full pipeline (Prometheus →
	// log scan → LLM CoT → DB persist) is a long-running task that must
	// not be killed by the HTTP server timeout.
	// 避免 ctx 被 HTTP 请求超时取消，整个分析流程是一个长时间运行的任务
	bgCtx := context.WithoutCancel(ctx)

	result, err := s.analyzer.Run(bgCtx, timeRange, req.Services)
	if err != nil {
		s.log.Errorf("Analysis failed: %v", err)
		return nil, err
	}

	report, err := s.reporter.GenerateReport(bgCtx, result)
	if err != nil {
		s.log.Errorf("Report generation failed: %v", err)
		return nil, err
	}

	return &pb.TriggerAnalysisReply{
		ReportId:           report.ID,
		Title:              report.Title,
		Summary:            report.Summary,
		NoiseCount:         int32(report.NoiseCount),
		RealRiskCount:      int32(report.RealRiskCount),
		FullReportMarkdown: report.FullReportMarkdown,
		CreatedAt:          report.CreatedAt,
	}, nil
}

// GetReport retrieves a single report by ID.
func (s *AIOpsService) GetReport(ctx context.Context, req *pb.GetReportRequest) (*pb.GetReportReply, error) {
	report, err := s.reporter.GetReport(ctx, req.ReportId)
	if err != nil {
		return nil, err
	}
	return &pb.GetReportReply{
		Id:                   report.ID,
		Title:                report.Title,
		Summary:              report.Summary,
		RootCausesJson:       report.RootCausesJSON,
		SuggestedActionsJson: report.SuggestedActionsJSON,
		NoiseCount:           int32(report.NoiseCount),
		RealRiskCount:        int32(report.RealRiskCount),
		Status:               report.Status,
		FullReportMarkdown:   report.FullReportMarkdown,
		CreatedAt:            report.CreatedAt,
		AcknowledgedAt:       report.AcknowledgedAt,
	}, nil
}

// ListReports lists historical reports.
func (s *AIOpsService) ListReports(ctx context.Context, req *pb.ListReportsRequest) (*pb.ListReportsReply, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 {
		pageSize = 20
	}

	reports, total, err := s.reporter.ListReports(ctx, page, pageSize, req.StatusFilter)
	if err != nil {
		return nil, err
	}

	var summaries []*pb.ReportSummary
	for _, r := range reports {
		summaries = append(summaries, &pb.ReportSummary{
			Id:            r.ID,
			Title:         r.Title,
			NoiseCount:    int32(r.NoiseCount),
			RealRiskCount: int32(r.RealRiskCount),
			Status:        r.Status,
			CreatedAt:     r.CreatedAt,
		})
	}

	return &pb.ListReportsReply{
		Reports: summaries,
		Total:   int32(total),
	}, nil
}

// GetMetricsSummary returns current metrics snapshot.
func (s *AIOpsService) GetMetricsSummary(ctx context.Context, req *pb.GetMetricsSummaryRequest) (*pb.GetMetricsSummaryReply, error) {
	snapshot, err := s.analyzer.GetMetricsSnapshot(ctx, req.Services)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

// AckReport acknowledges a report by updating its status to "acknowledged".
func (s *AIOpsService) AckReport(ctx context.Context, req *pb.AckReportRequest) (*pb.AckReportReply, error) {
	s.log.Infof("Acknowledging report %d", req.ReportId)
	if err := s.reporter.AckReport(ctx, req.ReportId); err != nil {
		s.log.Errorf("Failed to acknowledge report %d: %v", req.ReportId, err)
		return nil, err
	}
	return &pb.AckReportReply{}, nil
}
