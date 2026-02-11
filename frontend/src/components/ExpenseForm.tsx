'use client'

import { useEffect, useState } from 'react'
import { CreateExpenseInput, UpdateExpenseInput, Expense } from '@/lib/types/expense'
import { getCategories } from '@/lib/api/categories'
import { Category } from '@/lib/types/category'

type Props = {
    mode?: 'create' | 'edit'
    initialData?: Expense
    onSubmit: (input: CreateExpenseInput | UpdateExpenseInput) => Promise<void>
    onCancel?: () => void
    isSubmitting: boolean
}

export function ExpenseForm({ mode = 'create', initialData, onSubmit, onCancel, isSubmitting }: Props) {
    const [amount, setAmount] = useState(initialData?.amount.toString() || '')
    const [categoryId, setCategoryId] = useState(initialData?.category.id.toString() || '1')
    const [memo, setMemo] = useState(initialData?.memo || '')
    const [spentAt, setSpentAt] = useState(initialData?.spent_at || '')
    const [status, setStatus] = useState<'planned' | 'confirmed'>(initialData?.status || 'confirmed')
    const [categories, setCategories] = useState<Category[]>([])

    useEffect(() => {
        getCategories()
            .then(setCategories)
            .catch(console.error)
    }, [])

    // initialDataが変更されたらフォーム状態を更新（編集対象切り替え時）
    useEffect(() => {
        if (initialData) {
            setAmount(initialData.amount.toString())
            setCategoryId(initialData.category.id.toString())
            setMemo(initialData.memo || '')
            setSpentAt(initialData.spent_at)
            setStatus(initialData.status)
            setErrors({})
        }
    }, [initialData])

    const [errors, setErrors] = useState<{
        amount?: string
        spent_at?: string
        status?: string
    }>({})

    const validate = (): boolean => {
        const newErrors: typeof errors = {}

        const amountNumber = Number(amount)
        if (!amount || isNaN(amountNumber) || amountNumber <= 0) {
            newErrors.amount = '金額は0より大きい数値で入力してください'
        }

        if (!spentAt) {
            newErrors.spent_at = '日付を入力してください'
        }

        // 編集モードでstatus遷移ルールをチェック
        if (mode === 'edit' && initialData) {
            if (initialData.status === 'confirmed' && status === 'planned') {
                newErrors.status = '確定済みの支出を予定に戻すことはできません'
            }
        }

        setErrors(newErrors)
        return Object.keys(newErrors).length === 0
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        if (!validate()) return

        await onSubmit({
            amount: Number(amount),
            category_id: Number(categoryId),
            memo: memo || undefined,
            spent_at: spentAt,
            status: status,
        })

        // 作成モードの時のみリセット（statusは保持）
        if (mode === 'create') {
            setAmount('')
            setCategoryId('1')
            setMemo('')
            setSpentAt('')
            setErrors({})
        }
    }

    return (
        <form onSubmit={handleSubmit} className="w-full max-w-lg mx-auto p-3 sm:p-4 bg-card rounded-lg border border-border shadow-sm space-y-3 sm:space-y-4">
            <div className="space-y-2">
                <label className="block text-xs sm:text-sm font-medium text-foreground">
                    金額
                    <input
                        className="mt-1 block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
                        type="number"
                        value={amount}
                        onChange={(e) => setAmount(e.target.value)}
                    />
                </label>
                {errors.amount && <p className="text-xs sm:text-sm text-danger">{errors.amount}</p>}
            </div>

            <div className="space-y-2">
                <label className="block text-xs sm:text-sm font-medium text-foreground">
                    カテゴリ
                    <select
                        className="mt-1 block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
                        value={categoryId}
                        onChange={(e) => setCategoryId(e.target.value)}
                    >
                        <option value="">カテゴリを選択</option>
                        {categories.map((category) => (
                            <option key={category.id} value={category.id}>
                                {category.name}
                            </option>
                        ))}
                    </select>
                </label>
            </div>

            <div className="space-y-2">
                <label className="block text-xs sm:text-sm font-medium text-foreground">
                    メモ
                    <input
                        className="mt-1 block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
                        type="text"
                        value={memo}
                        onChange={(e) => setMemo(e.target.value)}
                    />
                </label>
            </div>

            <div className="space-y-2">
                <label className="block text-xs sm:text-sm font-medium text-foreground">
                    日付
                    <input
                        className="mt-1 block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
                        type="date"
                        value={spentAt}
                        onChange={(e) => setSpentAt(e.target.value)}
                    />
                </label>
                {errors.spent_at && (
                    <p className="text-xs sm:text-sm text-danger">{errors.spent_at}</p>
                )}
            </div>

            <div className="space-y-2">
                <label className="block text-xs sm:text-sm font-medium text-foreground">ステータス</label>
                <div className="flex gap-3 sm:gap-4">
                    <label className={`
                        flex items-center justify-center px-4 py-2.5 rounded-lg border-2 cursor-pointer
                        transition-all duration-200 flex-1 sm:flex-initial
                        ${status === 'confirmed' 
                            ? 'border-primary bg-primary/20 text-primary font-semibold shadow-md' 
                            : 'border-border bg-card hover:border-primary/50 hover:bg-primary/5 hover:shadow-sm text-foreground'}
                    `}>
                        <input
                            className="sr-only"
                            type="radio"
                            name="status"
                            value="confirmed"
                            checked={status === 'confirmed'}
                            onChange={(e) => setStatus(e.target.value as 'confirmed')}
                        />
                        <span className="text-xs sm:text-sm">確定</span>
                    </label>
                    <label className={`
                        flex items-center justify-center px-4 py-2.5 rounded-lg border-2 cursor-pointer
                        transition-all duration-200 flex-1 sm:flex-initial
                        ${status === 'planned' 
                            ? 'border-warning bg-warning/20 text-warning font-semibold shadow-md' 
                            : 'border-border bg-card hover:border-warning/50 hover:bg-warning/5 hover:shadow-sm text-foreground'}
                        ${mode === 'edit' && initialData?.status === 'confirmed' 
                            ? 'opacity-50 cursor-not-allowed' 
                            : ''}
                    `}>
                        <input
                            className="sr-only"
                            type="radio"
                            name="status"
                            value="planned"
                            checked={status === 'planned'}
                            onChange={(e) => setStatus(e.target.value as 'planned')}
                            disabled={mode === 'edit' && initialData?.status === 'confirmed'}
                        />
                        <span className="text-xs sm:text-sm">予定</span>
                    </label>
                </div>
                {errors.status && <p className="text-xs sm:text-sm text-danger">{errors.status}</p>}
            </div>

            <div className="flex flex-col sm:flex-row gap-2">
                <button 
                    className="flex-1 px-4 py-2 text-sm font-medium rounded bg-primary hover:bg-primary-hover text-primary-foreground disabled:opacity-60 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-ring" 
                    type="submit" 
                    disabled={isSubmitting}
                >
                    {isSubmitting ? '送信中...' : mode === 'edit' ? '更新' : '追加'}
                </button>
                {mode === 'edit' && onCancel && (
                    <button 
                        type="button" 
                        onClick={onCancel}
                        className="px-4 py-2 text-sm font-medium rounded bg-secondary hover:bg-secondary-hover text-secondary-foreground transition-colors focus:outline-none focus:ring-2 focus:ring-ring"
                    >
                        キャンセル
                    </button>
                )}
            </div>
        </form>
    )
}