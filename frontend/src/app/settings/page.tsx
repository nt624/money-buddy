'use client'

import { useState } from 'react'
import { useFixedCosts } from '@/hooks/useFixedCosts'
import { useDashboard } from '@/hooks/useDashboard'
import { FixedCostList } from '@/components/FixedCostList'
import { FixedCostForm } from '@/components/FixedCostForm'
import { UserForm } from '@/components/UserForm'
import { FixedCost, FixedCostInput } from '@/lib/types/fixed-cost'
import { UpdateUserInput } from '@/lib/types/user'
import Link from 'next/link'

export default function SettingsPage() {
  const { dashboard, updateUserSettings, isSubmitting: userSubmitting, isLoading: dashboardLoading, error: dashboardError } = useDashboard()
  const { fixedCosts, createFixedCost, updateFixedCost, deleteFixedCost, isSubmitting: fcSubmitting, isLoading: fcLoading, error: fcError } = useFixedCosts()
  const { refetch: refetchDashboard } = useDashboard({ enabled: false })
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
    <main style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ marginBottom: '2rem' }}>
        <Link href="/" style={{ color: '#3b82f6', textDecoration: 'none' }}>
          ← ホームに戻る
        </Link>
      </div>

      <h1 style={{ marginBottom: '2rem' }}>設定</h1>

      {/* 基本情報セクション */}
      <section style={{ marginBottom: '3rem', padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', border: '1px solid #e5e7eb' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h2 style={{ fontSize: '1.25rem', fontWeight: '600' }}>基本情報</h2>
          {!showUserForm && (
            <button 
              onClick={() => setShowUserForm(true)}
              disabled={dashboardLoading}
              style={{ 
                padding: '0.5rem 1rem', 
                backgroundColor: dashboardLoading ? '#9ca3af' : '#3b82f6', 
                color: 'white', 
                border: 'none', 
                borderRadius: '6px', 
                cursor: dashboardLoading ? 'not-allowed' : 'pointer' 
              }}
            >
              編集
            </button>
          )}
        </div>

        {dashboardLoading && <p>読み込み中...</p>}

        {!dashboardLoading && !showUserForm && dashboard && (
          <div style={{ color: '#374151' }}>
            <div style={{ marginBottom: '0.5rem' }}>
              <span style={{ fontWeight: 500 }}>月収（手取り）:</span> {dashboard.income.toLocaleString()}円
            </div>
            <div>
              <span style={{ fontWeight: 500 }}>貯金目標額（月）:</span> {dashboard.saving_goal.toLocaleString()}円
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
      <section style={{ marginBottom: '3rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h2 style={{ fontSize: '1.25rem', fontWeight: '600' }}>固定費</h2>
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
              onSubmit={handleFixedCostSubmit}
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
    </main>
  )
}
