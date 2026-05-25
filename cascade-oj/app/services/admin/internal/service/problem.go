package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/pkg/middleware/auth"
	"context"
	"errors"
)

func mapBizStatusToPbStatus(status biz.ProblemStatus) pb.ProblemStatus {
	switch status {
	case biz.ProblemStatusUnavailable:
		return pb.ProblemStatus_PROBLEM_STATUS_UNAVAILABLE
	case biz.ProblemStatusAvailable:
		return pb.ProblemStatus_PROBLEM_STATUS_AVAILABLE
	case biz.ProblemStatusUsing:
		return pb.ProblemStatus_PROBLEM_STATUS_USING
	case biz.ProblemStatusDeleted:
		return pb.ProblemStatus_PROBLEM_STATUS_DELETED
	default:
		return pb.ProblemStatus_PROBLEM_STATUS_UNAVAILABLE
	}
}

func (adminService *AdminService) GetProblems(ctx context.Context, request *pb.GetProblemsRequest) (*pb.GetProblemsReply, error) {
	problems, err := adminService.problemUseCase.GetProblems(
		ctx,
		request.ContestId,
	)
	if err != nil {
		return nil, err
	}
	pbProblems := make([]*pb.ProblemMetadata, 0, len(problems))
	for _, problem := range problems {
		pbProblems = append(pbProblems,
			&pb.ProblemMetadata{
				Id:            problem.ID,
				Title:         problem.Title,
				TimeLimitMs:   problem.TimeLimitMs,
				MemoryLimitMb: problem.MemoryLimitKB / 1024,
				Status:        mapBizStatusToPbStatus(problem.Status),
				Description:   problem.Description,
			},
		)
	}
	return &pb.GetProblemsReply{
		Problems: pbProblems,
	}, nil
}

func (adminService *AdminService) GetSingleProblem(ctx context.Context, request *pb.GetSingleProblemRequest) (*pb.GetSingleProblemReply, error) {
	problem, err := adminService.problemUseCase.GetSingleProblem(ctx, request.ProblemId)
	if err != nil {
		return nil, err
	}

	var templateStr string
	if len(problem.Problem.Templates) > 0 {
		templateStr = problem.Problem.Templates[0].Content
	}

	return &pb.GetSingleProblemReply{
		Metadata: &pb.ProblemMetadata{
			Id:            problem.Problem.ID,
			Title:         problem.Problem.Title,
			TimeLimitMs:   problem.Problem.TimeLimitMs,
			MemoryLimitMb: problem.Problem.MemoryLimitKB / 1024,
			Status:        mapBizStatusToPbStatus(problem.Problem.Status),
			Description:   problem.Problem.Description,
		},
		Creator:      problem.CreatorUsername,
		Description:  problem.Description,
		CodeTemplate: templateStr,
	}, nil
}

func (adminService *AdminService) PostProblem(ctx context.Context, request *pb.PostProblemRequest) (*pb.PostProblemReply, error) {
	claims, ok := auth.FromContext(ctx)
	if !ok {
		return nil, errors.New("unauthorized: metadata not found in context")
	}

	problemID, err := adminService.problemUseCase.PostProblem(ctx, claims.UserID,
		biz.ProblemCreateInfo{
			Title:         request.Metadata.Title,
			TimeLimitMs:   request.Metadata.TimeLimitMs,
			MemoryLimitKB: request.Metadata.MemoryLimitMb * 1024,
			Description:   request.Description,
			Templates: []*biz.ProblemTemplate{
				{
					Name:    "main.cpp",
					Content: request.CodeTemplate,
				},
			},
		})

	if err != nil {
		return nil, err
	}
	return &pb.PostProblemReply{
		ProblemId: problemID,
	}, nil
}

func (adminService *AdminService) PutProblem(ctx context.Context, request *pb.PutProblemRequest) (*pb.PutProblemReply, error) {
	isUpdated, err := adminService.problemUseCase.PutProblem(ctx,
		biz.ProblemEditInfo{
			ID:            request.ProblemId,
			Title:         request.Metadata.Title,
			TimeLimitMs:   request.Metadata.TimeLimitMs,
			MemoryLimitKB: request.Metadata.MemoryLimitMb * 1024,
			Description:   request.Description,
			Templates: []*biz.ProblemTemplate{
				{
					Name:    "main.cpp",
					Content: request.CodeTemplate,
				},
			},
		})
	if err != nil {
		return nil, err
	}
	return &pb.PutProblemReply{
		IsUpdated: isUpdated,
	}, nil
}

func (adminService *AdminService) PublishProblem(ctx context.Context, request *pb.PublishProblemRequest) (*pb.PublishProblemReply, error) {
	isSuccess, err := adminService.problemUseCase.PublishProblem(ctx, request.ProblemId)
	if err != nil {
		return nil, err
	}
	return &pb.PublishProblemReply{
		IsSuccess: isSuccess,
	}, nil
}

func (adminService *AdminService) DeleteProblem(ctx context.Context, request *pb.DeleteProblemRequest) (*pb.DeleteProblemReply, error) {
	isDeleted, err := adminService.problemUseCase.DeleteProblem(ctx, request.ProblemId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProblemReply{
		IsDeleted: isDeleted,
	}, nil
}

func (adminService *AdminService) DisableProblem(ctx context.Context, request *pb.DisableProblemRequest) (*pb.DisableProblemReply, error) {
	isSuccess, err := adminService.problemUseCase.DisableProblem(ctx, request.ProblemId)
	if err != nil {
		return nil, err
	}
	return &pb.DisableProblemReply{
		IsSuccess: isSuccess,
	}, nil
}
