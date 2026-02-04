import { Dashboard } from "@/lib/types/dashboard";

const API_BASE_URL = "http://localhost:8080";

export async function getDashboard(): Promise<Dashboard> {
  const res = await fetch(`${API_BASE_URL}/dashboard`, {
    method: "GET",
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Failed to fetch dashboard: ${res.status} ${text}`);
  }

  const data = await res.json();
  return data;
}
