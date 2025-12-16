package service

import (
	"context"
	"time"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"
	"cascade-oj/pkg/middleware/auth"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (gatewayService *GatewayService) PostSelfTest(ctx context.Context, req *pb.SelfTestRequest) (*pb.SelfTestReply, error) {
	selfTestID, err := gatewayService.judgeUsecase.CreateSelfTest(ctx, &biz.SelfTest{
		UUID:     uuid.New().String(),
		UserID:   ctx.Value("userInfo").(auth.Claims).UserID,
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
	submissionID, err := gatewayService.judgeUsecase.CreateSubmission(ctx, &biz.Submission{
		UUID:       uuid.New().String(),
		UserID:     ctx.Value("userInfo").(auth.Claims).UserID,
		ProblemID:  req.ProblemId,
		Code:       req.Code,
		Language:   req.Language,
		Status:     "Pending",
		CreateTime: time.Now(),
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

func (gatewayService *GatewayService) GetSubmissions(ctx context.Context, req *pb.GetSubmissionsRequest) (
	*pb.GetSubmissionsReply, error) {
	submissions, err := gatewayService.judgeUsecase.GetSubmissions(ctx, req.UserId, req.ProblemId)
	if err != nil {
		return nil, err
	}
	res := &pb.GetSubmissionsReply{}
	for _, v := range submissions {
		res.Submissions = append(res.Submissions, &pb.SubmissionMetadata{
			SubmissionUuid: v.SubmissionID,
			ProblemId:      v.ProblemID,
			UserId:         v.UserID,
			Status:         v.Status,
			SubmitTime:     timestamppb.New(v.SubmitTime),
			Score:          int32(v.Score),
		})
	}
	return res, nil
}

func (gatewayService *GatewayService) GetSingleSubmission(ctx context.Context, req *pb.GetSingleSubmissionRequest) (
	*pb.GetSingleSubmissionReply, error) {
	submission, err := gatewayService.judgeUsecase.GetSingleSubmission(ctx, req.SubmissionId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSingleSubmissionReply{
		Metadata: &pb.SubmissionMetadata{
			SubmissionUuid: submission.UUID,
			ProblemId:      submission.ProblemID,
			UserId:         submission.UserID,
			Status:         submission.Status,
			SubmitTime:     timestamppb.New(submission.CreateTime),
			Score:          int32(submission.Score),
		},
		Code:     submission.Code,
		Language: submission.Language,
	}, nil
}
