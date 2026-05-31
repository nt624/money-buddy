import type { DataSource } from "@pace/core/data";
import { api } from "@/lib/apiClient";

/**
 * REST-backed DataSource for the web app. Reads/writes go straight to the
 * backend through the shared API client (no local store). Mobile will provide
 * a local-first implementation backed by packages/sync.
 */
export const restDataSource: DataSource = {
  expenses: {
    list: (period) => api.expenses.getExpenses(period).then((r) => r.expenses),
    create: (input) => api.expenses.createExpense(input),
    update: (id, input) => api.expenses.updateExpense(id, input),
    remove: (id) => api.expenses.deleteExpense(id),
  },
  categories: {
    list: () => api.categories.getCategories(),
    create: (input) => api.categories.createCategory(input),
    update: (id, input) => api.categories.updateCategory(id, input),
    remove: (id) => api.categories.deleteCategory(id),
    reorder: (items) => api.categories.reorderCategories(items),
  },
  dashboard: {
    get: (period) => api.dashboard.getDashboard(period),
  },
  user: {
    getMe: () => api.user.getMe(),
    update: (input) => api.user.updateUser(input),
  },
  fixedCosts: {
    list: () => api.fixedCosts.getFixedCosts().then((r) => r.fixed_costs),
    create: (input) => api.fixedCosts.createFixedCost(input),
    update: (id, input) => api.fixedCosts.updateFixedCost(id, input),
    remove: (id) => api.fixedCosts.deleteFixedCost(id),
  },
  monthlySettings: {
    get: (year, month) => api.monthlySettings.getMonthlySettings(year, month),
    upsert: (input) => api.monthlySettings.upsertMonthlySettings(input),
    remove: (year, month) => api.monthlySettings.deleteMonthlySettings(year, month),
  },
};
