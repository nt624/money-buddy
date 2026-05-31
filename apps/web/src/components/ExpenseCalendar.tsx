import { Expense } from "@pace/core/types/expense";

type Props = {
  expenses: Expense[];
  year: number;
  month: number;
};

const WEEKDAY_LABELS = ['日', '月', '火', '水', '木', '金', '土'];

function getExpensesForDay(expenses: Expense[], year: number, month: number, day: number): Expense[] {
  const pad = (n: number) => String(n).padStart(2, '0');
  const prefix = `${year}-${pad(month)}-${pad(day)}`;
  return expenses.filter((e) => e.spent_at.startsWith(prefix));
}

export function ExpenseCalendar({ expenses, year, month }: Props) {
  const firstDay = new Date(year, month - 1, 1).getDay();
  const daysInMonth = new Date(year, month, 0).getDate();
  const today = new Date();

  const cells: (number | null)[] = [
    ...Array(firstDay).fill(null),
    ...Array.from({ length: daysInMonth }, (_, i) => i + 1),
  ];
  while (cells.length % 7 !== 0) cells.push(null);

  return (
    <div className="bg-card border border-border rounded-lg shadow-sm p-4">
      <div className="grid grid-cols-7 mb-2">
        {WEEKDAY_LABELS.map((label) => (
          <div key={label} className="text-center text-xs font-semibold text-muted-foreground py-1">
            {label}
          </div>
        ))}
      </div>
      <div className="grid grid-cols-7 gap-px bg-border">
        {cells.map((day, idx) => {
          if (!day) return <div key={idx} className="bg-card min-h-[80px]" />;

          const dayExpenses = getExpensesForDay(expenses, year, month, day);
          const totalAmount = dayExpenses.reduce((sum, e) => sum + e.amount, 0);
          const isToday =
            today.getFullYear() === year &&
            today.getMonth() + 1 === month &&
            today.getDate() === day;

          return (
            <div key={idx} className="bg-card min-h-[80px] p-1">
              <div className={`text-xs font-medium mb-1 ${isToday ? 'text-primary font-bold' : 'text-foreground'}`}>
                {day}
              </div>
              {dayExpenses.length > 0 && (
                <div className="space-y-0.5">
                  <div className="text-xs font-semibold text-foreground">
                    ¥{totalAmount.toLocaleString()}
                  </div>
                  {dayExpenses.slice(0, 2).map((e) => (
                    <div
                      key={e.id}
                      className={`text-xs truncate px-1 rounded ${
                        e.status === 'confirmed'
                          ? 'bg-success text-success-foreground'
                          : 'bg-warning text-warning-foreground'
                      }`}
                    >
                      {e.category.name}
                    </div>
                  ))}
                  {dayExpenses.length > 2 && (
                    <div className="text-xs text-muted-foreground">+{dayExpenses.length - 2}件</div>
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
