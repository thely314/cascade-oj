package service

import (
	"context"

	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"

	"google.golang.org/genproto/googleapis/api/httpbody"
)

func (s *AdminService) ListLogFiles(ctx context.Context, req *pb.ListLogFilesRequest) (*pb.ListLogFilesReply, error) {
	files, total, err := s.logUseCase.ListLogFiles(ctx, req.PageSize, req.Offset)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, len(files))
	for i, f := range files {
		filenames[i] = f.Name
	}

	return &pb.ListLogFilesReply{
		Filenames: filenames,
		Total:     int32(total),
	}, nil
}

func (s *AdminService) QueryLogContent(ctx context.Context, req *pb.QueryLogContentRequest) (*pb.QueryLogContentReply, error) {
	query := &biz.LogQuery{
		Filename: req.Filename,
		Level:    req.Level,
		PageSize: req.PageSize,
		Offset:   req.Offset,
	}
	if req.TimeStart != nil {
		t := req.TimeStart.AsTime()
		query.StartTime = &t
	}
	if req.TimeEnd != nil {
		t := req.TimeEnd.AsTime()
		query.EndTime = &t
	}

	lines, total, err := s.logUseCase.QueryLogContent(ctx, query)
	if err != nil {
		return nil, err
	}

	return &pb.QueryLogContentReply{
		Lines: lines,
		Total: int32(total),
	}, nil
}

func (s *AdminService) DownloadLogs(ctx context.Context, req *pb.DownloadLogsRequest) (*httpbody.HttpBody, error) {
	buf, err := s.logUseCase.DownloadLogs(ctx, req.Filenames)
	if err != nil {
		return nil, err
	}

	return &httpbody.HttpBody{
		ContentType: "application/zip",
		Data:        buf.Bytes(),
	}, nil
}
