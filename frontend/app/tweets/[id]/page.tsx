'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '../../context/AuthContext'
import { API_URL } from '../../../lib/api'

type Tweet = {
  id: string
  user_id: string
  content: string
  created_at: string
  updated_at: string
  like_count: number
  liked_by_me: boolean
}

export default function TweetDetail() {
  const params = useParams<{ id: string }>()
  const router = useRouter()
  const { token, userId } = useAuth()
  const [tweet, setTweet] = useState<Tweet | null>(null)
  const [isEditing, setIsEditing] = useState(false)
  const [editContent, setEditContent] = useState('')

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])

  useEffect(() => {
    if (!token) return

    const fetchTweet = async () => {
      const res = await fetch(`${API_URL}/tweets/${params.id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (res.ok) {
        const data = await res.json()
        setTweet(data)
      }
    }

    fetchTweet()
  }, [token, params.id])

  const handleLike = async () => {
    if (!tweet) return
    const wasLiked = tweet.liked_by_me

    setTweet({
      ...tweet,
      liked_by_me: !wasLiked,
      like_count: tweet.like_count + (wasLiked ? -1 : 1),
    })

    const res = await fetch(`${API_URL}/tweets/${tweet.id}/like`, {
      method: wasLiked ? 'DELETE' : 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (!res.ok) {
      setTweet({
        ...tweet,
        liked_by_me: wasLiked,
        like_count: tweet.like_count + (wasLiked ? 1 : -1),
      })
    }
  }

  const handleEditStart = () => {
    if (!tweet) return
    setEditContent(tweet.content)
    setIsEditing(true)
  }

  const handleEditCancel = () => {
    setIsEditing(false)
    setEditContent('')
  }

  const handleUpdate = async () => {
    if (!tweet) return

    const res = await fetch(`${API_URL}/tweets/${tweet.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ content: editContent }),
    })

    if (res.ok) {
      const data = await res.json()
      setTweet(data)
      setIsEditing(false)
      setEditContent('')
    }
  }

  const handleDelete = async () => {
    if (!tweet) return

    const res = await fetch(`${API_URL}/tweets/${tweet.id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.ok) {
      router.push('/timeline')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-md mx-auto bg-white p-6 rounded-lg shadow-md">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">ツイート詳細</h1>
          <Link href="/timeline" className="text-sm text-blue-500 hover:underline">
            タイムラインに戻る
          </Link>
        </div>

        {tweet ? (
          <div className="space-y-2">
            <Link href={`/users/${tweet.user_id}`} className="text-xs text-blue-500 hover:underline">
              投稿者のプロフィール
            </Link>

            {isEditing ? (
              <div className="space-y-2">
                <textarea
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows={3}
                />
                <div className="flex gap-2">
                  <button onClick={handleUpdate} className="text-xs text-blue-500 hover:text-blue-700">
                    保存
                  </button>
                  <button onClick={handleEditCancel} className="text-xs text-gray-500 hover:text-gray-700">
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
              </>
            )}

            <div className="flex items-center gap-3 mt-2">
              <button
                onClick={handleLike}
                className={`text-xs flex items-center gap-1 ${
                  tweet.liked_by_me ? 'text-pink-500' : 'text-gray-400 hover:text-pink-500'
                }`}
              >
                {tweet.liked_by_me ? '♥' : '♡'} {tweet.like_count}
              </button>

              {tweet.user_id === userId && !isEditing && (
                <>
                  <button onClick={handleEditStart} className="text-xs text-blue-500 hover:text-blue-700">
                    編集
                  </button>
                  <button onClick={handleDelete} className="text-xs text-red-500 hover:text-red-700">
                    削除
                  </button>
                </>
              )}
            </div>
          </div>
        ) : (
          <p className="text-sm text-gray-500">読み込み中...</p>
        )}
      </div>
    </div>
  )
}
