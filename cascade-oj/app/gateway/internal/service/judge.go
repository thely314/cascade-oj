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
		UUID:       uuid.New().String(),
		UserID:     ctx.Value("userInfo").(auth.Claims).UserID,
		ProblemID:  req.ProblemId,
		Code:       req.Code,
		Language:   req.Language,
		Input:      req.Input,
		IsCompiled: false,
		Stdout:     "",
		Stderr:     "",
		TimeCost:   0,
		MemoryCost: 0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SelfTestReply{
		Uuid: selfTestID,
	}, nil
}

func (gatewayService *GatewayService) PostSubmission(ctx context.Context, req *pb.PostSubmissionRequest) (*pb.PostSubmissionReply, error) {
	submissionID, err := gatewayService.judgeUsecase.CreateSubmission(ctx, &biz.Submission{
		UUID:        uuid.New().String(),
		UserID:      ctx.Value("userInfo").(auth.Claims).UserID,
		ProblemID:   req.ProblemId,
		Code:        req.Code,
		Language:    req.Language,
		Status:      "Pending",
		CreateTime:  time.Now(),
		Score:       0,
		TimeCost:    0,
		MemoryCost:  0,
		CaseVersion: 0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PostSubmissionReply{
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
	submission, err := gatewayService.judgeUsecase.GetSingleSubmission(ctx, req.SubmissionUuid)
	if err != nil {
		return nil, err
	}
	var casesResults []*pb.CaseMetadata
	for _, v := range submission.CasesResults {
		casesResults = append(casesResults, &pb.CaseMetadata{
			Score:      v.Score,
			Status:     v.Status,
			TimeCost:   int32(v.TimeCost),
			MemoryCost: int32(v.MemoryCost),
		})
	}
	return &pb.GetSingleSubmissionReply{
		Metadata: &pb.SubmissionMetadata{
			SubmissionUuid: submission.Submission.UUID,
			ProblemId:      submission.Submission.ProblemID,
			UserId:         submission.Submission.UserID,
			Status:         submission.Submission.Status,
			SubmitTime:     timestamppb.New(submission.Submission.CreateTime),
			Score:          int32(submission.Submission.Score),
		},
		Code:       submission.Submission.Code,
		Language:   submission.Submission.Language,
		TimeCost:   int32(submission.Submission.TimeCost),
		MemoryCost: int32(submission.Submission.MemoryCost),
		CaseResults: &pb.GetSingleSubmissionReply_CaseResults{
			Cases: casesResults,
		},
	}, nil
}

func (gatewayService *GatewayService) GetSelfTestResult(ctx context.Context, req *pb.GetSelfTestResultRequest) (*pb.GetSelfTestResultReply, error) {
	selfTest, err := gatewayService.judgeUsecase.GetSelfTest(ctx, req.SelftestUuid)
	if err != nil {
		return nil, err
	}
	return &pb.GetSelfTestResultReply{
		IsCompiled: selfTest.IsCompiled,
		Stdout:     selfTest.Stdout,
		Stderr:     selfTest.Stderr,
		TimeCost:   int32(selfTest.TimeCost),
		MemoryCost: int32(selfTest.MemoryCost),
	}, nil
}
