package service

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/pkg/middleware/auth"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (gatewayService *GatewayService) GetContests(ctx context.Context, req *pb.GetContestsRequest) (*pb.GetContestsReply, error) {
	contests, err := gatewayService.contestUsecase.GetContests(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.GetContestsReply{}
	for _, v := range contests {
		res.Contests = append(res.Contests, &pb.ContestMetadata{
			Id:        v.ID,
			Title:     v.Title,
			StartTime: timestamppb.New(v.StartTime),
			EndTime:   timestamppb.New(v.EndTime),
			Status:    v.Status,
		})
	}
	return res, nil
}

func (gatewayService *GatewayService) GetSingleContest(ctx context.Context, req *pb.GetSingleContestRequest) (
	*pb.GetSingleContestReply, error) {
	contest, err := gatewayService.contestUsecase.GetSingleContest(ctx, req.ContestId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSingleContestReply{
		Metadata: &pb.ContestMetadata{
			Id:        contest.Metadata.ID,
			Title:     contest.Metadata.Title,
			StartTime: timestamppb.New(contest.Metadata.StartTime),
			EndTime:   timestamppb.New(contest.Metadata.EndTime),
			Status:    contest.Metadata.Status,
		},
		Description: contest.Description,
	}, nil
}

func (gatewayService *GatewayService) JoinContest(ctx context.Context, req *pb.JoinContestRequest) (*pb.JoinContestReply, error) {
	err := gatewayService.contestUsecase.JoinContest(ctx, req.ContestId, req.UserId)
	if err != nil {
		return &pb.JoinContestReply{
			IsJoin: false,
		}, err
	}
	return &pb.JoinContestReply{
		IsJoin: true,
	}, nil
}

func (gatewayService *GatewayService) QuitContest(ctx context.Context, req *pb.QuitContestRequest) (*pb.QuitContestReply, error) {
	err := gatewayService.contestUsecase.QuitContest(ctx, req.ContestId, req.UserId)
	if err != nil {
		return &pb.QuitContestReply{
			IsJoin: true,
		}, err
	}
	return &pb.QuitContestReply{
		IsJoin: false,
	}, nil
}

func (gatewayService *GatewayService) GetProblems(ctx context.Context, req *pb.GetProblemsRequest) (*pb.GetProblemsReply, error) {
	problems, err := gatewayService.contestUsecase.GetProblems(ctx, req.ContestId)
	if err != nil {
		return nil, err
	}
	res := &pb.GetProblemsReply{}
	for _, v := range problems {
		res.Problems = append(res.Problems, &pb.ProblemMetadata{
			Id:            v.ID,
			Title:         v.Title,
			ProblemType:   v.ProblemType,
			TimeLimitMs:   v.TimeLimitMs,
			MemoryLimitMb: v.MemoryLimitMb,
		})
	}
	return res, nil
}

func (gatewayService *GatewayService) GetSingleProblem(ctx context.Context, req *pb.GetSingleProblemRequest) (
	*pb.GetSingleProblemReply, error) {
	problem, err := gatewayService.contestUsecase.GetSingleProblem(ctx, ctx.Value("userInfo").(auth.Claims).UserID, req.ProblemId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSingleProblemReply{
		Metadata: &pb.ProblemMetadata{
			Id:            problem.Metadata.ID,
			Title:         problem.Metadata.Title,
			ProblemType:   problem.Metadata.ProblemType,
			TimeLimitMs:   problem.Metadata.TimeLimitMs,
			MemoryLimitMb: problem.Metadata.MemoryLimitMb,
		},
		Creator:     problem.Creator,
		Description: problem.Description,
	}, nil
}
