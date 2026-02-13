import { auth } from "@/lib/firebase/config";
import { signOut } from "firebase/auth";

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

/**
 * 認証ヘッダーを含むヘッダーオブジェクトを作成します
 * ユーザーがログインしている場合は、Firebase ID Tokenを含めます
 * @param hasBody リクエストボディがある場合true（Content-Typeヘッダーを付与）
 */
export async function getAuthHeaders(hasBody: boolean = false): Promise<Record<string, string>> {
  const headers: Record<string, string> = {};

  // リクエストボディがある場合のみContent-Typeを付与
  if (hasBody) {
    headers["Content-Type"] = "application/json";
  }

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

/**
 * APIレスポンスのエラーハンドリング
 * 401エラーの場合は自動的にログアウトします
 */
export async function handleApiError(response: Response, operation: string): Promise<never> {
  // 401 Unauthorized - トークンが無効または期限切れ
  if (response.status === 401) {
    console.warn("認証エラーが発生しました。ログアウトします。");
    
    // 自動ログアウト
    try {
      await signOut(auth);
    } catch (error) {
      console.error("ログアウトに失敗しました:", error);
    }
    
    // ログインページにリダイレクト
    if (typeof window !== "undefined") {
      window.location.href = "/login?reason=session_expired";
    }
    
    throw new Error("認証の有効期限が切れました。再度ログインしてください。");
  }

  // その他のエラー
  let errorMessage = `${operation}に失敗しました`;
  try {
    const errorData = await response.json();
    if (errorData.error) {
      errorMessage = errorData.error;
    }
  } catch {
    // JSONパースエラーは無視
  }

  throw new Error(errorMessage);
}
