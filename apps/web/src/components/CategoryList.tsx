'use client'

import {
  DndContext,
  closestCenter,
  MouseSensor,
  TouchSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core'
import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
  arrayMove,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

import { Category, ReorderCategoryItem } from "@pace/core/types/category"

type Props = {
  categories: Category[]
  onEdit: (category: Category) => void
  onDelete: (id: number) => void
  onReorder: (items: ReorderCategoryItem[]) => void
  isSubmitting: boolean
}

type ItemProps = {
  cat: Category
  isSubmitting: boolean
  onEdit: (category: Category) => void
  onDelete: (id: number) => void
}

function SortableCategoryItem({ cat, isSubmitting, onEdit, onDelete }: ItemProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: cat.id })

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    zIndex: isDragging ? 10 : undefined,
  }

  return (
    <li
      ref={setNodeRef}
      style={style}
      className="p-3 border border-border bg-card rounded-lg flex items-center gap-2 shadow-sm"
    >
      {/* ドラッグハンドル */}
      <button
        {...attributes}
        {...listeners}
        disabled={isSubmitting}
        aria-label="ドラッグして並び替え"
        className="p-2 -ml-1 flex-shrink-0 text-muted-foreground hover:text-foreground disabled:opacity-30 disabled:cursor-not-allowed touch-manipulation cursor-grab active:cursor-grabbing rounded"
      >
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
          <circle cx="5.5" cy="3.5" r="1.25"/>
          <circle cx="10.5" cy="3.5" r="1.25"/>
          <circle cx="5.5" cy="8" r="1.25"/>
          <circle cx="10.5" cy="8" r="1.25"/>
          <circle cx="5.5" cy="12.5" r="1.25"/>
          <circle cx="10.5" cy="12.5" r="1.25"/>
        </svg>
      </button>

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
  )
}

function StaticCategoryItem({ cat, isSubmitting, onEdit, onDelete }: ItemProps) {
  return (
    <li className="p-3 border border-border bg-card rounded-lg flex items-center gap-2 shadow-sm opacity-70">
      {/* 並び替え不可のプレースホルダー（ハンドルなし） */}
      <div className="w-8 flex-shrink-0" />

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
  )
}

export function CategoryList({ categories, onEdit, onDelete, onReorder, isSubmitting }: Props) {
  // スクロールと誤操作を区別するための activationConstraint
  const sensors = useSensors(
    useSensor(MouseSensor, {
      activationConstraint: { distance: 5 },
    }),
    useSensor(TouchSensor, {
      activationConstraint: { delay: 250, tolerance: 5 },
    })
  )

  if (categories.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        カテゴリが登録されていません
      </div>
    )
  }

  const mainCategories = categories.filter((c) => c.category_type !== "other")
  const otherCategories = categories.filter((c) => c.category_type === "other")

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over || active.id === over.id) return

    const oldIndex = mainCategories.findIndex((c) => c.id === active.id)
    const newIndex = mainCategories.findIndex((c) => c.id === over.id)
    if (oldIndex === -1 || newIndex === -1) return

    const reordered = arrayMove(mainCategories, oldIndex, newIndex)
    onReorder(reordered.map((c, i) => ({ id: c.id, sort_order: i + 1 })))
  }

  return (
    <ul className="list-none p-0 space-y-2">
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={handleDragEnd}
      >
        <SortableContext
          items={mainCategories.map((c) => c.id)}
          strategy={verticalListSortingStrategy}
        >
          {mainCategories.map((cat) => (
            <SortableCategoryItem
              key={cat.id}
              cat={cat}
              isSubmitting={isSubmitting}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </SortableContext>
      </DndContext>

      {otherCategories.map((cat) => (
        <StaticCategoryItem
          key={cat.id}
          cat={cat}
          isSubmitting={isSubmitting}
          onEdit={onEdit}
          onDelete={onDelete}
        />
      ))}
    </ul>
  )
}
