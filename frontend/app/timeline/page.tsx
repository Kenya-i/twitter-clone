'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuth } from '../context/AuthContext'

type Tweet = {
  id: string
  user_id: string
  content: string
  created_at: string
  updated_at: string
}

export default function Timeline() {
  const router = useRouter()
  const { token, userId, logout } = useAuth()
  const [content, setContent] = useState('')
  const [message, setMessage] = useState('')
  const [tweets, setTweets] = useState<Tweet[]>([])

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])

  const fetchTweets = async () => {
    const res = await fetch('http://localhost:8080/tweets', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.ok) {
      const data = await res.json()
      setTweets(data)
    }
  }

  useEffect(() => {
    if (token) fetchTweets()
  }, [token])

  const handlePost = async (e: React.FormEvent) => {
    e.preventDefault()
    setMessage('')

    const res = await fetch('http://localhost:8080/tweets', {
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
    const res = await fetch(`http://localhost:8080/tweets/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.ok) {
      fetchTweets()
    }
  }


  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-md mx-auto bg-white p-6 rounded-lg shadow-md">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">タイムライン</h1>
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
              <div>
                <p className="text-sm">{tweet.content}</p>
                <p className="text-xs text-gray-400 mt-1">
                  {new Date(tweet.created_at).toLocaleString()}
                </p>
              </div>
              {tweet.user_id === userId && (
                <button
                  onClick={() => handleDelete(tweet.id)}
                  className="text-xs text-red-500 hover:text-red-700 ml-2"
                >
                  削除
                </button>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
