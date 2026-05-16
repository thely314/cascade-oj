package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/pkg/middleware/auth"
	"context"
)

func (adminService *AdminService) GetAnnouncements(ctx context.Context, request *pb.GetAnnouncementsRequest) (*pb.GetAnnouncementsReply, error) {
	announcements, err := adminService.announcementUseCase.GetAnnouncements(ctx)
	if err != nil {
		return nil, err
	}
	pbAnnouncements := make([]*pb.GetAnnouncementsReply_Announcement, 0, len(announcements))
	for _, anannouncement := range announcements {
		pbAnnouncements = append(pbAnnouncements,
			&pb.GetAnnouncementsReply_Announcement{
				Id:            anannouncement.ID,
				PublisherName: anannouncement.PublisherName,
				Title:         anannouncement.Title,
				Content:       anannouncement.Content,
			},
		)
	}
	return &pb.GetAnnouncementsReply{
		Announcements: pbAnnouncements,
	}, nil
}

func (adminService *AdminService) PostAnnouncement(ctx context.Context, request *pb.PostAnnouncementRequest) (*pb.PostAnnouncementReply, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	announcementID, err := adminService.announcementUseCase.PostAnnouncement(ctx, biz.AnnouncementCreateInfo{
		PublisherID: userID,
		Title:       request.Title,
		Content:     request.Content,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PostAnnouncementReply{AnnouncementId: announcementID}, nil
}

func (adminService *AdminService) PutAnnouncement(ctx context.Context, request *pb.PutAnnouncementRequest) (*pb.PutAnnouncementReply, error) {
	isUpdated, err := adminService.announcementUseCase.PutAnnouncement(ctx, biz.AnnouncementEditInfo{
		AnnouncementID: request.AnnouncementId,
		Title:          request.Title,
		Content:        request.Content,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PutAnnouncementReply{IsUpdated: isUpdated}, nil
}

func (adminService *AdminService) DeleteAnnouncement(ctx context.Context, request *pb.DeleteAnnouncementRequest) (*pb.DeleteAnnouncementReply, error) {
	isDeleted, err := adminService.announcementUseCase.DeleteAnnouncement(ctx, request.AnnouncementId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAnnouncementReply{IsDeleted: isDeleted}, nil
}
