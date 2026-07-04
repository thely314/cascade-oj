// import client from './client'
import request from '@/utils/request';

import type { ContestMetadata, ProblemMetadata } from './types'
export type { ContestMetadata, ProblemMetadata } from './types'

export type GetSingleContestReply = {
  metadata: ContestMetadata
  description: string
}

export async function getContest(contestId: string): Promise<GetSingleContestReply> {
  return (await request.get(`/user/contests/${contestId}`))?.data;
}

// 获取比赛题目列表
export async function getContestProblems(contestId: string): Promise<{ problems: ProblemMetadata[] }> {
  return (await request.get(`/user/contests/${contestId}/problems`))?.data;
}

// adjust API according to proto definition
export async function joinContest(contestId: string, userId?: string): Promise<{ isJoin: boolean }>{
  const body: Record<string, any> = { contestId }
  if (typeof userId !== 'undefined') body.userId = userId
  return (await request.post(`/user/contests/${contestId}/join`, body))?.data;
}

export async function quitContest(contestId: string, userId?: string): Promise<{ isJoin: boolean }>{
  return (await request.delete(`/user/contests/${contestId}/join`, { params: { userId } }))?.data;
}

// 查询是否已加入比赛（后端将提供该接口）
export async function getJoinStatus(contestId: string, userId?: string): Promise<{ isJoin: boolean }>{
  return (await request.get(`/user/contests/${contestId}/join`, { params: { userId } }))?.data;
}
