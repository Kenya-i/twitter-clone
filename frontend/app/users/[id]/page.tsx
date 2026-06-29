'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '../../context/AuthContext'
import { API_URL } from '../../../lib/api'

type User = {
  id: string
  username: string
  email: string
  display_name: string
  bio: string
  avatar_url: string
  created_at: string
}

type FollowInfo = {
  followers_count: number
  following_count: number
  is_following: boolean
}

export default function UserProfile() {
  const params = useParams<{ id: string }>()
  const router = useRouter()
  const { token, userId } = useAuth()
  const [user, setUser] = useState<User | null>(null)
  const [followInfo, setFollowInfo] = useState<FollowInfo | null>(null)
  const [avatarFile, setAvatarFile] = useState<File | null>(null)
  const [uploading, setUploading] = useState(false)

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])

  useEffect(() => {
    if (!token) return

    const fetchProfile = async () => {
      const [userRes, followRes] = await Promise.all([
        fetch(`${API_URL}/users/${params.id}`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch(`${API_URL}/users/${params.id}/follow`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ])

      if (userRes.ok) {
        const data = await userRes.json()
        setUser(data)
      }

      if (followRes.ok) {
        const data = await followRes.json()
        setFollowInfo(data)
      }
    }

    fetchProfile()
  }, [token, params.id])

  const handleFollow = async () => {
    if (!followInfo) return
    const wasFollowing = followInfo.is_following

    setFollowInfo({
      ...followInfo,
      is_following: !wasFollowing,
      followers_count: followInfo.followers_count + (wasFollowing ? -1 : 1),
    })

    const res = await fetch(`${API_URL}/users/${params.id}/follow`, {
      method: wasFollowing ? 'DELETE' : 'POST',
      headers: { Authorization: `Bearer ${token}` },
    })

    if (!res.ok) {
      setFollowInfo({
        ...followInfo,
        is_following: wasFollowing,
        followers_count: followInfo.followers_count + (wasFollowing ? 1 : -1),
      })
    }
  }

  const handleAvatarUpload = async () => {
    if (!avatarFile) return
    setUploading(true)

    const formData = new FormData()
    formData.append('avatar', avatarFile)

    const res = await fetch(`${API_URL}/users/me/avatar`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: formData,
    })

    if (res.ok) {
      const data = await res.json()
      setUser(data)
      setAvatarFile(null)
    }

    setUploading(false)
  }


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
            <div className="flex justify-between items-start">
              <img
                src={user.avatar_url || '/default-avatar.svg'}
                alt="avatar"
                className="w-16 h-16 rounded-full object-cover bg-gray-200"
              />
              <div>
                <p className="text-xs text-gray-400">表示名</p>
                <p className="text-lg font-bold">{user.display_name}</p>
                <p className="text-sm text-gray-500">@{user.username}</p>
              </div>
              {params.id !== userId && followInfo && (
                <button
                  onClick={handleFollow}
                  className={`text-sm px-4 py-1 rounded-full border transition-colors ${
                    followInfo.is_following
                      ? 'border-gray-400 text-gray-600 hover:border-red-400 hover:text-red-500'
                      : 'bg-blue-500 text-white border-blue-500 hover:bg-blue-600'
                  }`}
                >
                  {followInfo.is_following ? 'フォロー中' : 'フォロー'}
                </button>
              )}
            </div>

            {followInfo && (
              <div className="flex gap-4 text-sm">
                <span><strong>{followInfo.following_count}</strong> フォロー中</span>
                <span><strong>{followInfo.followers_count}</strong> フォロワー</span>
              </div>
            )}

            {params.id === userId && (
              <div>
                <p className="text-xs text-gray-400">メールアドレス</p>
                <p>{user.email}</p>
              </div>
            )}

            {user.bio && (
              <div>
                <p className="text-xs text-gray-400">自己紹介</p>
                <p>{user.bio}</p>
              </div>
            )}

            {params.id === userId && (
              <div className="space-y-1">
                <input
                  type="file"
                  accept="image/*"
                  onChange={(e) => setAvatarFile(e.target.files?.[0] ?? null)}
                  className="text-xs"
                />
                <button
                  onClick={handleAvatarUpload}
                  disabled={!avatarFile || uploading}
                  className="text-xs text-blue-500 hover:text-blue-700 disabled:text-gray-400"
                >
                  {uploading ? 'アップロード中...' : '画像を更新'}
                </button>
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
