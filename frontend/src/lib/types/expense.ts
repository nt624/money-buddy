export type CreateExpenseInput = {
    amount: number
    category_id: number
    memo?: string
    spent_at: string // yyyy-mm-dd
    status?: 'planned' | 'confirmed'
}

export type UpdateExpenseInput = {
    amount: number
    category_id: number
    memo?: string
    spent_at: string // yyyy-mm-dd or date-time
    status?: 'planned' | 'confirmed'
}

export type Expense = {
    id: number;
    amount: number;
    memo: string | null;
    spent_at: string; // YYYY-MM-DD
    status: 'planned' | 'confirmed';
    category: {
        id: number;
        name: string;
    }
};

export type GetExpensesResponse = {
    expenses: Expense[];
};