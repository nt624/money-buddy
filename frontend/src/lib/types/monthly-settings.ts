export type MonthlySettings = {
  year: number
  month: number
  income: number
  saving_goal: number
  is_custom: boolean
}

export type UpsertMonthlySettingsInput = {
  year: number
  month: number
  income: number
  saving_goal: number
}
