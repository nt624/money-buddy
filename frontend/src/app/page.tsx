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

export default function Home() {
  const { user, needsSetup, isLoading: userLoading, error: userError, refetchUser } = useUser()
  const { expenses, createExpense, updateExpense, deleteExpense, isSubmitting, isLoading, error } = useExpenses()
  const { dashboard, isLoading: dashboardLoading, error: dashboardError, refetch: refetchDashboard } = useDashboard()
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
      setSetupError(err instanceof Error ? err.message : 'unknown error')
    } finally {
      setSetupSubmitting(false)
    }
  }

  const handleEdit = (expense: Expense) => {
    setEditingExpense(expense)
  }

  const handleUpdateSubmit = async (input: UpdateExpenseInput) => {
    if (!editingExpense) return
    
    await updateExpense(editingExpense.id, input)
    // 成功時のみ編集モード解除（エラー時はuseExpensesのerror stateに格納される）
    if (!error) {
      setEditingExpense(null)
      refetchDashboard() // 更新後にダッシュボードを再取得
    }
  }

  const handleCancelEdit = () => {
    setEditingExpense(null)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('本当に削除しますか？')) return
    
    await deleteExpense(id)
    // 削除成功時、編集中の支出が削除された場合は編集モードを解除
    if (editingExpense?.id === id) {
      setEditingExpense(null)
    }
    refetchDashboard() // 削除後にダッシュボードを再取得
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

      {/* Dashboard Section */}
      {dashboardLoading && <p>ダッシュボード読み込み中...</p>}
      {dashboardError && <p style={{ color: 'red' }}>ダッシュボードエラー: {dashboardError}</p>}
      {dashboard && <Dashboard dashboard={dashboard} />}

      {editingExpense ? (
        <>
          <h2>支出を編集</h2>
          <ExpenseForm 
            mode="edit"
            initialData={editingExpense}
            onSubmit={handleUpdateSubmit}
            onCancel={handleCancelEdit}
            isSubmitting={isSubmitting}
          />
        </>
      ) : (
        <ExpenseForm 
          onSubmit={async (input) => {
            await createExpense(input)
            refetchDashboard() // 作成後にダッシュボードを再取得
          }}
          isSubmitting={isSubmitting}
        />
      )}

      {isLoading && <p>読み込み中...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <h2>支出一覧</h2>
      <ExpenseList expenses={expenses} onEdit={handleEdit} onDelete={handleDelete} isSubmitting={isSubmitting} />
    </main>
  )
}
