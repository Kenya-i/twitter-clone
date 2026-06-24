'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuth } from '../context/AuthContext'
import Link from 'next/link'
import { API_URL } from '../../lib/api'

type Tweet = {
  id: string
  user_id: string
  content: string
  created_at: string
  updated_at: string
  like_count: number
  liked_by_me: boolean
}

export default function Timeline() {
  const router = useRouter()
  const { token, userId, logout } = useAuth()
  const [content, setContent] = useState('')
  const [message, setMessage] = useState('')
  const [tweets, setTweets] = useState<Tweet[]>([])
  const [nextCursor, setNextCursor] = useState<string | null>(null)
  const [loadingMore, setLoadingMore] = useState(false)


  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])

  const fetchTweets = async (cursor?: string) => {
    const url = cursor
      ? `${API_URL}/tweets?cursor=${encodeURIComponent(cursor)}`
      : `${API_URL}/tweets`

    const res = await fetch(url, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.ok) {
      const data = await res.json()
      setTweets((prev) => (cursor ? [...prev, ...data.tweets] : data.tweets))
      setNextCursor(data.next_cursor)
    }
  }

  const handleLoadMore = async () => {
    if (!nextCursor) return
    setLoadingMore(true)
    await fetchTweets(nextCursor)
    setLoadingMore(false)
  }

  useEffect(() => {
    if (token) fetchTweets()
  }, [token])

  const handlePost = async (e: React.FormEvent) => {
    e.preventDefault()
    setMessage('')

    const res = await fetch(`${API_URL}/tweets`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ content }),
    })

    if (!res.ok) {
      setMessage('投稿に失敗しました')
      return
    }

    setContent('')
    setMessage('投稿しました！')
    fetchTweets()
  }

  const handleDelete = async (id: string) => {
    const res = await fetch(`${API_URL}/tweets/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.ok) {
      fetchTweets()
    }
  }

  const [editingId, setEditingId] = useState<string | null>(null)
  const [editContent, setEditContent] = useState('')

  const handleEditStart = (tweet: Tweet) => {
    setEditingId(tweet.id)
    setEditContent(tweet.content)
  }

  const handleEditCancel = () => {
    setEditingId(null)
    setEditContent('')
  }

  const handleUpdate = async (id: string) => {
    const res = await fetch(`${API_URL}/tweets/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ content: editContent }),
    })

    if (res.ok) {
      setEditingId(null)
      setEditContent('')
      fetchTweets()
    }
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

    const res = await fetch(`${API_URL}/tweets/${tweet.id}/like`, {
      method: wasLiked ? 'DELETE' : 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (!res.ok) {
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
          <h1 className="text-xl font-bold">タイムライン</h1>
          <div className="flex items-center gap-3">
            <Link href="/search" className="text-sm text-blue-500 hover:underline">
              検索
            </Link>
            <Link href="/users" className="text-sm text-blue-500 hover:underline">
              ユーザーを探す
            </Link>
            <Link href="/profile" className="text-sm text-blue-500 hover:underline">
              プロフィール
            </Link>
            <button
              onClick={() => {
                logout()
                router.push('/')
              }}
              className="text-sm text-gray-500 hover:text-gray-700"
            >
              ログアウト
            </button>
          </div>
        </div>
        <form onSubmit={handlePost} className="space-y-2">
          {message && <p className="text-sm text-green-600">{message}</p>}
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="いまどうしてる？"
            rows={3}
          />
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 transition-colors"
          >
            ツイートする
          </button>
        </form>
        <div className="mt-6 space-y-3">
          {tweets.map((tweet) => (
            <div key={tweet.id} className="border-b border-gray-200 pb-2 flex justify-between items-start">
              <div className="flex-1">
                <Link href={`/users/${tweet.user_id}`} className="text-xs text-blue-500 hover:underline">
                  投稿者のプロフィール
                </Link>
                {editingId === tweet.id ? (
                  <div className="mt-1 space-y-2">
                    <textarea
                      value={editContent}
                      onChange={(e) => setEditContent(e.target.value)}
                      className="w-full border border-gray-300 rounded-md px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      rows={2}
                    />
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleUpdate(tweet.id)}
                        className="text-xs text-blue-500 hover:text-blue-700"
                      >
                        保存
                      </button>
                      <button
                        onClick={handleEditCancel}
                        className="text-xs text-gray-500 hover:text-gray-700"
                      >
                        キャンセル
                      </button>
                    </div>
                  </div>
                ) : (
                  <>
                    <p className="text-sm mt-1">{tweet.content}</p>
                    <p className="text-xs text-gray-400 mt-1">
                      {new Date(tweet.created_at).toLocaleString()}
                    </p>
                    <div className="flex items-center gap-3 mt-1">
                      <button
                        onClick={() => handleLike(tweet)}
                        className={`text-xs flex items-center gap-1 ${
                          tweet.liked_by_me ? 'text-pink-500' : 'text-gray-400 hover:text-pink-500'
                        }`}
                      >
                        {tweet.liked_by_me ? '♥' : '♡'} {tweet.like_count}
                      </button>
                      <Link href={`/tweets/${tweet.id}`} className="text-xs text-blue-500 hover:underline">
                        詳細
                      </Link>
                    </div>
                  </>
                )}
              </div>
              {tweet.user_id === userId && editingId !== tweet.id && (
                <div className="flex gap-2 ml-2">
                  <button
                    onClick={() => handleEditStart(tweet)}
                    className="text-xs text-blue-500 hover:text-blue-700"
                  >
                    編集
                  </button>
                  <button
                    onClick={() => handleDelete(tweet.id)}
                    className="text-xs text-red-500 hover:text-red-700"
                  >
                    削除
                  </button>
                </div>
              )}
            </div>
          ))}
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
