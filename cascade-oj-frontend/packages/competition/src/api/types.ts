// Shared API types aligned with OpenAPI

export type ContestMetadata = {
  id: string
  title: string
  startTime: string // ISO date-time
  endTime: string   // ISO date-time
  status: string
}

export type ProblemMetadata = {
  id: string
  title: string
  timeLimitMs: number
  memoryLimitMb: number
}

export type Announcement = {
  id: string
  publisherName: string
  title: string
  content: string
}

export type RankItem = {
  userId: string
  username: string
  rank: number
  score: number
}
