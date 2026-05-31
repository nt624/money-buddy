import { useEffect, useState } from "react";
import { createExpense, getExpenses, updateExpense, deleteExpense } from "@/lib/api/expenses";
import { CreateExpenseInput, UpdateExpenseInput, Expense } from "@pace/core/types/expense";

type SelectedMonth = { year: number; month: number };

function currentMonth(): SelectedMonth {
    const now = new Date();
    return { year: now.getFullYear(), month: now.getMonth() + 1 };
}

export function useExpenses() {
    const [selectedMonth, setSelectedMonth] = useState<SelectedMonth>(currentMonth);
    const [expenses, setExpenses] = useState<Expense[]>([]);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // GET
    useEffect(() => {
        const year = selectedMonth.year;
        const month = selectedMonth.month;
        const fetchExpenses = async () => {
            setIsLoading(true);
            setError(null);

            try {
                const data = await getExpenses({ year, month });
                setExpenses(data.expenses);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'エラーが発生しました');
            } finally {
                setIsLoading(false);
            }
        }

        fetchExpenses();
    }, [selectedMonth.year, selectedMonth.month]);

    const navigateMonth = (direction: 'prev' | 'next') => {
        setSelectedMonth(({ year, month }) => {
            if (direction === 'prev') {
                return month === 1 ? { year: year - 1, month: 12 } : { year, month: month - 1 };
            }
            return month === 12 ? { year: year + 1, month: 1 } : { year, month: month + 1 };
        });
    };

    // POST
    const handleCreateExpense = async (input: CreateExpenseInput): Promise<boolean> => {
        setIsSubmitting(true);
        setError(null);

        try {
            const expense = await createExpense(input);
            setExpenses((prevExpenses) => [...prevExpenses, expense]);
            return true;
        } catch (err) {
            setError(err instanceof Error ? err.message : 'エラーが発生しました');
            return false;
        } finally {
            setIsSubmitting(false);
        }
    }

    // PUT
    const handleUpdateExpense = async (id: number, input: UpdateExpenseInput): Promise<boolean> => {
        setIsSubmitting(true);
        setError(null);

        try {
            const updatedExpense = await updateExpense(id, input);
            setExpenses((prevExpenses) =>
                prevExpenses.map((exp) => (exp.id === id ? updatedExpense : exp))
            );
            return true;
        } catch (err) {
            setError(err instanceof Error ? err.message : 'エラーが発生しました');
            return false;
        } finally {
            setIsSubmitting(false);
        }
    }

    // DELETE
    const handleDeleteExpense = async (id: number): Promise<boolean> => {
        setIsSubmitting(true);
        setError(null);

        try {
            await deleteExpense(id);
            setExpenses((prevExpenses) => prevExpenses.filter((exp) => exp.id !== id));
            return true;
        } catch (err) {
            setError(err instanceof Error ? err.message : 'エラーが発生しました');
            return false;
        } finally {
            setIsSubmitting(false);
        }
    }

    return {
        expenses,
        selectedMonth,
        navigateMonth,
        isLoading,
        isSubmitting,
        error,
        createExpense: handleCreateExpense,
        updateExpense: handleUpdateExpense,
        deleteExpense: handleDeleteExpense,
    }
}
