import client from './client'
import type { RankItem } from './types'
export type { RankItem } from './types'

export async function getRanks(contestId: string): Promise<{ ranks: RankItem[] }> {
  return client.get(`/user/${contestId}/ranks`)
}
