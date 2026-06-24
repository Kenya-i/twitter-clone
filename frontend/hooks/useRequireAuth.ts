'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

export function useRequireAuth(token: string | null) {
  const router = useRouter()

  useEffect(() => {
    if (token === null) {
      const saved = localStorage.getItem('token')
      if (!saved) router.push('/')
    }
  }, [token, router])
}
