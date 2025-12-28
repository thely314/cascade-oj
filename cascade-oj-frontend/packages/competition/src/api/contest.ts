import client from './client'
import type { ContestMetadata, ProblemMetadata } from './types'
export type { ContestMetadata, ProblemMetadata } from './types'

export type GetSingleContestReply = {
  metadata: ContestMetadata
  description: string
}

export async function getContest(contestId: string): Promise<GetSingleContestReply> {
  return client.get(`/user/contests/${contestId}`)
}

// 获取比赛题目列表
export async function getContestProblems(contestId: string): Promise<{ problems: ProblemMetadata[] }> {
  return client.get(`/user/contests/${contestId}/problems`)
}

export async function joinContest(contestId: string, userId?: string): Promise<{ isJoin: boolean }>{
  const body: Record<string, any> = { contestId }
  if (typeof userId !== 'undefined') body.userId = userId
  return client.post(`/user/contests/${contestId}/join`, body)
}

export async function quitContest(contestId: string, userId?: string): Promise<{ isJoin: boolean }>{
  return client.delete(`/user/contests/${contestId}/join`, { params: { userId } })
}

// 占位：查询是否已加入比赛（后端将提供该接口）
export async function getJoinStatus(contestId: string, userId?: string): Promise<{ isJoin: boolean }>{
  return client.get(`/user/contests/${contestId}/join`, { params: { userId } })
}
