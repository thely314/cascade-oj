import { get, post, put, del } from './request';
import type {
    GetContestsReply,
    GetSingleContestReply,
    PostContestRequest,
    PostContestReply,
    PutContestRequest,
    PutContestReply,
    DeleteContestReply,
    GetProblemsReply,
    GetSingleProblemReply,
    PostProblemRequest,
    PostProblemReply,
    PutProblemRequest,
    PutProblemReply,
    DeleteProblemReply,
    PublishProblemReply,
    GetSubmissionsRequest,
    GetSubmissionsReply,
    GetSingleSubmissionReply,
    RejudgeSubmissionReply,
    GetRanksReply,
    GetAnnouncementsReply,
    PostAnnouncementRequest,
    PostAnnouncementReply,
    PutAnnouncementRequest,
    PutAnnouncementReply,
    DeleteAnnouncementReply,
    GetUsersRequest,
    GetUsersReply,
    UpdateUserInfoRequest,
    UpdateUserInfoReply,
    DeleteUserReply,
    GetContestStatisticsReply,
    GetLogsRequest,
    GetLogsReply
} from './types';

// Contests
export const getContests = () =>
    get<GetContestsReply>('/admin/contests');

export const getSingleContest = (contestId: number) =>
    get<GetSingleContestReply>(`/admin/contests/${contestId}`);

export const postContest = (data: PostContestRequest) =>
    post<PostContestReply>('/admin/contests', data);

export const putContest = (contestId: number, data: PutContestRequest) =>
    put<PutContestReply>(`/admin/contests/${contestId}`, data);

export const deleteContest = (contestId: number) =>
    del<DeleteContestReply>(`/admin/contests/${contestId}`);

// Problems
export const getProblems = (contestId?: number) =>
    get<GetProblemsReply>('/admin/problems', { contest_id: contestId });

export const getSingleProblem = (problemId: number) =>
    get<GetSingleProblemReply>(`/admin/problems/${problemId}`);

export const postProblem = (data: PostProblemRequest) =>
    post<PostProblemReply>('/admin/problems', data);

export const putProblem = (problemId: number, data: PutProblemRequest) =>
    put<PutProblemReply>(`/admin/problems/${problemId}`, data);

export const deleteProblem = (problemId: number) =>
    del<DeleteProblemReply>(`/admin/problems/${problemId}`);

export const publishProblem = (problemId: number) =>
    post<PublishProblemReply>(`/admin/problems/${problemId}/publish`);

// Submissions
export const getSubmissions = (params: GetSubmissionsRequest) =>
    get<GetSubmissionsReply>('/admin/submissions', params as Record<string, any>);

export const getSingleSubmission = (submissionUuid: string) =>
    get<GetSingleSubmissionReply>(`/admin/submissions/${submissionUuid}`);

export const rejudgeSubmission = (submissionUuid: string) =>
    post<RejudgeSubmissionReply>(`/admin/submissions/${submissionUuid}/rejudge`);

// Ranks
export const getRanks = (contestId: number) =>
    get<GetRanksReply>(`/admin/${contestId}/ranks`);

// Announcements
export const getAnnouncements = () =>
    get<GetAnnouncementsReply>('/admin/announcements');

export const postAnnouncement = (data: PostAnnouncementRequest) =>
    post<PostAnnouncementReply>('/admin/announcements', data);

export const putAnnouncement = (announcementId: number, data: PutAnnouncementRequest) =>
    put<PutAnnouncementReply>(`/admin/announcements/${announcementId}`, data);

export const deleteAnnouncement = (announcementId: number) =>
    del<DeleteAnnouncementReply>(`/admin/announcements/${announcementId}`);

// Users
export const getUsers = (params: GetUsersRequest) =>
    get<GetUsersReply>('/admin/users', params as Record<string, any>);

export const updateUserInfo = (userId: number, data: UpdateUserInfoRequest) =>
    put<UpdateUserInfoReply>(`/admin/users/${userId}`, data);

export const deleteUser = (userId: number) =>
    del<DeleteUserReply>(`/admin/users/${userId}`);

// Statistics
export const getContestStatistics = (contestId: number) =>
    get<GetContestStatisticsReply>(`/admin/contests/${contestId}/statistics`);

// Logs
export const getLogs = (params: GetLogsRequest) =>
    get<GetLogsReply>('/admin/logs', params as Record<string, any>);
