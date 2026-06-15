'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '../context/AuthContext'

export default function Profile() {
  const router = useRouter()
  const { token, userId } = useAuth()

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) {
        router.push('/')
        return
      }
    }

    if (userId) {
      router.replace(`/users/${userId}`)
    }
  }, [token, userId, router])

  return <p className="text-center mt-8 text-sm text-gray-500">読み込み中...</p>
}