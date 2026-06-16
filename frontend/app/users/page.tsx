'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '../context/AuthContext'

type User = {
  id: string
  username: string
  display_name: string
  bio: string
}

export default function UserList() {
  const router = useRouter()
  const { token, userId } = useAuth()
  const [users, setUsers] = useState<User[]>([])

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])

  useEffect(() => {
    if (!token) return

    const fetchUsers = async () => {
      const res = await fetch('http://localhost:8080/users', {
        headers: { Authorization: `Bearer ${token}` },
      })

      if (res.ok) {
        const data = await res.json()
        setUsers(data.filter((u: User) => u.id !== userId))
      }
    }

    fetchUsers()
  }, [token, userId])

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-md mx-auto bg-white p-6 rounded-lg shadow-md">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">ユーザー一覧</h1>
          <Link href="/timeline" className="text-sm text-blue-500 hover:underline">
            タイムラインに戻る
          </Link>
        </div>

        <div className="space-y-3">
          {users.map((user) => (
            <Link
              key={user.id}
              href={`/users/${user.id}`}
              className="block border-b border-gray-200 pb-3 hover:bg-gray-50 rounded px-2 transition-colors"
            >
              <p className="font-bold">{user.display_name}</p>
              <p className="text-sm text-gray-500">@{user.username}</p>
              {user.bio && <p className="text-xs text-gray-400 mt-1">{user.bio}</p>}
            </Link>
          ))}
          {users.length === 0 && (
            <p className="text-sm text-gray-400 text-center">他のユーザーがいません</p>
          )}
        </div>
      </div>
    </div>
  )
}
