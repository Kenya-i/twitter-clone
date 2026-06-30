'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '../context/AuthContext'
import Link from 'next/link'
import { API_URL } from '../../lib/api'
import { toggleLike } from '../../lib/tweets'
import { useRequireAuth } from '../../hooks/useRequireAuth'
import { Tweet } from '../../types/tweet'
import TweetCard from '../../components/TweetCard'

export default function Timeline() {
  const router = useRouter()
  const { token, userId, logout } = useAuth()
  const [content, setContent] = useState('')
  const [message, setMessage] = useState('')
  const [tweets, setTweets] = useState<Tweet[]>([])
  const [nextCursor, setNextCursor] = useState<string | null>(null)
  const [loadingMore, setLoadingMore] = useState(false)
  const [imageFiles, setImageFiles] = useState<File[]>([])

  useRequireAuth(token)

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

    const formData = new FormData()
    formData.append('content', content)
    imageFiles.forEach((file) => formData.append('images', file))

    const res = await fetch(`${API_URL}/tweets`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    })

    if (!res.ok) {
      setMessage('投稿に失敗しました')
      return
    }

    setContent('')
    setImageFiles([])
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
          <input
            type="file"
            accept="image/*"
            multiple
            onChange={(e) => setImageFiles(Array.from(e.target.files ?? []))}
            className="text-xs"
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
            <TweetCard
              key={tweet.id}
              tweet={tweet}
              currentUserId={userId}
              onLike={handleLike}
              isEditing={editingId === tweet.id}
              editContent={editContent}
              onEditContentChange={setEditContent}
              onEditStart={handleEditStart}
              onEditCancel={handleEditCancel}
              onEditSave={handleUpdate}
              onDelete={handleDelete}
            />
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