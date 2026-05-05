package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/pkg/middleware/auth"
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (adminService *AdminService) GetContests(ctx context.Context, request *pb.GetContestsRequest) (*pb.GetContestsReply, error) {
	userInfo, ok := ctx.Value("userInfo").(*auth.Claims)
	if !ok {
		log.Error("Failed to get user info from context")
		return nil, errors.New("failed to get user info")
	}

	contests, err := adminService.contestUseCase.GetContests(ctx, userInfo.UserID)
	if err != nil {
		return nil, err
	}

	pbContests := make([]*pb.ContestMetadata, 0, len(contests))
	for _, contest := range contests {
		pbContests = append(pbContests,
			&pb.ContestMetadata{
				Id:        contest.ID,
				Title:     contest.Title,
				StartTime: timestamppb.New(contest.StartTime),
				EndTime:   timestamppb.New(contest.EndTime),
				Status:    contest.Status,
			},
		)
	}

	return &pb.GetContestsReply{
		Contests: pbContests,
	}, nil
}

func (adminService *AdminService) GetSingleContest(ctx context.Context, request *pb.GetSingleContestRequest) (*pb.GetSingleContestReply, error) {
	contest, err := adminService.contestUseCase.GetSingleContest(ctx, request.ContestId)
	if err != nil {
		return nil, err
	}

	log.Infof("GetSingleContest - contestID: %d, problemIDs: %v", request.ContestId, contest.ProblemIDs)

	return &pb.GetSingleContestReply{
		Metadata: &pb.ContestMetadata{
			Id:        contest.Contest.ID,
			Title:     contest.Contest.Title,
			StartTime: timestamppb.New(contest.Contest.StartTime),
			EndTime:   timestamppb.New(contest.Contest.EndTime),
			Status:    contest.Contest.Status,
		},
		Description: contest.Description,
		ProblemIds:  contest.ProblemIDs,
	}, nil
}

func (adminService *AdminService) PostContest(ctx context.Context, request *pb.PostContestRequest) (*pb.PostContestReply, error) {
	userInfo, ok := ctx.Value("userInfo").(*auth.Claims)
	if !ok {
		log.Error("Failed to get user info from context")
		return nil, errors.New("failed to get user info")
	}

	if request.Problems == nil || len(request.Problems.ProblemIds) == 0 {
		return nil, errors.New("no problems selected")
	}

	contestID, err := adminService.contestUseCase.PostContest(
		ctx,
		biz.ContestCreateInfo{
			Title:         request.Title,
			Description:   request.Description,
			StartTime:     request.StartTime.AsTime(),
			EndTime:       request.EndTime.AsTime(),
			ProblemIDList: request.Problems.ProblemIds,
		},
		userInfo.UserID,
	)
	if err != nil {
		return nil, err
	}

	if contestID == -1 {
		return nil, errors.New("failed to create contest")
	}

	return &pb.PostContestReply{
		ContestId: contestID,
	}, nil
}

func (adminService *AdminService) PutContest(ctx context.Context, request *pb.PutContestRequest) (*pb.PutContestReply, error) {
	log.Infof("PutContest - contestID: %d, Problems: %v", request.ContestId, request.Problems)
	var problemIDList []int64
	if request.Problems != nil {
		log.Infof("PutContest - ProblemIds: %v", request.Problems.ProblemIds)
		problemIDList = request.Problems.ProblemIds
	}

	isUpdated, err := adminService.contestUseCase.PutContest(
		ctx,
		biz.ContestEditInfo{
			ID:            request.ContestId,
			Title:         request.Title,
			Description:   request.Description,
			StartTime:     request.StartTime.AsTime(),
			EndTime:       request.EndTime.AsTime(),
			ProblemIDList: problemIDList,
		},
	)
	if err != nil {
		return nil, err
	}

	return &pb.PutContestReply{
		IsUpdated: isUpdated,
	}, nil
}

func (adminService *AdminService) DeleteContest(ctx context.Context, request *pb.DeleteContestRequest) (*pb.DeleteContestReply, error) {
	isDeleted, err := adminService.contestUseCase.DeleteContest(ctx, request.ContestId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteContestReply{
		IsDeleted: isDeleted,
	}, nil
}

func (adminService *AdminService) GetContestCompetitors(ctx context.Context, request *pb.GetContestCompetitorsRequest) (*pb.GetContestCompetitorsReply, error) {
	competitors, err := adminService.contestUseCase.GetContestCompetitors(ctx, request.ContestId)
	if err != nil {
		return nil, err
	}

	pbCompetitors := make([]*pb.UserInfo, 0, len(competitors))
	for _, c := range competitors {
		pbCompetitors = append(pbCompetitors,
			&pb.UserInfo{
				UserId:   c.UserID,
				Username: c.Username,
				Email:    c.Email,
			})
	}

	return &pb.GetContestCompetitorsReply{
		Competitors: pbCompetitors,
	}, nil
}

func (adminService *AdminService) PutContestCompetitors(ctx context.Context, request *pb.PutContestCompetitorsRequest) (*pb.PutContestCompetitorsReply, error) {
	isUpdated, err := adminService.contestUseCase.PutContestCompetitors(ctx, request.ContestId, request.UserIds)
	if err != nil {
		return nil, err
	}

	return &pb.PutContestCompetitorsReply{
		IsUpdated: isUpdated,
	}, nil
}
