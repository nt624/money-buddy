export type User = {
  id: string // Firebase UID
  income: number
  saving_goal: number
  created_at: string
  updated_at: string
}

export type UpdateUserInput = {
  income: number
  saving_goal: number
}
