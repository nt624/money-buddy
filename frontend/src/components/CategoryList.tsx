'use client'

import { Category, ReorderCategoryItem } from "@/lib/types/category";

type Props = {
  categories: Category[];
  onEdit: (category: Category) => void;
  onDelete: (id: number) => void;
  onReorder: (items: ReorderCategoryItem[]) => void;
  isSubmitting: boolean;
};

const GROUP_LABELS: Record<string, string> = {
  default: "デフォルト",
  user: "ユーザー定義",
  other: "その他",
};

export function CategoryList({ categories, onEdit, onDelete, onReorder, isSubmitting }: Props) {
  if (categories.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        カテゴリが登録されていません
      </div>
    );
  }

  const groups: Array<{ type: string; items: Category[] }> = [
    { type: "default", items: categories.filter((c) => c.category_type === "default") },
    { type: "user", items: categories.filter((c) => c.category_type === "user") },
    { type: "other", items: categories.filter((c) => c.category_type === "other") },
  ].filter((g) => g.items.length > 0);

  const handleMove = (group: Category[], current: Category, direction: "up" | "down") => {
    const idx = group.indexOf(current);
    const swapIdx = direction === "up" ? idx - 1 : idx + 1;
    if (swapIdx < 0 || swapIdx >= group.length) return;

    const target = group[swapIdx];
    onReorder([
      { id: current.id, sort_order: target.sort_order },
      { id: target.id, sort_order: current.sort_order },
    ]);
  };

  return (
    <div className="space-y-6">
      {groups.map(({ type, items }) => (
        <div key={type}>
          <div className="text-xs font-semibold text-muted-foreground uppercase tracking-wide mb-2 px-1">
            {GROUP_LABELS[type]}
          </div>
          <ul className="list-none p-0 space-y-2">
            {items.map((cat, idx) => (
              <li
                key={cat.id}
                className="p-3 border border-border bg-card rounded-lg flex items-center gap-2 shadow-sm"
              >
                {/* 並び替えボタン (その他グループは非表示) */}
                {type !== "other" && (
                  <div className="flex flex-col gap-0.5 flex-shrink-0">
                    <button
                      onClick={() => handleMove(items, cat, "up")}
                      disabled={isSubmitting || idx === 0}
                      aria-label="上に移動"
                      className="p-0.5 text-muted-foreground hover:text-foreground disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                    >
                      ▲
                    </button>
                    <button
                      onClick={() => handleMove(items, cat, "down")}
                      disabled={isSubmitting || idx === items.length - 1}
                      aria-label="下に移動"
                      className="p-0.5 text-muted-foreground hover:text-foreground disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                    >
                      ▼
                    </button>
                  </div>
                )}

                <span className="flex-1 text-sm font-medium text-foreground">{cat.name}</span>

                <div className="flex gap-2 flex-shrink-0">
                  <button
                    onClick={() => onEdit(cat)}
                    disabled={isSubmitting}
                    className="px-3 py-1.5 text-xs font-medium rounded bg-primary hover:bg-primary-hover text-primary-foreground disabled:opacity-60 disabled:cursor-not-allowed transition-colors"
                  >
                    編集
                  </button>
                  <button
                    onClick={() => onDelete(cat.id)}
                    disabled={isSubmitting}
                    className="px-3 py-1.5 text-xs font-medium rounded bg-danger hover:bg-danger-hover text-danger-foreground disabled:opacity-60 disabled:cursor-not-allowed transition-colors"
                  >
                    削除
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
}
