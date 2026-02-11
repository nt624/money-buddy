'use client'

import { useState } from 'react'
import { useFixedCosts } from '@/hooks/useFixedCosts'
import { useDashboard } from '@/hooks/useDashboard'
import { FixedCostList } from '@/components/FixedCostList'
import { FixedCostForm } from '@/components/FixedCostForm'
import { UserForm } from '@/components/UserForm'
import { FixedCost, FixedCostInput } from '@/lib/types/fixed-cost'
import { UpdateUserInput } from '@/lib/types/user'
import { Container } from '@/components/Layout/Container'
import { Button } from '@/components/ui/Button'

export default function SettingsPage() {
  const { dashboard, updateUserSettings, isSubmitting: userSubmitting, isLoading: dashboardLoading, error: dashboardError, refetch: refetchDashboard } = useDashboard()
  const { fixedCosts, createFixedCost, updateFixedCost, deleteFixedCost, isSubmitting: fcSubmitting, isLoading: fcLoading, error: fcError } = useFixedCosts()
  const [editingFixedCost, setEditingFixedCost] = useState<FixedCost | null>(null)
  const [showFixedCostForm, setShowFixedCostForm] = useState(false)
  const [showUserForm, setShowUserForm] = useState(false)

  const handleFixedCostEdit = (fixedCost: FixedCost) => {
    setEditingFixedCost(fixedCost)
    setShowFixedCostForm(true)
  }

  const handleFixedCostCreateSubmit = async (input: FixedCostInput) => {
    const success = await createFixedCost(input)
    // 成功時のみフォーム閉じてダッシュボード再取得
    if (success) {
      setShowFixedCostForm(false)
      refetchDashboard()
    }
  }

  const handleFixedCostUpdateSubmit = async (input: FixedCostInput) => {
    if (!editingFixedCost) return
    
    const success = await updateFixedCost(editingFixedCost.id, input)
    // 成功時のみ編集モード解除とダッシュボード再取得
    if (success) {
      setEditingFixedCost(null)
      setShowFixedCostForm(false)
      refetchDashboard()
    }
  }

  const handleFixedCostSubmit = async (input: FixedCostInput) => {
    if (editingFixedCost) {
      await handleFixedCostUpdateSubmit(input)
    } else {
      await handleFixedCostCreateSubmit(input)
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

  const handleUserSubmit = async (input: UpdateUserInput) => {
    const success = await updateUserSettings(input)
    if (success) {
      setShowUserForm(false)
    }
  }

  const handleCancelUserEdit = () => {
    setShowUserForm(false)
  }

  return (
    <Container className="py-6 space-y-6" maxWidth="lg">
      <h1 className="text-2xl sm:text-3xl font-bold text-foreground">設定</h1>

      {/* 基本情報セクション */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold text-foreground">基本情報</h2>
          {!showUserForm && (
            <Button 
              onClick={() => setShowUserForm(true)}
              disabled={dashboardLoading}
              size="sm"
            >
              編集
            </Button>
          )}
        </div>

        {dashboardLoading && <p className="text-muted-foreground">読み込み中...</p>}

        {!dashboardLoading && !showUserForm && dashboard && (
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <div className="text-sm text-muted-foreground">月収（手取り）</div>
              <div className="text-xl font-semibold text-foreground">¥{dashboard.income.toLocaleString()}</div>
            </div>
            <div className="space-y-1">
              <div className="text-sm text-muted-foreground">貯金目標額（月）</div>
              <div className="text-xl font-semibold text-foreground">¥{dashboard.saving_goal.toLocaleString()}</div>
            </div>
          </div>
        )}

        {showUserForm && dashboard && (
          <UserForm
            initialIncome={dashboard.income}
            initialSavingGoal={dashboard.saving_goal}
            onSubmit={handleUserSubmit}
            onCancel={handleCancelUserEdit}
            isSubmitting={userSubmitting}
            error={dashboardError}
          />
        )}
      </section>

      {/* 固定費セクション */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold text-foreground">固定費</h2>
          {!showFixedCostForm && (
            <Button 
              onClick={() => setShowFixedCostForm(true)}
              size="sm"
            >
              固定費を追加
            </Button>
          )}
        </div>
        
        {showFixedCostForm && (
          <div className="mb-4">
            <FixedCostForm 
              fixedCost={editingFixedCost}
              onSubmit={handleFixedCostSubmit}
              onCancel={handleCancelFixedCostEdit}
              isSubmitting={fcSubmitting}
            />
          </div>
        )}

        {fcLoading && <p className="text-muted-foreground">読み込み中...</p>}
        {fcError && <p className="text-danger">{fcError}</p>}

        <FixedCostList 
          fixedCosts={fixedCosts} 
          onEdit={handleFixedCostEdit} 
          onDelete={handleFixedCostDelete} 
          isSubmitting={fcSubmitting} 
        />
      </section>
    </Container>
  )
}
