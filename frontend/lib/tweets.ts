import { API_URL } from './api'

export async function toggleLike(
  tweetId: string,
  isLiked: boolean,
  token: string | null
): Promise<boolean> {
  const res = await fetch(`${API_URL}/tweets/${tweetId}/like`, {
    method: isLiked ? 'DELETE' : 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  return res.ok
}
