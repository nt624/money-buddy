import { MonthlySettings, UpsertMonthlySettingsInput } from "@pace/core/types/monthly-settings";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

export async function getMonthlySettings(year: number, month: number): Promise<MonthlySettings> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/monthly-settings?year=${year}&month=${month}`, {
    method: "GET",
    headers,
  });

  if (!res.ok) {
    await handleApiError(res, '月設定の取得');
  }

  return res.json();
}

export async function upsertMonthlySettings(input: UpsertMonthlySettingsInput): Promise<MonthlySettings> {
  const headers = await getAuthHeaders(true);
  const res = await fetch(`${API_BASE_URL}/monthly-settings`, {
    method: "PUT",
    headers,
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    await handleApiError(res, '月設定の保存');
  }

  return res.json();
}

export async function deleteMonthlySettings(year: number, month: number): Promise<void> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/monthly-settings?year=${year}&month=${month}`, {
    method: "DELETE",
    headers,
  });

  if (!res.ok) {
    await handleApiError(res, '月設定の削除');
  }
}
