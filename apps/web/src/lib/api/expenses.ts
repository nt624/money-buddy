// Thin web facade over the shared API client (see @/lib/apiClient).
// Implementations live in @pace/core/api.
import { api } from "@/lib/apiClient";

export const createExpense = api.expenses.createExpense;
export const getExpenses = api.expenses.getExpenses;
export const updateExpense = api.expenses.updateExpense;
export const deleteExpense = api.expenses.deleteExpense;
