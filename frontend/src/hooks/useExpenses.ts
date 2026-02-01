import { useEffect, useState } from "react";
import { createExpense, getExpenses } from "@/lib/api/expenses";
import { CreateExpenseInput, Expense } from "@/lib/types/expense";

export function useExpenses() {
    const [expenses, setExpenses] = useState<Expense[]>([]);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // GET
    useEffect(() => {
        const fetchExpenses = async () => {
            setIsLoading(true);
            setError(null);

            try {
                const data = await getExpenses();
                setExpenses(data.expenses);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'unknown error');
            } finally {
                setIsLoading(false);
            }
        }

        fetchExpenses();
    }, []);

    // POST
    const addExpense = async (input: CreateExpenseInput) => {
        setIsSubmitting(true);
        setError(null);

        try {
            const expense = await createExpense(input);
            setExpenses((prevExpenses) => [...prevExpenses, expense]);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'unknown error');
        } finally {
            setIsSubmitting(false);
        }
    }

    return {
        expenses,
        isLoading,
        isSubmitting,
        error,
        addExpense,
    }
}
