export type Category = {
  id: number
  name: string
  category_type: 'default' | 'user' | 'other'
  sort_order: number
}

export type CreateCategoryInput = {
  name: string
}

export type UpdateCategoryInput = {
  name: string
}

export type ReorderCategoryItem = {
  id: number
  sort_order: number
}