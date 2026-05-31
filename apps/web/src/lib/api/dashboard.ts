import { Dashboard } from "@/lib/types/dashboard";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

export async function getDashboard(params?: { year: number; month: number }): Promise<Dashboard> {
  const headers = await getAuthHeaders();
  const query = params ? `?year=${params.year}&month=${params.month}` : '';
  const res = await fetch(`${API_BASE_URL}/dashboard${query}`, {
    method: "GET",
    headers,
  });

  if (!res.ok) {
    await handleApiError(res, 'ダッシュボードの取得');
  }

  const data = await res.json();
  return data;
}
