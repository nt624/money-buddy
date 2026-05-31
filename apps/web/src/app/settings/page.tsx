'use client'

import { useState } from 'react'
import { useFixedCosts, useCategories, useDashboard, useUser } from '@pace/core/hooks'
import { FixedCostList } from '@/components/FixedCostList'
import { FixedCostForm } from '@/components/FixedCostForm'
import { CategoryList } from '@/components/CategoryList'
import { CategoryForm } from '@/components/CategoryForm'
import { UserForm } from '@/components/UserForm'
import { FixedCost, FixedCostInput } from '@pace/core/types/fixed-cost'
import { Category, ReorderCategoryItem } from '@pace/core/types/category'
import { UpdateUserInput } from '@pace/core/types/user'
import { updateUser } from '@/lib/api/users'
import { Container } from '@/components/Layout/Container'
import { Button } from '@/components/ui/Button'

export default function SettingsPage() {
  // グローバルデフォルト値の取得・編集には useUser を使う
  // （useDashboard は COALESCE で月別値を返す場合があるためデフォルト表示に不適）
  const { user, refetchUser, isLoading: userLoading } = useUser()
  const { refetch: refetchDashboard } = useDashboard({ enabled: false })
  const { fixedCosts, createFixedCost, updateFixedCost, deleteFixedCost, isSubmitting: fcSubmitting, isLoading: fcLoading, error: fcError } = useFixedCosts()
  const { categories, createCategory, updateCategory, deleteCategory, reorderCategories, isSubmitting: catSubmitting, isLoading: catLoading, error: catError } = useCategories()

  const [editingFixedCost, setEditingFixedCost] = useState<FixedCost | null>(null)
  const [showFixedCostForm, setShowFixedCostForm] = useState(false)
  const [showUserForm, setShowUserForm] = useState(false)
  const [userSubmitting, setUserSubmitting] = useState(false)
  const [userError, setUserError] = useState<string | null>(null)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [showCategoryForm, setShowCategoryForm] = useState(false)

  const handleFixedCostEdit = (fixedCost: FixedCost) => {
    setEditingFixedCost(fixedCost)
    setShowFixedCostForm(true)
  }

  const handleFixedCostCreateSubmit = async (input: FixedCostInput) => {
    const success = await createFixedCost(input)
    if (success) {
      setShowFixedCostForm(false)
      refetchDashboard()
    }
  }

  const handleFixedCostUpdateSubmit = async (input: FixedCostInput) => {
    if (!editingFixedCost) return
    const success = await updateFixedCost(editingFixedCost.id, input)
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
    if (success) {
      if (editingFixedCost?.id === id) {
        setEditingFixedCost(null)
        setShowFixedCostForm(false)
      }
      refetchDashboard()
    }
  }

  const handleUserSubmit = async (input: UpdateUserInput) => {
    setUserSubmitting(true)
    setUserError(null)
    try {
      await updateUser(input)
      await refetchUser()
      refetchDashboard()
      setShowUserForm(false)
    } catch (err) {
      setUserError(err instanceof Error ? err.message : 'ユーザー設定の更新に失敗しました')
    } finally {
      setUserSubmitting(false)
    }
  }

  const handleCancelUserEdit = () => {
    setUserError(null)
    setShowUserForm(false)
  }

  // Category handlers
  const handleCategoryEdit = (category: Category) => {
    setEditingCategory(category)
    setShowCategoryForm(true)
  }

  const handleCategoryFormSubmit = async (name: string) => {
    if (editingCategory) {
      const success = await updateCategory(editingCategory.id, { name })
      if (success) {
        setEditingCategory(null)
        setShowCategoryForm(false)
      }
    } else {
      const success = await createCategory({ name })
      if (success) setShowCategoryForm(false)
    }
  }

  const handleCancelCategoryEdit = () => {
    setEditingCategory(null)
    setShowCategoryForm(false)
  }

  const handleCategoryDelete = async (id: number) => {
    if (!confirm('本当に削除しますか？')) return
    const success = await deleteCategory(id)
    if (success && editingCategory?.id === id) {
      setEditingCategory(null)
      setShowCategoryForm(false)
    }
  }

  const handleCategoryReorder = async (items: ReorderCategoryItem[]) => {
    await reorderCategories(items)
  }

  return (
    <Container className="py-6 space-y-6" maxWidth="lg">
      <h1 className="text-2xl sm:text-3xl font-bold text-foreground">設定</h1>

      {/* 基本情報セクション */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-4 sm:p-6">
        <div className="flex justify-between items-center mb-1 sm:mb-2">
          <h2 className="text-base sm:text-lg font-semibold text-foreground">基本情報</h2>
          {!showUserForm && (
            <Button
              onClick={() => setShowUserForm(true)}
              disabled={userLoading}
              size="sm"
            >
              編集
            </Button>
          )}
        </div>

        <p className="text-xs text-muted-foreground mb-3 sm:mb-4">
          月毎の設定がない月に適用されるデフォルト値です
        </p>

        {userLoading && <p className="text-muted-foreground">読み込み中...</p>}

        {!userLoading && !showUserForm && user && (
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
            <div className="space-y-1">
              <div className="text-xs sm:text-sm text-muted-foreground">月収（手取り）</div>
              <div className="text-lg sm:text-xl font-semibold text-foreground">¥{user.income.toLocaleString()}</div>
            </div>
            <div className="space-y-1">
              <div className="text-xs sm:text-sm text-muted-foreground">貯金目標額（月）</div>
              <div className="text-lg sm:text-xl font-semibold text-foreground">¥{user.saving_goal.toLocaleString()}</div>
            </div>
          </div>
        )}

        {showUserForm && user && (
          <UserForm
            initialIncome={user.income}
            initialSavingGoal={user.saving_goal}
            onSubmit={handleUserSubmit}
            onCancel={handleCancelUserEdit}
            isSubmitting={userSubmitting}
            error={userError}
          />
        )}
      </section>

      {/* カテゴリセクション */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-4 sm:p-6">
        <div className="flex justify-between items-center mb-3 sm:mb-4">
          <h2 className="text-base sm:text-lg font-semibold text-foreground">カテゴリ</h2>
          {!showCategoryForm && (
            <Button
              onClick={() => { setEditingCategory(null); setShowCategoryForm(true) }}
              size="sm"
            >
              カテゴリを追加
            </Button>
          )}
        </div>

        {showCategoryForm && (
          <div className="mb-4">
            <CategoryForm
              category={editingCategory}
              onSubmit={handleCategoryFormSubmit}
              onCancel={handleCancelCategoryEdit}
              isSubmitting={catSubmitting}
            />
          </div>
        )}

        {catLoading && <p className="text-muted-foreground">読み込み中...</p>}
        {catError && <p className="text-danger">{catError}</p>}

        <CategoryList
          categories={categories}
          onEdit={handleCategoryEdit}
          onDelete={handleCategoryDelete}
          onReorder={handleCategoryReorder}
          isSubmitting={catSubmitting}
        />
      </section>

      {/* 固定費セクション */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-4 sm:p-6">
        <div className="flex justify-between items-center mb-3 sm:mb-4">
          <h2 className="text-base sm:text-lg font-semibold text-foreground">固定費</h2>
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
