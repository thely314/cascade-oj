package service

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"
)

func (gatewayService *GatewayService) GetRanks(ctx context.Context, req *pb.GetRanksRequest) (*pb.GetRanksReply, error) {
	ranks, err := gatewayService.miscUsecase.GetRanks(ctx, req.ContestId)
	if err != nil {
		return nil, err
	}
	res := &pb.GetRanksReply{}
	for _, rank := range ranks {
		res.Ranks = append(res.Ranks, &pb.GetRanksReply_RankItem{
			UserId:   rank.UserID,
			Username: rank.Username,
			Rank:     rank.UserRank,
			Score:    rank.Score,
		})
	}
	return res, nil
}

func (gatewayService *GatewayService) GetAnnouncements(ctx context.Context, req *pb.GetAnnouncementsRequest) (
	*pb.GetAnnouncementsReply, error) {
	announcements, err := gatewayService.miscUsecase.GetAnnouncements(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.GetAnnouncementsReply{}
	for _, announcement := range announcements {
		res.Announcements = append(res.Announcements, &pb.GetAnnouncementsReply_Announcement{
			Id:            announcement.ID,
			PublisherName: announcement.PublisherName,
			Title:         announcement.Title,
			Content:       announcement.Content,
		})
	}
	return res, nil
}

func (gatewayService *GatewayService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	userInfo, err := gatewayService.miscUsecase.GetUserInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserInfoReply{
		UserId:   userInfo.UserID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}, nil
}

func (gatewayService *GatewayService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoReply, error) {
	err := gatewayService.miscUsecase.UpdateUserInfo(ctx, &biz.UserInfo{
		UserID:   req.UserId,
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		return &pb.UpdateUserInfoReply{
			IsUpdated: false,
		}, err
	}
	return &pb.UpdateUserInfoReply{
		IsUpdated: true,
	}, nil
}
