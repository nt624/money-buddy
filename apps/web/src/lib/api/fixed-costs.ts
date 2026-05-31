import {
  FixedCost,
  CreateFixedCostInput,
  UpdateFixedCostInput,
  CreateFixedCostResponse,
  GetFixedCostsResponse,
  UpdateFixedCostResponse,
} from "@/lib/types/fixed-cost";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

export async function createFixedCost(
  input: CreateFixedCostInput
): Promise<FixedCost> {
  const headers = await getAuthHeaders(true);
  const res = await fetch(`${API_BASE_URL}/fixed-costs`, {
    method: "POST",
    headers,
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    await handleApiError(res, '固定費の作成');
  }

  const data: CreateFixedCostResponse = await res.json();
  if (!data || !data.fixed_cost) {
    console.error('固定費のレスポンスが正しくありません', { data });
    throw new Error('固定費のレスポンスが正しくありません');
  }

  return data.fixed_cost;
}

export async function getFixedCosts(): Promise<GetFixedCostsResponse> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/fixed-costs`, {
    method: "GET",
    headers,
  });

  if (!res.ok) {
    await handleApiError(res, '固定費の取得');
  }

  const data = await res.json();
  if (!data || !Array.isArray(data.fixed_costs)) {
    console.error('固定費のレスポンスが正しくありません', { data });
    throw new Error('固定費のレスポンスが正しくありません');
  }

  return data;
}

export async function updateFixedCost(
  id: number,
  input: UpdateFixedCostInput
): Promise<FixedCost> {
  const headers = await getAuthHeaders(true);
  const res = await fetch(`${API_BASE_URL}/fixed-costs/${id}`, {
    method: "PUT",
    headers,
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    await handleApiError(res, '固定費の更新');
  }

  const data: UpdateFixedCostResponse = await res.json();
  if (!data || !data.fixed_cost) {
    console.error('固定費のレスポンスが正しくありません', { data });
    throw new Error('固定費のレスポンスが正しくありません');
  }

  return data.fixed_cost;
}

export async function deleteFixedCost(id: number): Promise<void> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/fixed-costs/${id}`, {
    method: "DELETE",
    headers,
  });

  if (!res.ok) {
    await handleApiError(res, '固定費の削除');
  }
}
