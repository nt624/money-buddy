'use client'

import { useState } from 'react'
import { useExpenses, useUser, useDashboard, useMonthlySettings } from '@pace/core/hooks'
import { ExpenseForm } from '@/components/ExpenseForm'
import { ExpenseList } from '@/components/ExpenseList'
import { ExpenseCalendar } from '@/components/ExpenseCalendar'
import { Button } from '@/components/ui/Button'
import { InitialSetupForm } from '@/components/InitialSetupForm'
import { Dashboard } from '@/components/Dashboard'
import { MonthlySettingsModal } from '@/components/MonthlySettingsModal'
import { Container } from '@/components/Layout/Container'
import { submitInitialSetup } from '@/lib/api/setup'
import { InitialSetupRequest } from '@pace/core/types/setup'
import { Expense, UpdateExpenseInput } from '@pace/core/types/expense'

export default function Home() {
  const { user, needsSetup, isLoading: userLoading, error: userError, refetchUser } = useUser()
  const { expenses, selectedMonth, navigateMonth, createExpense, updateExpense, deleteExpense, isSubmitting, isLoading, error } = useExpenses()
  const { dashboard, isLoading: dashboardLoading, error: dashboardError, refetch: refetchDashboard } = useDashboard({ enabled: !needsSetup, selectedMonth })
  const {
    settings: monthlySettings,
    isLoading: msLoading,
    isSubmitting: msSubmitting,
    error: msError,
    isModalOpen,
    openModal,
    closeModal,
    saveSettings,
    resetSettings,
  } = useMonthlySettings(selectedMonth, user)
  const [setupSubmitting, setSetupSubmitting] = useState(false)
  const [setupError, setSetupError] = useState<string | null>(null)
  const [editingExpense, setEditingExpense] = useState<Expense | null>(null)
  const [viewMode, setViewMode] = useState<'list' | 'calendar'>('list')

  const MONTH_NAMES = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
  const monthLabel = `${selectedMonth.year}年${MONTH_NAMES[selectedMonth.month - 1]}`

  const handleSetupSubmit = async (input: InitialSetupRequest) => {
    setSetupSubmitting(true)
    setSetupError(null)

    try {
      await submitInitialSetup(input)
      await refetchUser()
      refetchDashboard()
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
      <Container className="py-12">
        <p className="text-center text-muted-foreground">読み込み中...</p>
      </Container>
    )
  }

  // ユーザー情報取得エラー（404以外）
  if (userError) {
    return (
      <Container className="py-12">
        <p className="text-danger text-center">エラー: {userError}</p>
      </Container>
    )
  }

  // 初期セットアップが必要
  if (needsSetup) {
    return (
      <div className="min-h-screen py-12 px-4">
        <InitialSetupForm onSubmit={handleSetupSubmit} isSubmitting={setupSubmitting} />
        {setupError && <p className="text-danger text-center mt-4">エラー: {setupError}</p>}
      </div>
    )
  }

  // 通常の支出入力画面
  return (
    <Container className="py-6 space-y-6" maxWidth="lg">
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
          <div className="flex items-center justify-between gap-2 flex-wrap">
            <div className="flex items-center gap-2">
              <Button variant="ghost" size="sm" aria-label="前月へ" onClick={() => navigateMonth('prev')}>‹</Button>
              <span className="text-base font-semibold text-foreground min-w-[120px] text-center">{monthLabel}</span>
              <Button variant="ghost" size="sm" aria-label="翌月へ" onClick={() => navigateMonth('next')}>›</Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={openModal}
                disabled={msLoading}
                className="text-xs text-muted-foreground border border-border"
              >
                {monthlySettings?.is_custom ? '設定済' : '設定'}
              </Button>
            </div>
            <div className="flex gap-1">
              <Button
                variant={viewMode === 'list' ? 'primary' : 'secondary'}
                size="sm"
                onClick={() => setViewMode('list')}
              >
                リスト
              </Button>
              <Button
                variant={viewMode === 'calendar' ? 'primary' : 'secondary'}
                size="sm"
                onClick={() => setViewMode('calendar')}
              >
                カレンダー
              </Button>
            </div>
          </div>
          {viewMode === 'list' ? (
            <ExpenseList expenses={expenses} onEdit={handleEdit} onDelete={handleDelete} isSubmitting={isSubmitting} />
          ) : (
            <ExpenseCalendar expenses={expenses} year={selectedMonth.year} month={selectedMonth.month} />
          )}
        </div>
      </section>

      <MonthlySettingsModal
        isOpen={isModalOpen}
        onClose={closeModal}
        settings={monthlySettings}
        year={selectedMonth.year}
        month={selectedMonth.month}
        onSave={async (input) => {
          const success = await saveSettings(input, () => refetchDashboard())
          return success
        }}
        onReset={async () => {
          const success = await resetSettings(() => refetchDashboard())
          return success
        }}
        isSubmitting={msSubmitting}
        error={msError}
      />
    </Container>
  )
}
