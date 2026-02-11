import { auth } from "@/lib/firebase/config";

export const API_BASE_URL = "http://localhost:8080";

/**
 * 認証ヘッダーを含むヘッダーオブジェクトを作成します
 * ユーザーがログインしている場合は、Firebase ID Tokenを含めます
 */
export async function getAuthHeaders(): Promise<Record<string, string>> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  // 現在のユーザーがログインしている場合、ID Tokenを取得
  const user = auth.currentUser;
  if (user) {
    try {
      const token = await user.getIdToken();
      headers["Authorization"] = `Bearer ${token}`;
    } catch (error) {
      console.error("Failed to get ID token:", error);
      // トークン取得失敗時はヘッダーなしで続行
      // （バックエンドが401を返す）
    }
  }

  return headers;
}
