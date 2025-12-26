package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (adminService *AdminService) GetLogs(ctx context.Context, request *pb.GetLogsRequest) (*pb.GetLogsReply, error) {
	logs, err := adminService.logUseCase.GetLogs(ctx,
		biz.LogRequestInfo{
			Page:     request.Page,
			PageSize: request.PageSize,
		})
	if err != nil {
		return nil, err
	}
	pbLogs := make([]*pb.GetLogsReply_LogEntry, 0, len(logs))
	for _, log := range logs {
		pbLogs = append(pbLogs,
			&pb.GetLogsReply_LogEntry{
				Id:        log.ID,
				Message:   log.Message,
				Timestamp: timestamppb.New(log.TimeStamp),
			},
		)
	}
	return &pb.GetLogsReply{
		Logs: pbLogs,
	}, nil
}
