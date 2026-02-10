import { InitialSetupRequest, InitialSetupResponse } from "@/lib/types/setup";

const API_BASE_URL = "http://localhost:8080";

export async function submitInitialSetup(
  input: InitialSetupRequest
): Promise<InitialSetupResponse> {
  const res = await fetch(`${API_BASE_URL}/setup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`初期設定の送信に失敗しました: ${res.status} ${text}`);
  }

  const data = await res.json();
  return data;
}
