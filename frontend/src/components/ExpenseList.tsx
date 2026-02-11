import { Expense } from "@/lib/types/expense"

type Props = {
  expenses: Expense[]
  onEdit: (expense: Expense) => void
  onDelete: (id: number) => void
  isSubmitting: boolean
}

export function ExpenseList({ expenses, onEdit, onDelete, isSubmitting }: Props) {
  return (
    <ul className="list-none p-0 space-y-4">
      {expenses.map((e) => (
        <li 
          key={e.id} 
          className="p-4 border border-border bg-card rounded-lg flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 shadow-sm hover:shadow-md transition-shadow"
        >
          <div className="flex-1 space-y-2">
            <div>
              <span className={`
                inline-block px-2 py-1 rounded text-xs font-bold
                ${e.status === 'confirmed' 
                  ? 'bg-success text-success-foreground' 
                  : 'bg-warning text-warning-foreground'}
              `}>
                {e.status === 'confirmed' ? '確定' : '予定'}
              </span>
            </div>
            <div className="text-sm text-foreground space-x-2">
              <span className="font-medium">{e.spent_at}</span>
              <span className="text-muted-foreground">/</span>
              <span className="font-semibold">¥{e.amount.toLocaleString()}</span>
              <span className="text-muted-foreground">/</span>
              <span className="text-muted-foreground">カテゴリ: {e.category.name}</span>
              {e.memo && (
                <>
                  <span className="text-muted-foreground">/</span>
                  <span className="text-sm">{e.memo}</span>
                </>
              )}
            </div>
          </div>
          <div className="flex gap-2 sm:flex-shrink-0">
            <button
              onClick={() => onEdit(e)}
              disabled={isSubmitting}
              className="
                px-4 py-2 text-sm font-medium rounded
                bg-primary hover:bg-primary-hover text-primary-foreground
                disabled:opacity-60 disabled:cursor-not-allowed
                transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
              "
            >
              編集
            </button>
            <button
              onClick={() => onDelete(e.id)}
              disabled={isSubmitting}
              className="
                px-4 py-2 text-sm font-medium rounded
                bg-danger hover:bg-danger-hover text-danger-foreground
                disabled:opacity-60 disabled:cursor-not-allowed
                transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
              "
            >
              削除
            </button>
          </div>
        </li>
      ))}
    </ul>
  )
}

