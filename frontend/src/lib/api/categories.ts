import { Category } from "@/lib/types/category";

export async function getCategories(): Promise<Category[]> {
  const res = await fetch("http://localhost:8080/categories")

  if (!res.ok) {
    throw new Error("Failed to fetch categories")
  }

  const data = await res.json()
  return data.categories
}