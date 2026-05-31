import { FixedCost } from "@/lib/types/fixed-cost";

type Props = {
  fixedCosts: FixedCost[];
  onEdit: (fixedCost: FixedCost) => void;
  onDelete: (id: number) => void;
  isSubmitting: boolean;
};

export function FixedCostList({
  fixedCosts,
  onEdit,
  onDelete,
  isSubmitting,
}: Props) {
  const totalAmount = fixedCosts.reduce((sum, fc) => sum + fc.amount, 0);

  if (fixedCosts.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        固定費が登録されていません
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* 固定費合計 */}
      <div className="flex justify-between items-center p-4 bg-secondary/20 border border-border rounded-lg">
        <span className="font-semibold text-foreground">固定費合計</span>
        <span className="text-xl font-bold text-foreground">¥{totalAmount.toLocaleString()}</span>
      </div>
      
      {/* 固定費一覧 */}
      <ul className="list-none p-0 space-y-4">
        {fixedCosts.map((fc) => (
          <li 
            key={fc.id} 
            className="p-4 border border-border bg-card rounded-lg flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 shadow-sm hover:shadow-md transition-shadow"
          >
            <div className="flex-1 space-y-2">
              <div className="text-lg font-semibold text-foreground">{fc.name}</div>
              <div className="text-xl font-bold text-primary">¥{fc.amount.toLocaleString()}</div>
            </div>
            <div className="flex gap-2 sm:flex-shrink-0">
              <button
                onClick={() => onEdit(fc)}
                disabled={isSubmitting}
                className="
                  px-4 py-2 text-sm font-medium rounded
                  bg-primary hover:bg-primary-hover text-primary-foreground
                  disabled:opacity-60 disabled:cursor-not-allowed
                  transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
                "
              >
                編集
              </button>
              <button
                onClick={() => onDelete(fc.id)}
                disabled={isSubmitting}
                className="
                  px-4 py-2 text-sm font-medium rounded
                  bg-danger hover:bg-danger-hover text-danger-foreground
                  disabled:opacity-60 disabled:cursor-not-allowed
                  transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
                "
              >
                削除
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
