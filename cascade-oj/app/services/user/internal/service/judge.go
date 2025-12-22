package service

import (
	"context"
	"time"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/pkg/middleware/auth"
	"cascade-oj/pkg/util"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) PostSubmission(ctx context.Context, req *pb.PostSubmissionRequest) (*pb.PostSubmissionReply, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	// TODO: validate user not banned and contest/problem access
	// exist, err := s.TODO.CheckBanned(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }
	// if exist {
	// 	return nil, pb.ErrorUserBanned("user is banned")
	// }
	// _, err = s.TODO.Find(ctx, req.ContestId)
	// if err != nil {
	// 	return nil, pb.ErrorContestEnd("contest is end")
	// }
	id := uuid.NewString()
	_, err := s.judgeUsecase.CreateSubmission(ctx, &biz.Submission{
		UUID:         id,
		UserID:       userID,
		ProblemSetID: req.ContestId,
		ProblemID:    req.ProblemId,
		Code:         util.CRLF2LF(req.Code),
		Language:     req.Language,
		Status:       "",
		CreateTime:   time.Now(),
		Score:        0,
		TimeCost:     0,
		MemoryCost:   0,
		CaseVersion:  0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PostSubmissionReply{Uuid: id}, nil
}

func (s *UserService) PostSelfTest(ctx context.Context, req *pb.SelfTestRequest) (*pb.SelfTestReply, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	// TODO: Add ban check
	// exist, err := s.TODO.CheckBanned(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }
	// if exist {
	// 	return nil, pb.ErrorUserBanned("user is banned")
	// }
	id := uuid.NewString()
	_, err := s.judgeUsecase.CreateSelfTest(ctx, &biz.SelfTest{
		UUID:       id,
		UserID:     userID,
		ProblemID:  req.ProblemId,
		Code:       util.CRLF2LF(req.Code),
		Language:   req.Language,
		Input:      util.CRLF2LF(req.Input),
		IsCompiled: false,
		Stdout:     "",
		Stderr:     "",
		TimeCost:   0,
		MemoryCost: 0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SelfTestReply{Uuid: id}, nil
}

func (s *UserService) GetSingleSubmission(ctx context.Context, req *pb.GetSingleSubmissionRequest) (*pb.GetSingleSubmissionReply, error) {
	submission, err := s.judgeUsecase.GetSingleSubmission(ctx, req.SubmissionUuid)
	if err != nil {
		return nil, err
	}
	var casesResults []*pb.CaseMetadata
	cases, err := s.judgeUsecase.GetCases(ctx, req.SubmissionUuid)
	if err != nil {
		return nil, err
	}
	for _, c := range cases {
		casesResults = append(casesResults, &pb.CaseMetadata{
			Score:      c.Score,
			Status:     c.Status,
			TimeCost:   int32(c.TimeCost),
			MemoryCost: int32(c.MemoryCost),
		})
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
		Code:       submission.Code,
		Language:   submission.Language,
		TimeCost:   int32(submission.TimeCost),
		MemoryCost: int32(submission.MemoryCost),
		CaseResults: &pb.GetSingleSubmissionReply_CaseResults{
			Cases: casesResults,
		},
	}, nil
}

func (s *UserService) GetSubmissions(ctx context.Context, req *pb.GetSubmissionsRequest) (*pb.GetSubmissionsReply, error) {
	submissions, err := s.judgeUsecase.GetSubmissions(ctx, req.GetContestId(), req.GetProblemId())
	if err != nil {
		return nil, err
	}
	var pbSubmissions []*pb.SubmissionMetadata
	for _, submission := range submissions {
		pbSubmissions = append(pbSubmissions, &pb.SubmissionMetadata{
			SubmissionUuid: submission.UUID,
			ProblemId:      submission.ProblemID,
			UserId:         submission.UserID,
			Status:         submission.Status,
			SubmitTime:     timestamppb.New(submission.CreateTime),
			Score:          int32(submission.Score),
		})
	}
	return &pb.GetSubmissionsReply{Submissions: pbSubmissions}, nil
}

func (s *UserService) GetSelfTestResult(ctx context.Context, req *pb.GetSelfTestResultRequest) (*pb.GetSelfTestResultReply, error) {
	selfTest, err := s.judgeUsecase.GetSelfTest(ctx, req.SelftestUuid)
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
