export type FixedCost = {
  id: number
  user_id: string
  name: string
  amount: number
  created_at: string
  updated_at: string
}

export type CreateFixedCostInput = {
  name: string
  amount: number
}

export type UpdateFixedCostInput = {
  name: string
  amount: number
}

export type CreateFixedCostResponse = {
  fixed_cost: FixedCost
}

export type GetFixedCostsResponse = {
  fixed_costs: FixedCost[]
}

export type UpdateFixedCostResponse = {
  fixed_cost: FixedCost
}
