import client from './client'
import type { Announcement } from './types'
export type { Announcement } from './types'

export async function getAnnouncements(): Promise<{ announcements: Announcement[] }> {
  return client.get('/user/announcements')
}
