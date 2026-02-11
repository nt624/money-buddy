'use client'

import { useState } from 'react'
import { useExpenses } from '@/hooks/useExpenses'
import { useUser } from '@/hooks/useUser'
import { useDashboard } from '@/hooks/useDashboard'
import { ExpenseForm } from '@/components/ExpenseForm'
import { ExpenseList } from '@/components/ExpenseList'
import { InitialSetupForm } from '@/components/InitialSetupForm'
import { Dashboard } from '@/components/Dashboard'
import { submitInitialSetup } from '@/lib/api/setup'
import { InitialSetupRequest } from '@/lib/types/setup'
import { Expense, UpdateExpenseInput } from '@/lib/types/expense'
import Link from 'next/link'

export default function Home() {
  const { user, needsSetup, isLoading: userLoading, error: userError, refetchUser } = useUser()
  const { expenses, createExpense, updateExpense, deleteExpense, isSubmitting, isLoading, error } = useExpenses()
  const { dashboard, isLoading: dashboardLoading, error: dashboardError, refetch: refetchDashboard } = useDashboard({ enabled: !needsSetup })
  const [setupSubmitting, setSetupSubmitting] = useState(false)
  const [setupError, setSetupError] = useState<string | null>(null)
  const [editingExpense, setEditingExpense] = useState<Expense | null>(null)

  const handleSetupSubmit = async (input: InitialSetupRequest) => {
    setSetupSubmitting(true)
    setSetupError(null)

    try {
      await submitInitialSetup(input)
      await refetchUser() // セットアップ完了後、ユーザー情報を再取得
      refetchDashboard() // ダッシュボードも更新
    } catch (err) {
      setSetupError(err instanceof Error ? err.message : 'エラーが発生しました')
    } finally {
      setSetupSubmitting(false)
    }
  }

  const handleEdit = (expense: Expense) => {
    setEditingExpense(expense)
  }

  const handleUpdateSubmit = async (input: UpdateExpenseInput) => {
    if (!editingExpense) return
    
    const success = await updateExpense(editingExpense.id, input)
    // 成功時のみ編集モード解除とダッシュボード再取得
    if (success) {
      setEditingExpense(null)
      refetchDashboard()
    }
  }

  const handleCancelEdit = () => {
    setEditingExpense(null)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('本当に削除しますか？')) return
    
    const success = await deleteExpense(id)
    // 削除成功時のみダッシュボード再取得と編集モード解除
    if (success) {
      if (editingExpense?.id === id) {
        setEditingExpense(null)
      }
      refetchDashboard()
    }
  }
  

  // 初回読み込み中
  if (userLoading) {
    return (
      <main className="max-w-4xl mx-auto p-6">
        <p className="text-center text-muted-foreground">読み込み中...</p>
      </main>
    )
  }

  // ユーザー情報取得エラー（404以外）
  if (userError) {
    return (
      <main className="max-w-4xl mx-auto p-6">
        <p className="text-danger text-center">エラー: {userError}</p>
      </main>
    )
  }

  // 初期セットアップが必要
  if (needsSetup) {
    return (
      <main className="min-h-screen py-12 px-4">
        <InitialSetupForm onSubmit={handleSetupSubmit} isSubmitting={setupSubmitting} />
        {setupError && <p className="text-danger text-center mt-4">エラー: {setupError}</p>}
      </main>
    )
  }

  // 通常の支出入力画面
  return (
    <main className="max-w-4xl mx-auto p-4 sm:p-6 space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
        <h1 className="text-2xl sm:text-3xl font-bold text-foreground">ホーム</h1>
        <Link 
          href="/settings" 
          className="inline-flex items-center gap-1 text-primary hover:text-primary-hover text-sm font-medium transition-colors"
        >
          ⚙️ 設定
        </Link>
      </div>
      
      {/* User Info Card */}
      {user && (
        <div className="p-4 bg-card border border-border rounded-lg space-y-2 shadow-sm">
          <p className="text-foreground">
            <span className="font-medium">月収:</span> ¥{user.income.toLocaleString()}
          </p>
          <p className="text-foreground">
            <span className="font-medium">貯蓄目標:</span> ¥{user.saving_goal.toLocaleString()}
          </p>
        </div>
      )}

      {/* Dashboard Section */}
      {dashboardLoading && <p className="text-muted-foreground">ダッシュボード読み込み中...</p>}
      {dashboardError && <p className="text-danger">ダッシュボードエラー: {dashboardError}</p>}
      {dashboard && <Dashboard dashboard={dashboard} />}

      {/* Expense Section */}
      <section className="space-y-4">
        <h2 className="text-xl sm:text-2xl font-bold text-foreground">支出入力</h2>
        {editingExpense ? (
          <div className="space-y-2">
            <h3 className="text-lg font-semibold text-foreground">支出を編集</h3>
            <ExpenseForm 
              mode="edit"
              initialData={editingExpense}
              onSubmit={handleUpdateSubmit}
              onCancel={handleCancelEdit}
              isSubmitting={isSubmitting}
            />
          </div>
        ) : (
          <ExpenseForm 
            onSubmit={async (input) => {
              const success = await createExpense(input)
              if (success) {
                refetchDashboard()
              }
            }}
            isSubmitting={isSubmitting}
          />
        )}

        {isLoading && <p className="text-muted-foreground">読み込み中...</p>}
        {error && <p className="text-danger">{error}</p>}

        <div className="space-y-2">
          <h3 className="text-lg font-semibold text-foreground">支出一覧</h3>
          <ExpenseList expenses={expenses} onEdit={handleEdit} onDelete={handleDelete} isSubmitting={isSubmitting} />
        </div>
      </section>
    </main>
  )
}
