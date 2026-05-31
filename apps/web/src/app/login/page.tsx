'use client'

export const dynamic = 'force-dynamic'

import { useState, useEffect, useCallback } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useAuth } from '@/contexts/AuthContext'
import { Container } from '@/components/Layout/Container'

type FirebaseError = { code: string; message: string }

function isFirebaseError(err: unknown): err is FirebaseError {
  return (
    typeof err === 'object' &&
    err !== null &&
    'code' in err &&
    typeof (err as Record<string, unknown>).code === 'string'
  )
}

function getFirebaseErrorMessage(err: unknown): string {
  if (!isFirebaseError(err)) return 'エラーが発生しました'
  const code = err.code
  if (code === 'auth/email-already-in-use') return 'このメールアドレスは既に使用されています'
  if (code === 'auth/invalid-email') return 'メールアドレスの形式が正しくありません'
  if (code === 'auth/weak-password') return 'パスワードは6文字以上にしてください'
  if (code === 'auth/user-not-found' || code === 'auth/wrong-password') return 'メールアドレスまたはパスワードが間違っています'
  if (code === 'auth/invalid-credential') return 'メールアドレスまたはパスワードが間違っています'
  if (code === 'auth/popup-blocked') return 'ポップアップがブロックされました。ブラウザの設定を確認してください'
  if (code === 'auth/popup-closed-by-user') return 'ログインがキャンセルされました'
  if (code === 'auth/unauthorized-domain') return 'このドメインは許可されていません'
  if (code === 'auth/operation-not-allowed') return 'Google認証が有効化されていません'
  return 'エラーが発生しました'
}

export default function LoginPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { user, signIn, signUp, signInWithGoogle } = useAuth()
  const [isSignUp, setIsSignUp] = useState(false)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState(() =>
    searchParams.get('reason') === 'session_expired'
      ? '認証の有効期限が切れました。再度ログインしてください。'
      : ''
  )
  const [loading, setLoading] = useState(false)

  /**
   * リダイレクト先を安全に取得する
   * オープンリダイレクト脆弱性対策: アプリ内パス（/で始まる）のみ許可
   */
  const getSafeRedirect = useCallback((): string => {
    const redirect = searchParams.get('redirect')
    if (!redirect) return '/'
    if (redirect.startsWith('/') && !redirect.startsWith('//')) return redirect
    console.warn(`Invalid redirect parameter detected: ${redirect}`)
    return '/'
  }, [searchParams])

  useEffect(() => {
    // 既にログイン済みの場合、リダイレクト先またはホームへ
    if (user) {
      router.push(getSafeRedirect())
    }
  }, [user, router, getSafeRedirect])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      if (isSignUp) {
        await signUp(email, password)
      } else {
        await signIn(email, password)
      }
      router.push(getSafeRedirect())
    } catch (err: unknown) {
      setError(getFirebaseErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleGoogleSignIn = async () => {
    setError('')
    setLoading(true)

    try {
      await signInWithGoogle()
      router.push(getSafeRedirect())
    } catch (err: unknown) {
      console.error('Googleログインエラー:', err)
      if (isFirebaseError(err)) {
        setError(getFirebaseErrorMessage(err))
      } else {
        setError('Googleログインに失敗しました')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Container className="py-12 flex items-center justify-center min-h-[calc(100vh-3.5rem)]">
      <div className="max-w-md w-full bg-card p-6 sm:p-8 rounded-lg border border-border shadow-lg">
        <h1 className="text-2xl sm:text-3xl font-bold mb-6 text-center text-foreground">
          {isSignUp ? '新規登録' : 'ログイン'}
        </h1>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-foreground mb-1">
              メールアドレス
            </label>
            <input
              id="email"
              type="email"
              placeholder="example@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
              required
            />
          </div>
          
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-foreground mb-1">
              パスワード
            </label>
            <input
              id="password"
              type="password"
              placeholder="6文字以上"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
              required
              minLength={6}
            />
          </div>
          
          {error && (
            <div className="p-3 bg-danger/10 border border-danger rounded-md">
              <p className="text-danger text-sm">{error}</p>
            </div>
          )}
          
          <button 
            type="submit"
            disabled={loading}
            className="w-full bg-primary hover:bg-primary-hover text-primary-foreground py-2.5 rounded-md font-medium transition-colors disabled:opacity-60 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-ring"
          >
            {loading ? '処理中...' : isSignUp ? '登録' : 'ログイン'}
          </button>
        </form>

        <div className="relative my-6">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-border"></div>
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-card text-muted-foreground">または</span>
          </div>
        </div>

        <button
          onClick={handleGoogleSignIn}
          disabled={loading}
          className="w-full bg-background hover:bg-muted border border-border py-2.5 rounded-md font-medium transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2 focus:outline-none focus:ring-2 focus:ring-ring"
        >
          <svg className="w-5 h-5" viewBox="0 0 24 24">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
          {isSignUp ? 'Googleで登録' : 'Googleでログイン'}
        </button>

        <p className="mt-6 text-center text-sm text-muted-foreground">
          {isSignUp ? 'アカウントをお持ちですか？' : 'アカウントをお持ちでない方'}
          <button
            type="button"
            onClick={() => {
              setIsSignUp(!isSignUp)
              setError('')
            }}
            className="text-primary hover:underline ml-2 font-medium"
          >
            {isSignUp ? 'ログイン' : '新規登録'}
          </button>
        </p>
      </div>
    </Container>
  )
}
