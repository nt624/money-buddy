'use client'

import { useState } from 'react'
import { useExpenses } from '@/hooks/useExpenses'
import { useUser } from '@/hooks/useUser'
import { ExpenseForm } from '@/components/ExpenseForm'
import { ExpenseList } from '@/components/ExpenseList'
import { InitialSetupForm } from '@/components/InitialSetupForm'
import { submitInitialSetup } from '@/lib/api/setup'
import { InitialSetupRequest } from '@/lib/types/setup'

export default function Home() {
  const { user, needsSetup, isLoading: userLoading, error: userError, refetchUser } = useUser()
  const { expenses, addExpense, isSubmitting, isLoading, error } = useExpenses()
  const [setupSubmitting, setSetupSubmitting] = useState(false)
  const [setupError, setSetupError] = useState<string | null>(null)

  const handleSetupSubmit = async (input: InitialSetupRequest) => {
    setSetupSubmitting(true)
    setSetupError(null)

    try {
      await submitInitialSetup(input)
      await refetchUser() // セットアップ完了後、ユーザー情報を再取得
    } catch (err) {
      setSetupError(err instanceof Error ? err.message : 'unknown error')
      throw err
    } finally {
      setSetupSubmitting(false)
    }
  }

  // 初回読み込み中
  if (userLoading) {
    return (
      <main>
        <p>読み込み中...</p>
      </main>
    )
  }

  // ユーザー情報取得エラー（404以外）
  if (userError) {
    return (
      <main>
        <p style={{ color: 'red' }}>エラー: {userError}</p>
      </main>
    )
  }

  // 初期セットアップが必要
  if (needsSetup) {
    return (
      <main>
        <InitialSetupForm onSubmit={handleSetupSubmit} isSubmitting={setupSubmitting} />
        {setupError && <p style={{ color: 'red' }}>エラー: {setupError}</p>}
      </main>
    )
  }

  // 通常の支出入力画面
  return (
    <main>
      <h1>支出入力</h1>
      {user && (
        <div style={{ marginBottom: '1rem', padding: '1rem', backgroundColor: '#f5f5f5', borderRadius: '4px' }}>
          <p>月収: ¥{user.income.toLocaleString()}</p>
          <p>貯蓄目標: ¥{user.saving_goal.toLocaleString()}</p>
        </div>
      )}

      <ExpenseForm onSubmit={addExpense} isSubmitting={isSubmitting} />

      {isLoading && <p>読み込み中...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <h2>支出一覧</h2>
      <ExpenseList expenses={expenses} />
    </main>
  )
}
