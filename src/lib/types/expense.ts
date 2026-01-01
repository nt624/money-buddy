export type CreateExpenseInput = {
    amount: number
    category_id: number
    memo?: string
    spent_at: string // yyyy-mm-dd
}

export type Expense = {
    id: number;
    amount: number;
    category_id: number;
    memo: string | null;
    spent_at: string; // YYYY-MM-DD
};

export type GetExpensesResponse = {
    expenses: Expense[];
};