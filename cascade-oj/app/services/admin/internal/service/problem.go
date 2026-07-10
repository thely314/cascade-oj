package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/pkg/middleware/auth"
	"context"
	"errors"
	"strconv"
	"strings"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
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

	return &pb.GetSingleProblemReply{
		Metadata: &pb.ProblemMetadata{
			Id:            problem.Problem.ID,
			Title:         problem.Problem.Title,
			TimeLimitMs:   problem.Problem.TimeLimitMs,
			MemoryLimitMb: problem.Problem.MemoryLimitKB / 1024,
			Status:        mapBizStatusToPbStatus(problem.Problem.Status),
			Description:   problem.Problem.Description,
		},
		Creator:     problem.CreatorUsername,
		Description: problem.Description,
		Templates:   mapBizTemplatesToPb(problem.Problem.Templates),
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
			Templates:     mapPbTemplatesToBiz(request.Templates),
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
			Templates:     mapPbTemplatesToBiz(request.Templates),
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

func mapPbTemplatesToBiz(pbTemplates []*pb.CodeTemplate) []*biz.ProblemTemplate {
	bizTemplates := make([]*biz.ProblemTemplate, 0, len(pbTemplates))
	for _, t := range pbTemplates {
		// 将语言和文件名拼接，例如 "cpp/main.cpp"
		fileName := t.Language + "/" + t.Name
		bizTemplates = append(bizTemplates, &biz.ProblemTemplate{
			Name:    fileName,
			Content: t.Code,
		})
	}
	return bizTemplates
}

func mapBizTemplatesToPb(bizTemplates []*biz.ProblemTemplate) []*pb.CodeTemplate {
	pbTemplates := make([]*pb.CodeTemplate, 0, len(bizTemplates))
	for _, t := range bizTemplates {
		parts := strings.SplitN(t.Name, "/", 2)
		language := "cpp"
		name := t.Name

		if len(parts) == 2 {
			language = parts[0]
			name = parts[1]
		}

		pbTemplates = append(pbTemplates, &pb.CodeTemplate{
			Language: language,
			Name:     name,
			Code:     t.Content,
		})
	}
	return pbTemplates
}

// UploadTestCasesRaw 处理原生的文件流上传
func (adminService *AdminService) UploadTestCasesRaw(ctx khttp.Context) error {
	req := ctx.Request()
	problemIdStr := ctx.Vars().Get("id")

	// [REVIEW FIX 4]: 修复路由参数解析错误忽略的问题 (Copilot 提示)
	// 防止 ID 无效时覆盖 cases/0/testcase 目录
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	if err != nil || problemId <= 0 {
		return ctx.Result(400, map[string]interface{}{
			"isSuccess": false,
			"message":   "invalid problem id",
		})
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		return ctx.Result(400, map[string]interface{}{
			"isSuccess": false,
			"message":   "missing file field: " + err.Error(),
		})
	}
	defer file.Close()

	err = adminService.problemUseCase.HandleTestCasesUpload(ctx, problemId, file, handler.Filename)
	if err != nil {
		return ctx.Result(500, map[string]interface{}{
			"isSuccess": false,
			"message":   err.Error(),
		})
	}

	return ctx.Result(200, map[string]interface{}{
		"isSuccess": true,
		"message":   "upload success",
	})
}
