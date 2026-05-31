// Thin web facade over the shared API client (see @/lib/apiClient).
// Implementations live in @pace/core/api.
import { api } from "@/lib/apiClient";

export const getCategories = api.categories.getCategories;
export const createCategory = api.categories.createCategory;
export const updateCategory = api.categories.updateCategory;
export const deleteCategory = api.categories.deleteCategory;
export const reorderCategories = api.categories.reorderCategories;
