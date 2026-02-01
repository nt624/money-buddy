import { Expense } from "@/lib/types/expense"

type Props = {
  expenses: Expense[]
}

export function ExpenseList({ expenses }: Props) {
  return (
    <ul>
      {expenses.map((e) => (
        <li key={e.id}>
          <span>{e.spent_at}</span> / 
          <span>¥{e.amount}</span> / 
          <span>カテゴリ: {e.category.name}</span> / 
          <span>{e.memo}</span>
        </li>
      ))}
    </ul>
  )
}
