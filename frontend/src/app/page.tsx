'use client'

import { useState } from 'react'
import { useExpenses } from '@/hooks/useExpenses'
import { useFixedCosts } from '@/hooks/useFixedCosts'
import { useUser } from '@/hooks/useUser'
import { useDashboard } from '@/hooks/useDashboard'
import { ExpenseForm } from '@/components/ExpenseForm'
import { ExpenseList } from '@/components/ExpenseList'
import { FixedCostList } from '@/components/FixedCostList'
import { FixedCostForm } from '@/components/FixedCostForm'
import { InitialSetupForm } from '@/components/InitialSetupForm'
import { Dashboard } from '@/components/Dashboard'
import { submitInitialSetup } from '@/lib/api/setup'
import { InitialSetupRequest } from '@/lib/types/setup'
import { Expense, UpdateExpenseInput } from '@/lib/types/expense'
import { FixedCost, UpdateFixedCostInput } from '@/lib/types/fixed-cost'

export default function Home() {
  const { user, needsSetup, isLoading: userLoading, error: userError, refetchUser } = useUser()
  const { expenses, createExpense, updateExpense, deleteExpense, isSubmitting, isLoading, error } = useExpenses()
  const { dashboard, isLoading: dashboardLoading, error: dashboardError, refetch: refetchDashboard } = useDashboard({ enabled: !needsSetup })
  const { fixedCosts, updateFixedCost, deleteFixedCost, isSubmitting: fcSubmitting, isLoading: fcLoading, error: fcError, refetch: refetchFixedCosts } = useFixedCosts()
  const [setupSubmitting, setSetupSubmitting] = useState(false)
  const [setupError, setSetupError] = useState<string | null>(null)
  const [editingExpense, setEditingExpense] = useState<Expense | null>(null)
  const [editingFixedCost, setEditingFixedCost] = useState<FixedCost | null>(null)
  const [showFixedCostForm, setShowFixedCostForm] = useState(false)

  const handleSetupSubmit = async (input: InitialSetupRequest) => {
    setSetupSubmitting(true)
    setSetupError(null)

    try {
      await submitInitialSetup(input)
      await refetchUser() // セットアップ完了後、ユーザー情報を再取得
      refetchDashboard() // ダッシュボードも更新
      await refetchFixedCosts() // 固定費も再取得
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

  const handleFixedCostEdit = (fixedCost: FixedCost) => {
    setEditingFixedCost(fixedCost)
    setShowFixedCostForm(true)
  }

  const handleFixedCostUpdateSubmit = async (input: UpdateFixedCostInput) => {
    if (!editingFixedCost) return
    
    const success = await updateFixedCost(editingFixedCost.id, input)
    // 成功時のみ編集モード解除とダッシュボード再取得
    if (success) {
      setEditingFixedCost(null)
      setShowFixedCostForm(false)
      refetchDashboard()
    }
  }

  const handleCancelFixedCostEdit = () => {
    setEditingFixedCost(null)
    setShowFixedCostForm(false)
  }

  const handleFixedCostDelete = async (id: number) => {
    if (!confirm('本当に削除しますか？')) return
    
    const success = await deleteFixedCost(id)
    // 削除成功時のみダッシュボード再取得と編集モード解除
    if (success) {
      if (editingFixedCost?.id === id) {
        setEditingFixedCost(null)
        setShowFixedCostForm(false)
      }
      refetchDashboard()
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

      {/* Dashboard Section */}
      {dashboardLoading && <p>ダッシュボード読み込み中...</p>}
      {dashboardError && <p style={{ color: 'red' }}>ダッシュボードエラー: {dashboardError}</p>}
      {dashboard && <Dashboard dashboard={dashboard} />}

      {/* 固定費セクション */}
      <section style={{ marginBottom: '2rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h2>固定費</h2>
          {!showFixedCostForm && (
            <button 
              onClick={() => setShowFixedCostForm(true)}
              style={{ padding: '0.5rem 1rem', backgroundColor: '#3b82f6', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer' }}
            >
              固定費を追加
            </button>
          )}
        </div>
        
        {showFixedCostForm && (
          <div style={{ marginBottom: '1rem' }}>
            <FixedCostForm 
              fixedCost={editingFixedCost}
              onSubmit={handleFixedCostUpdateSubmit}
              onCancel={handleCancelFixedCostEdit}
              isSubmitting={fcSubmitting}
            />
          </div>
        )}

        {fcLoading && <p>読み込み中...</p>}
        {fcError && <p style={{ color: 'red' }}>{fcError}</p>}

        <FixedCostList 
          fixedCosts={fixedCosts} 
          onEdit={handleFixedCostEdit} 
          onDelete={handleFixedCostDelete} 
          isSubmitting={fcSubmitting} 
        />
      </section>

      {/* 支出セクション */}
      <section>
        <h2>支出入力</h2>
        {editingExpense ? (
          <>
            <h3>支出を編集</h3>
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
              const success = await createExpense(input)
              if (success) {
                refetchDashboard() // 作成成功時のみダッシュボードを再取得
              }
            }}
            isSubmitting={isSubmitting}
          />
        )}

        {isLoading && <p>読み込み中...</p>}
        {error && <p style={{ color: 'red' }}>{error}</p>}

        <h3>支出一覧</h3>
        <ExpenseList expenses={expenses} onEdit={handleEdit} onDelete={handleDelete} isSubmitting={isSubmitting} />
      </section>
    </main>
  )
}
