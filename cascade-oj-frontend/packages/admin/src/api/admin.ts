import service from './request';
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
    DisableProblemReply,
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
    GetContestUsersReply,
    AddContestUserReply,
    RemoveContestUserReply,
    UpdateUserPasswordRequest,
    UpdateUserPasswordReply,
    GetContestStatisticsReply,
    ListLogFilesRequest,
    ListLogFilesReply,
    QueryLogContentRequest,
    QueryLogContentReply,
    DownloadLogsRequest,
} from './types';

// Contests
export async function getContests(): Promise<GetContestsReply> {
    return service.get('/admin/contests');
}

export async function getSingleContest(contestId: number): Promise<GetSingleContestReply> {
    return service.get(`/admin/contests/${contestId}`);
}

export async function postContest(data: PostContestRequest): Promise<PostContestReply> {
    return service.post('/admin/contests', data);
}

export async function putContest(contestId: number, data: PutContestRequest): Promise<PutContestReply> {
    return service.put(`/admin/contests/${contestId}`, data);
}

export async function deleteContest(contestId: number): Promise<DeleteContestReply> {
    return service.delete(`/admin/contests/${contestId}`);
}

// Problems
export async function getProblems(contestId?: number): Promise<GetProblemsReply> {
    return service.get('/admin/problems', { params: { contest_id: contestId } });
}

export async function getSingleProblem(problemId: number): Promise<GetSingleProblemReply> {
    return service.get(`/admin/problems/${problemId}`);
}

export async function postProblem(data: PostProblemRequest): Promise<PostProblemReply> {
    return service.post('/admin/problems', data);
}

export async function putProblem(problemId: number, data: PutProblemRequest): Promise<PutProblemReply> {
    return service.put(`/admin/problems/${problemId}`, data);
}

export async function deleteProblem(problemId: number): Promise<DeleteProblemReply> {
    return service.delete(`/admin/problems/${problemId}`);
}

export async function publishProblem(problemId: number): Promise<PublishProblemReply> {
    return service.post(`/admin/problems/${problemId}/publish`, {});
}

export async function disableProblem(problemId: number): Promise<DisableProblemReply> {
    return service.post(`/admin/problems/${problemId}/disable`, {});
}

// Submissions
export async function getSubmissions(params: GetSubmissionsRequest): Promise<GetSubmissionsReply> {
    return service.get('/admin/submissions', { params: params as Record<string, any> });
}

export async function getSingleSubmission(submissionUuid: string): Promise<GetSingleSubmissionReply> {
    return service.get(`/admin/submissions/${submissionUuid}`);
}

export async function rejudgeSubmission(submissionUuid: string): Promise<RejudgeSubmissionReply> {
    return service.post(`/admin/submissions/${submissionUuid}/rejudge`);
}

// Ranks
export async function getRanks(contestId: number): Promise<GetRanksReply> {
    return service.get(`/admin/${contestId}/ranks`);
}

// Announcements
export async function getAnnouncements(): Promise<GetAnnouncementsReply> {
    return service.get('/admin/announcements');
}

export async function postAnnouncement(data: PostAnnouncementRequest): Promise<PostAnnouncementReply> {
    return service.post('/admin/announcements', data);
}

export async function putAnnouncement(announcementId: number, data: PutAnnouncementRequest): Promise<PutAnnouncementReply> {
    return service.put(`/admin/announcements/${announcementId}`, data);
}

export async function deleteAnnouncement(announcementId: number): Promise<DeleteAnnouncementReply> {
    return service.delete(`/admin/announcements/${announcementId}`);
}

// Users
export async function getUsers(params: GetUsersRequest): Promise<GetUsersReply> {
    return service.get('/admin/users', { params: params as Record<string, any> });
}

export async function updateUserInfo(userId: number, data: UpdateUserInfoRequest): Promise<UpdateUserInfoReply> {
    return service.put(`/admin/users/${userId}`, data);
}

export async function deleteUser(userId: number): Promise<DeleteUserReply> {
    return service.delete(`/admin/users/${userId}`);
}

export async function getContestUsers(contestId: number): Promise<GetContestUsersReply> {
    return service.get(`/admin/contests/${contestId}/users`);
}

export async function addContestUser(contestId: number, userId: number): Promise<AddContestUserReply> {
    return service.post(`/admin/contests/${contestId}/users/${userId}`);
}

export async function removeContestUser(contestId: number, userId: number): Promise<RemoveContestUserReply> {
    return service.delete(`/admin/contests/${contestId}/users/${userId}`);
}

export async function updateUserPassword(userId: number, data: UpdateUserPasswordRequest): Promise<UpdateUserPasswordReply> {
    return service.put(`/admin/users/${userId}/password`, data);
}
// Statistics
export async function getContestStatistics(contestId: number): Promise<GetContestStatisticsReply> {
    return service.get(`/admin/contests/${contestId}/statistics`);
}

// Logs
export async function listLogFiles(params: ListLogFilesRequest): Promise<ListLogFilesReply> {
    return service.get('/admin/logs/files', { params: params as Record<string, any> });
}

export async function queryLogContent(params: QueryLogContentRequest): Promise<QueryLogContentReply> {
    return service.get('/admin/logs/content', { params: params as Record<string, any> });
}

export async function downloadLogs(data: DownloadLogsRequest): Promise<any> {
    return service.post('/admin/logs/download', data, { responseType: 'json' });
}
