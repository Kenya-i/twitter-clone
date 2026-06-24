'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '../context/AuthContext'
import { API_URL } from '../../lib/api'
import { toggleLike } from '../../lib/tweets'
import { useRequireAuth } from '../../hooks/useRequireAuth'
import { Tweet } from '../../types/tweet'
import TweetCard from '../../components/TweetCard'

export default function Search() {
  const router = useRouter()
  const { token, userId } = useAuth()
  const [query, setQuery] = useState('')
  const [tweets, setTweets] = useState<Tweet[]>([])
  const [nextCursor, setNextCursor] = useState<string | null>(null)
  const [loadingMore, setLoadingMore] = useState(false)
  const [searched, setSearched] = useState(false)

  useRequireAuth(token)

  const fetchResults = async (q: string, cursor?: string) => {
    const params = new URLSearchParams({ q })
    if (cursor) params.set('cursor', cursor)

    const res = await fetch(`${API_URL}/tweets/search?${params.toString()}`, {
      headers: { Authorization: `Bearer ${token}` },
    })

    if (res.ok) {
      const data = await res.json()
      setTweets((prev) => (cursor ? [...prev, ...data.tweets] : data.tweets))
      setNextCursor(data.next_cursor)
    }
  }

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!query.trim()) return
    setSearched(true)
    await fetchResults(query)
  }

  const handleLoadMore = async () => {
    if (!nextCursor) return
    setLoadingMore(true)
    await fetchResults(query, nextCursor)
    setLoadingMore(false)
  }

  const handleLike = async (tweet: Tweet) => {
    const wasLiked = tweet.liked_by_me

    setTweets((prev) =>
      prev.map((t) =>
        t.id === tweet.id
          ? { ...t, liked_by_me: !wasLiked, like_count: t.like_count + (wasLiked ? -1 : 1) }
          : t
      )
    )

    const ok = await toggleLike(tweet.id, wasLiked, token)

    if (!ok) {
      setTweets((prev) =>
        prev.map((t) =>
          t.id === tweet.id
            ? { ...t, liked_by_me: wasLiked, like_count: t.like_count + (wasLiked ? 1 : -1) }
            : t
        )
      )
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-md mx-auto bg-white p-6 rounded-lg shadow-md">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">検索</h1>
          <Link href="/timeline" className="text-sm text-blue-500 hover:underline">
            タイムラインに戻る
          </Link>
        </div>

        <form onSubmit={handleSearch} className="flex gap-2 mb-4">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="flex-1 border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="キーワードを入力"
          />
          <button
            type="submit"
            className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 transition-colors"
          >
            検索
          </button>
        </form>

        <div className="space-y-3">
          {tweets.map((tweet) => (
            <TweetCard key={tweet.id} tweet={tweet} currentUserId={userId} onLike={handleLike} />
          ))}
          {searched && tweets.length === 0 && (
            <p className="text-sm text-gray-400 text-center">該当するツイートが見つかりません</p>
          )}
        </div>

        {nextCursor && (
          <button
            onClick={handleLoadMore}
            disabled={loadingMore}
            className="w-full text-sm text-blue-500 hover:underline mt-4 disabled:text-gray-400"
          >
            {loadingMore ? '読み込み中...' : 'もっと読み込む'}
          </button>
        )}
      </div>
    </div>
  )
}
