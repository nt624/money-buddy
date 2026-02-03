import { useEffect, useState } from "react";
import { createExpense, getExpenses, updateExpense, deleteExpense } from "@/lib/api/expenses";
import { CreateExpenseInput, UpdateExpenseInput, Expense } from "@/lib/types/expense";

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
    const handleCreateExpense = async (input: CreateExpenseInput) => {
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

    // PUT
    const handleUpdateExpense = async (id: number, input: UpdateExpenseInput) => {
        setIsSubmitting(true);
        setError(null);

        try {
            const updatedExpense = await updateExpense(id, input);
            setExpenses((prevExpenses) =>
                prevExpenses.map((exp) => (exp.id === id ? updatedExpense : exp))
            );
        } catch (err) {
            setError(err instanceof Error ? err.message : 'unknown error');
        } finally {
            setIsSubmitting(false);
        }
    }

    // DELETE
    const handleDeleteExpense = async (id: number) => {
        setIsSubmitting(true);
        setError(null);

        try {
            await deleteExpense(id);
            setExpenses((prevExpenses) => prevExpenses.filter((exp) => exp.id !== id));
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
        createExpense: handleCreateExpense,
        updateExpense: handleUpdateExpense,
        deleteExpense: handleDeleteExpense,
    }
}
