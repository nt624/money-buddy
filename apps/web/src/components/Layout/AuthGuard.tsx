'use client'

import { useAuth } from '@/contexts/AuthContext'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect } from 'react'

// 認証不要な公開ページのパス
const PUBLIC_PATHS = ['/login']

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth()
  const pathname = usePathname()
  const router = useRouter()

  useEffect(() => {
    // ローディング中は何もしない
    if (loading) return

    // 公開ページの場合は認証チェックをスキップ
    if (PUBLIC_PATHS.includes(pathname)) return

    // 未ログインの場合、ログインページにリダイレクト
    if (!user) {
      // 現在のパスをredirectパラメータとして保存
      const redirectUrl = `/login?redirect=${encodeURIComponent(pathname)}`
      router.push(redirectUrl)
    }
  }, [user, loading, pathname, router])

  // ローディング中
  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">読み込み中...</p>
        </div>
      </div>
    )
  }

  // 公開ページまたはログイン済み
  if (PUBLIC_PATHS.includes(pathname) || user) {
    return <>{children}</>
  }

  // 未ログインでリダイレクト中
  return null
}
