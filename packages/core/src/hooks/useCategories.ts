import { useEffect, useState } from "react";
import { useDataSource } from "../data";
import { Category, CreateCategoryInput, UpdateCategoryInput, ReorderCategoryItem } from "../types/category";

export function useCategories() {
  const ds = useDataSource();
  const [categories, setCategories] = useState<Category[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // GET
  useEffect(() => {
    const fetchCategories = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const data = await ds.categories.list();
        setCategories(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "エラーが発生しました");
      } finally {
        setIsLoading(false);
      }
    };
    fetchCategories();
  }, [ds]);

  const refetchCategories = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await ds.categories.list();
      setCategories(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  };

  // POST
  const handleCreateCategory = async (input: CreateCategoryInput): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);
    try {
      const newCategory = await ds.categories.create(input);
      setCategories((prev) => [...prev, newCategory]);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  // PUT
  const handleUpdateCategory = async (id: number, input: UpdateCategoryInput): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);
    try {
      const updated = await ds.categories.update(id, input);
      setCategories((prev) => prev.map((c) => (c.id === id ? updated : c)));
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  // DELETE
  const handleDeleteCategory = async (id: number): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);
    try {
      await ds.categories.remove(id);
      setCategories((prev) => prev.filter((c) => c.id !== id));
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  // REORDER
  const handleReorderCategories = async (items: ReorderCategoryItem[]): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);
    try {
      await ds.categories.reorder(items);
      await refetchCategories();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    categories,
    isLoading,
    isSubmitting,
    error,
    createCategory: handleCreateCategory,
    updateCategory: handleUpdateCategory,
    deleteCategory: handleDeleteCategory,
    reorderCategories: handleReorderCategories,
    refetch: refetchCategories,
  };
}
