import { Dashboard } from "@/lib/types/dashboard";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

export async function getDashboard(): Promise<Dashboard> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/dashboard`, {
    method: "GET",
    headers,
  });

  if (!res.ok) {
    await handleApiError(res, 'ダッシュボードの取得');
  }

  const data = await res.json();
  return data;
}
