import { InitialSetupRequest, InitialSetupResponse } from "@/lib/types/setup";
import { API_BASE_URL, getAuthHeaders } from "./client";

export async function submitInitialSetup(
  input: InitialSetupRequest
): Promise<InitialSetupResponse> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/setup`, {
    method: "POST",
    headers,
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    const text = await res.text();
    console.error("submitInitialSetup failed", {
      status: res.status,
      body: text,
    });
    throw new Error("初期設定の送信に失敗しました。時間をおいて再度お試しください。");
  }

  const data = await res.json();
  return data;
}
