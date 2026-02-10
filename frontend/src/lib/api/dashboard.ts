import { Dashboard } from "@/lib/types/dashboard";

const API_BASE_URL = "http://localhost:8080";

export async function getDashboard(): Promise<Dashboard> {
  const res = await fetch(`${API_BASE_URL}/dashboard`, {
    method: "GET",
  });

  if (!res.ok) {
    const text = await res.text();
    console.error("ダッシュボード取得に失敗しました", {
      status: res.status,
      body: text,
    });
    throw new Error("ダッシュボードの取得に失敗しました。時間をおいて再度お試しください。");
  }

  const data = await res.json();
  return data;
}
