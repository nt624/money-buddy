'use client'

import { useState } from 'react'
import { MonthlySettings, UpsertMonthlySettingsInput } from '@/lib/types/monthly-settings'
import { BUSINESS_MAX_AMOUNT } from '@/lib/constants'

type MonthlySettingsModalProps = {
  isOpen: boolean
  onClose: () => void
  settings: MonthlySettings | null
  year: number
  month: number
  onSave: (input: UpsertMonthlySettingsInput) => Promise<boolean>
  onReset: () => Promise<boolean>
  isSubmitting: boolean
  error: string | null
}

const MONTH_NAMES = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']

export function MonthlySettingsModal({
  isOpen,
  onClose,
  settings,
  year,
  month,
  onSave,
  onReset,
  isSubmitting,
  error,
}: MonthlySettingsModalProps) {
  const [prevSettings, setPrevSettings] = useState(settings)
  const [prevIsOpen, setPrevIsOpen] = useState(isOpen)
  const [income, setIncome] = useState(settings?.income ?? 0)
  const [savingGoal, setSavingGoal] = useState(settings?.saving_goal ?? 0)

  // Sync form values when settings or open state changes (React-recommended pattern)
  if (prevSettings !== settings || prevIsOpen !== isOpen) {
    setPrevSettings(settings)
    setPrevIsOpen(isOpen)
    if (settings) {
      setIncome(settings.income)
      setSavingGoal(settings.saving_goal)
    }
  }

  if (!isOpen) return null

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const success = await onSave({ year, month, income, saving_goal: savingGoal })
    if (success) onClose()
  }

  const handleReset = async () => {
    if (!confirm('この月の設定を削除してデフォルト値に戻しますか？')) return
    const success = await onReset()
    if (success) onClose()
  }

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-full items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />
      <div className="relative z-10 bg-card border border-border rounded-xl shadow-xl w-full max-w-sm p-6">
        <h2 className="text-lg font-semibold text-foreground mb-1">
          {year}年{MONTH_NAMES[month - 1]}の設定
        </h2>
        {settings && !settings.is_custom && (
          <p className="text-xs text-muted-foreground mb-4">
            現在はデフォルト値を表示しています
          </p>
        )}
        {settings?.is_custom && (
          <p className="text-xs text-primary mb-4">
            この月専用の設定が適用されています
          </p>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="p-3 bg-danger/10 border border-danger rounded-lg text-danger text-sm">
              {error}
            </div>
          )}

          <div className="space-y-2">
            <label htmlFor="ms-income" className="block text-sm font-medium text-foreground">
              手取り（月収）
            </label>
            <input
              type="number"
              id="ms-income"
              value={income}
              onChange={(e) => setIncome(Number(e.target.value))}
              disabled={isSubmitting}
              min="1"
              max={BUSINESS_MAX_AMOUNT}
              required
              className="w-full px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="ms-saving-goal" className="block text-sm font-medium text-foreground">
              目標貯金額
            </label>
            <input
              type="number"
              id="ms-saving-goal"
              value={savingGoal}
              onChange={(e) => setSavingGoal(Number(e.target.value))}
              disabled={isSubmitting}
              min="0"
              max={BUSINESS_MAX_AMOUNT}
              required
              className="w-full px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>

          <div className="flex gap-2 pt-2">
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 px-4 py-2 text-sm font-medium rounded-lg bg-primary hover:bg-primary-hover text-primary-foreground disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isSubmitting ? '保存中...' : '保存'}
            </button>
            <button
              type="button"
              onClick={onClose}
              disabled={isSubmitting}
              className="px-4 py-2 text-sm font-medium rounded-lg bg-secondary hover:bg-secondary-hover text-secondary-foreground disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              キャンセル
            </button>
          </div>
        </form>

        {settings?.is_custom && (
          <div className="mt-3 pt-3 border-t border-border">
            <button
              type="button"
              onClick={handleReset}
              disabled={isSubmitting}
              className="w-full px-4 py-2 text-sm font-medium rounded-lg text-muted-foreground hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              デフォルトに戻す
            </button>
          </div>
        )}
      </div>
      </div>
    </div>
  )
}
