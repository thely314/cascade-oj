package service

import (
	"context"
	"strconv"
	"time"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/pkg/middleware/auth"
	"cascade-oj/pkg/util"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) PostSubmission(ctx context.Context, req *pb.SubmissionRequest) (*pb.SubmissionReply, error) {
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
	// TODO: get case version
	// caseVer, err := s.getProblemCaseVer(ctx, req.ProblemId)
	// if err != nil {
	// 	return nil, err
	// }
	_, err := s.judgeUsecase.CreateSubmission(ctx, &biz.Submission{
		UUID:        id,
		UserID:      userID,
		ProblemID:   req.ProblemId,
		Code:        util.Crlf2lf(req.Code),
		Language:    req.Language,
		Status:      "",
		Score:       0,
		CreateTime:  time.Now(),
		TimeCost:    0,
		MemoryCost:  0,
		CaseVersion: 0, //TODO: caseVer
	})
	if err != nil {
		return nil, err
	}
	return &pb.SubmissionReply{Uuid: id}, nil
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
		UUID:      id,
		UserID:    userID,
		ProblemID: req.ProblemId,
		Code:      util.Crlf2lf(req.Code),
		Language:  req.Language,
		Input:     util.Crlf2lf(req.SelfCase),
	})
	if err != nil {
		return nil, err
	}
	return &pb.SelfTestReply{Uuid: id}, nil
}

func (s *UserService) GetSingleSubmission(ctx context.Context, req *pb.GetSingleSubmissionRequest) (*pb.GetSingleSubmissionReply, error) {
	submission, err := s.judgeUsecase.GetSingleSubmission(ctx, req.SubmissionId)
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(submission.UUID, 10, 64)
	if err != nil {
		return nil, err
	}
	return &pb.GetSingleSubmissionReply{
		Metadata: &pb.SubmissionMetadata{
			SubmissionId: id,
			ProblemId:    submission.ProblemID,
			UserId:       submission.UserID,
			Status:       submission.Status,
			Result:       "",
			SubmitTime:   timestamppb.New(submission.CreateTime),
			Score:        int32(submission.Score),
		},
		Code:     submission.Code,
		Language: submission.Language,
	}, nil
}

func (s *UserService) GetSubmissions(ctx context.Context, req *pb.GetSubmissionsRequest) (*pb.GetSubmissionsReply, error) {
	submissions, err := s.judgeUsecase.GetSubmissions(ctx, req.GetContestId(), req.GetProblemId())
	if err != nil {
		return nil, err
	}
	var pbSubmissions []*pb.SubmissionMetadata
	for _, submission := range submissions {
		id, err := strconv.ParseInt(submission.UUID, 10, 64)
		if err != nil {
			return nil, err
		}
		pbSubmissions = append(pbSubmissions, &pb.SubmissionMetadata{
			SubmissionId: id,
			ProblemId:    submission.ProblemID,
			UserId:       submission.UserID,
			Status:       submission.Status,
			Result:       "",
			SubmitTime:   timestamppb.New(submission.CreateTime),
			Score:        int32(submission.Score),
		})
	}
	return &pb.GetSubmissionsReply{Submissions: pbSubmissions}, nil
}

// TODO: pb.GetSelfTestRequest and pb.GetSelfTestReply needed
// func (s *UserService) GetSelfTest(ctx context.Context, req *pb.GetSelfTestRequest) (*pb.GetSelfTestReply, error) {
// 	selfTest, err := s.judgeUsecase.GetSelfTest(ctx, req.GetSelfTestId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.GetSelfTestReply{
// 		IsCompiled: selfTest.IsCompiled,
// 		CompileMsg: selfTest.CompileMsg,
// 		Stdout:     selfTest.Stdout,
// 		Stderr:     selfTest.Stderr,
// 		TimeCost:   selfTest.TimeCost,
// 		MemoryCost: selfTest.MemoryCost,
// 	}, nil
// }

// TODO: pb.GetCasesRequest and pb.GetCasesReply needed
// func (s *UserService) GetCases(ctx context.Context, req *pb.GetCasesRequest) (*pb.GetCasesReply, error) {
// 	cases, err := s.judgeUsecase.GetCases(ctx, req.GetSubmissionId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	var pbCases []*pb.GetCasesReply_Case
// 	for _, c := range cases {
// 		pbCases = append(pbCases, &pb.GetCasesReply_Case{
// 			Index:      c.Index,
// 			Score:      c.Score,
// 			Status:     c.Status,
// 			TimeCost:   c.TimeCost,
// 			MemoryCost: c.MemoryCost,
// 		})
// 	}
// 	return &pb.GetCasesReply{
// 		Cases: pbCases,
// 	}, nil
// }
