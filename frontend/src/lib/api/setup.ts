import { InitialSetupRequest, InitialSetupResponse } from "@/lib/types/setup";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

export async function submitInitialSetup(
  input: InitialSetupRequest
): Promise<InitialSetupResponse> {
  const headers = await getAuthHeaders(true);
  const res = await fetch(`${API_BASE_URL}/setup`, {
    method: "POST",
    headers,
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    await handleApiError(res, '初期設定の送信');
  }

  const data = await res.json();
  return data;
}
