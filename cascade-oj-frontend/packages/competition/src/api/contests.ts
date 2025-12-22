import client from './client'
import type { ContestMetadata } from './types'
export type { ContestMetadata } from './types'

export async function getContests(): Promise<{ contests: ContestMetadata[] }> {
  return client.get('/user/contests')
}
