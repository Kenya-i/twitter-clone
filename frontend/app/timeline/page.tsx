'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuth } from '../context/AuthContext'

export default function Timeline() {
  const router = useRouter()
  const { token, logout } = useAuth()
  const [content, setContent] = useState('')
  const [message, setMessage] = useState('')

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])

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
      </div>
    </div>
  )
}
