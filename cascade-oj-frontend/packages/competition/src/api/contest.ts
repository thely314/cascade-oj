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

export async function getContestProblems(contestId: string): Promise<{ problems: ProblemMetadata[] }> {
  return client.get(`/user/problems/${contestId}`)
}

export async function joinContest(contestId: string, userId: string): Promise<{ isJoin: boolean }>{
  return client.post(`/user/contests/${contestId}/join`, { contestId, userId })
}

export async function quitContest(contestId: string, userId?: string): Promise<{ isJoin: boolean }>{
  return client.delete(`/user/contests/${contestId}/join`, { params: { userId } })
}
