import { createApiClient } from "@pace/core/api";
import { firebaseWebAuthPort } from "@/auth/firebase-web";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

/**
 * The web app's singleton API client. Token retrieval is delegated to the
 * Firebase Web AuthPort; a 401 signs the user out and redirects to /login.
 */
export const api = createApiClient({
  baseUrl: API_BASE_URL,
  getToken: () => firebaseWebAuthPort.getToken(),
  onUnauthorized: () => {
    firebaseWebAuthPort.signOut().catch((error) => {
      console.error("ログアウトに失敗しました:", error);
    });
    if (typeof window !== "undefined") {
      window.location.href = "/login?reason=session_expired";
    }
  },
});
