'use client'

import { useState, useEffect } from "react";
import { Category } from "@/lib/types/category";

const CATEGORY_NAME_MAX_LENGTH = 50;

type Props = {
  category: Category | null;
  onSubmit: (name: string) => void;
  onCancel: () => void;
  isSubmitting: boolean;
};

export function CategoryForm({ category, onSubmit, onCancel, isSubmitting }: Props) {
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setName(category ? category.name : "");
    setError(null);
  }, [category]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const trimmedName = name.trim();
    if (!trimmedName) {
      setError("カテゴリ名を入力してください");
      return;
    }
    if ([...trimmedName].length > CATEGORY_NAME_MAX_LENGTH) {
      setError(`カテゴリ名は${CATEGORY_NAME_MAX_LENGTH}文字以内で入力してください`);
      return;
    }

    onSubmit(trimmedName);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h3 className="text-base font-semibold text-foreground">
        {category ? "カテゴリを編集" : "カテゴリを追加"}
      </h3>

      {error && (
        <div className="p-3 bg-danger/10 border border-danger rounded-lg text-danger text-sm">
          {error}
        </div>
      )}

      <div className="space-y-2">
        <label htmlFor="category-name" className="block text-sm font-medium text-foreground">
          カテゴリ名 <span className="text-danger">*</span>
        </label>
        <input
          id="category-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={isSubmitting}
          className="w-full px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed"
          placeholder="例: 外食費、ジム代"
          maxLength={CATEGORY_NAME_MAX_LENGTH}
        />
      </div>

      <div className="flex gap-2 pt-2">
        <button
          type="button"
          onClick={onCancel}
          disabled={isSubmitting}
          className="px-4 py-2 text-sm font-medium rounded-lg bg-secondary hover:bg-secondary-hover text-secondary-foreground disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          キャンセル
        </button>
        <button
          type="submit"
          disabled={isSubmitting}
          className="px-4 py-2 text-sm font-medium rounded-lg bg-primary hover:bg-primary-hover text-primary-foreground disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {isSubmitting ? "保存中..." : "保存"}
        </button>
      </div>
    </form>
  );
}
