export type FixedCost = {
  id: number
  user_id: string
  name: string
  amount: number
  created_at: string
  updated_at: string
}

// フォーム入力用の共通型
export type FixedCostInput = {
  name: string
  amount: number
}

// API型定義（互換性のため残す）
export type CreateFixedCostInput = FixedCostInput
export type UpdateFixedCostInput = FixedCostInput

export type CreateFixedCostResponse = {
  fixed_cost: FixedCost
}

export type GetFixedCostsResponse = {
  fixed_costs: FixedCost[]
}

export type UpdateFixedCostResponse = {
  fixed_cost: FixedCost
}
