'use client'

import { useEffect, useState } from 'react'
import styles from './ExpenseForm.module.css'
import { CreateExpenseInput } from '@/lib/types/expense'
import { getCategories } from '@/lib/api/categories'
import { Category } from '@/lib/types/category'

type Props = {
    onSubmit: (input: CreateExpenseInput) => Promise<void>
    isSubmitting: boolean
}

export function ExpenseForm({ onSubmit, isSubmitting }: Props) {
    const [amount, setAmount] = useState('')
    const [categoryId, setCategoryId] = useState('1')
    const [memo, setMemo] = useState('')
    const [spentAt, setSpentAt] = useState('')
    const [categories, setCategories] = useState<Category[]>([])

    useEffect(() => {
        getCategories()
            .then(setCategories)
            .catch(console.error)
    }, [])

    const [errors, setErrors] = useState<{
        amount?: string
        spent_at?: string
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
        })

        // 成功したらリセット
        setAmount('')
        setCategoryId('1')
        setMemo('')
        setSpentAt('')
        setErrors({})
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

            <button className={styles.button} type="submit" disabled={isSubmitting}>
                {isSubmitting ? '送信中...' : '追加'}
            </button>
        </form>
    )
}