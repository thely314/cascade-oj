package service

import (
	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"
	"context"

	"github.com/google/uuid"
)

type GatewayService struct {
	pb.UnimplementedUserServer
	judgeUsecase *biz.JudgeUsecase
}

func NewGatewayService(judgeUsecase *biz.JudgeUsecase) *GatewayService {
	return &GatewayService{
		judgeUsecase: judgeUsecase,
	}
}

func (gatewayService *GatewayService) PostSelfTest(ctx context.Context, req *pb.SelfTestRequest) (*pb.SelfTestReply, error) {
	// TODO: Implementation of PostSelfTest
	id := uuid.New().String()
	selfTestID, err := gatewayService.judgeUsecase.CreateSelfTest(ctx, &biz.SelfTest{
		ID:       id,
		UserID:   0, // TODO: get from context
		Code:     req.Code,
		Language: req.Language,
		Input:    req.SelfCase,
		Status:   "Pending",
	})
	if err != nil {
		return nil, err
	}
	return &pb.SelfTestReply{
		Uuid: selfTestID,
	}, nil
}

func (gatewayService *GatewayService) PostSubmission(ctx context.Context, req *pb.SubmissionRequest) (*pb.SubmissionReply, error) {
	// TODO: Implementation of PostSubmission
	submissionID, err := gatewayService.judgeUsecase.CreateSubmission(ctx, &biz.Submission{
		ID:         uuid.New().String(),
		UserID:     0, // TODO: get from context
		ProblemID:  req.ProblemId,
		Code:       req.Code,
		Language:   req.Language,
		Status:     "Pending",
		CreateTime: 0,
		Score:      0,
		TimeCost:   0,
		MemoryCost: 0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SubmissionReply{
		Uuid: submissionID,
	}, nil
}
