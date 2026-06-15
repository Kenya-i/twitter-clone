'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '../context/AuthContext'

type User = {
  id: string
  username: string
  email: string
  display_name: string
  bio: string
  created_at: string
}

export default function Profile() {
  const router = useRouter()
  const { token, userId } = useAuth()
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])

  useEffect(() => {
    if (!token || !userId) return

    const fetchProfile = async () => {
      const res = await fetch(`http://localhost:8080/users/${userId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (res.ok) {
        const data = await res.json()
        setUser(data)
      }
    }

    fetchProfile()
  }, [token, userId])

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-md mx-auto bg-white p-6 rounded-lg shadow-md">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">プロフィール</h1>
          <Link href="/timeline" className="text-sm text-blue-500 hover:underline">
            タイムラインに戻る
          </Link>
        </div>

        {user ? (
          <div className="space-y-3">
            <div>
              <p className="text-xs text-gray-400">表示名</p>
              <p className="text-lg font-bold">{user.display_name}</p>
            </div>
            <div>
              <p className="text-xs text-gray-400">ユーザー名</p>
              <p>@{user.username}</p>
            </div>
            <div>
              <p className="text-xs text-gray-400">メールアドレス</p>
              <p>{user.email}</p>
            </div>
            {user.bio && (
              <div>
                <p className="text-xs text-gray-400">自己紹介</p>
                <p>{user.bio}</p>
              </div>
            )}
            <div>
              <p className="text-xs text-gray-400">登録日</p>
              <p>{new Date(user.created_at).toLocaleDateString()}</p>
            </div>
          </div>
        ) : (
          <p className="text-sm text-gray-500">読み込み中...</p>
        )}
      </div>
    </div>
  )
}
