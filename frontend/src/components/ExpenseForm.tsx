'use client'

import { useEffect, useState } from 'react'
import styles from './ExpenseForm.module.css'
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
        <form onSubmit={handleSubmit} className={styles.form}>
            <div className={styles.field}>
                <label className={styles.label}>
                    金額
                    <input
                        className={styles.input}
                        type="number"
                        value={amount}
                        onChange={(e) => setAmount(e.target.value)}
                    />
                </label>
                {errors.amount && <p className={styles.error}>{errors.amount}</p>}
            </div>

            <div className={styles.field}>
                <label className={styles.label}>
                    カテゴリ
                    <select
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

            <div className={styles.field}>
                <label className={styles.label}>
                    メモ
                    <input
                        className={styles.input}
                        type="text"
                        value={memo}
                        onChange={(e) => setMemo(e.target.value)}
                    />
                </label>
            </div>

            <div className={styles.field}>
                <label className={styles.label}>
                    日付
                    <input
                        className={styles.input}
                        type="date"
                        value={spentAt}
                        onChange={(e) => setSpentAt(e.target.value)}
                    />
                </label>
                {errors.spent_at && (
                    <p className={styles.error}>{errors.spent_at}</p>
                )}
            </div>

            <div className={styles.field}>
                <label className={styles.label}>ステータス</label>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <input
                            type="radio"
                            name="status"
                            value="confirmed"
                            checked={status === 'confirmed'}
                            onChange={(e) => setStatus(e.target.value as 'confirmed')}
                        />
                        確定
                    </label>
                    <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <input
                            type="radio"
                            name="status"
                            value="planned"
                            checked={status === 'planned'}
                            onChange={(e) => setStatus(e.target.value as 'planned')}
                            disabled={mode === 'edit' && initialData?.status === 'confirmed'}
                        />
                        予定
                    </label>
                </div>
                {errors.status && <p className={styles.error}>{errors.status}</p>}
            </div>

            <div style={{ display: 'flex', gap: '0.5rem' }}>
                <button className={styles.button} type="submit" disabled={isSubmitting}>
                    {isSubmitting ? '送信中...' : mode === 'edit' ? '更新' : '追加'}
                </button>
                {mode === 'edit' && onCancel && (
                    <button 
                        type="button" 
                        onClick={onCancel}
                        style={{
                            padding: '0.5rem 1rem',
                            backgroundColor: '#6c757d',
                            color: 'white',
                            border: 'none',
                            borderRadius: '4px',
                            cursor: 'pointer'
                        }}
                    >
                        キャンセル
                    </button>
                )}
            </div>
        </form>
    )
}