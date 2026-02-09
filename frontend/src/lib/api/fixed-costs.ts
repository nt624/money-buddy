import {
  FixedCost,
  UpdateFixedCostInput,
  GetFixedCostsResponse,
  UpdateFixedCostResponse,
} from "@/lib/types/fixed-cost";

const API_BASE_URL = "http://localhost:8080";

export async function getFixedCosts(): Promise<GetFixedCostsResponse> {
  const res = await fetch(`${API_BASE_URL}/fixed-costs`, {
    method: "GET",
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Failed to fetch fixed costs: ${res.status} ${text}`);
  }

  const data = await res.json();
  if (!data || !Array.isArray(data.fixed_costs)) {
    throw new Error(
      `Invalid getFixedCosts response: ${JSON.stringify(data)}`
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
    throw new Error(`Failed to update fixed cost: ${res.status} ${text}`);
  }

  const data: UpdateFixedCostResponse = await res.json();
  if (!data || !data.fixed_cost) {
    throw new Error(
      `Invalid updateFixedCost response: ${JSON.stringify(data)}`
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
    throw new Error(`Failed to delete fixed cost: ${res.status} ${text}`);
  }
}
