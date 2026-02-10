import {
  FixedCost,
  CreateFixedCostInput,
  UpdateFixedCostInput,
  CreateFixedCostResponse,
  GetFixedCostsResponse,
  UpdateFixedCostResponse,
} from "@/lib/types/fixed-cost";

const API_BASE_URL = "http://localhost:8080";

export async function createFixedCost(
  input: CreateFixedCostInput
): Promise<FixedCost> {
  const res = await fetch(`${API_BASE_URL}/fixed-costs`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`固定費の作成に失敗しました: ${res.status} ${text}`);
  }

  const data: CreateFixedCostResponse = await res.json();
  if (!data || !data.fixed_cost) {
    throw new Error(
      `固定費のレスポンスが正しくありません: ${JSON.stringify(data)}`
    );
  }

  return data.fixed_cost;
}

export async function getFixedCosts(): Promise<GetFixedCostsResponse> {
  const res = await fetch(`${API_BASE_URL}/fixed-costs`, {
    method: "GET",
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`固定費の取得に失敗しました: ${res.status} ${text}`);
  }

  const data = await res.json();
  if (!data || !Array.isArray(data.fixed_costs)) {
    throw new Error(
      `固定費のレスポンスが正しくありません: ${JSON.stringify(data)}`
    );
  }

  return data;
}

export async function updateFixedCost(
  id: number,
  input: UpdateFixedCostInput
): Promise<FixedCost> {
  const res = await fetch(`${API_BASE_URL}/fixed-costs/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`固定費の更新に失敗しました: ${res.status} ${text}`);
  }

  const data: UpdateFixedCostResponse = await res.json();
  if (!data || !data.fixed_cost) {
    throw new Error(
      `固定費のレスポンスが正しくありません: ${JSON.stringify(data)}`
    );
  }

  return data.fixed_cost;
}

export async function deleteFixedCost(id: number): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/fixed-costs/${id}`, {
    method: "DELETE",
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`固定費の削除に失敗しました: ${res.status} ${text}`);
  }
}
