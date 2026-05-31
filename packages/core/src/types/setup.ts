export type FixedCostInput = {
  name: string
  amount: number
}

export type InitialSetupRequest = {
  income: number
  savingGoal: number
  fixedCosts: FixedCostInput[]
}

export type InitialSetupResponse = {
  status: string
}
