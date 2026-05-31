// Thin web facade over the shared API client (see @/lib/apiClient).
// Implementations live in @pace/core/api.
import { api } from "@/lib/apiClient";

export const getMonthlySettings = api.monthlySettings.getMonthlySettings;
export const upsertMonthlySettings = api.monthlySettings.upsertMonthlySettings;
export const deleteMonthlySettings = api.monthlySettings.deleteMonthlySettings;
