import { Dashboard as DashboardType } from "@/lib/types/dashboard";

type DashboardProps = {
  dashboard: DashboardType;
};

function getRemainingColor(remaining: number, variableBudget: number): string {
  // 残額がマイナスの場合は常に赤
  if (remaining < 0) return "red";
  
  // 変動費がゼロ以下の場合は赤（使える余裕がない状態）
  if (variableBudget <= 0) {
    return "red";
  }
  
  // 通常の割合ベースの色分け
  const percentage = (remaining / variableBudget) * 100;
  
  if (percentage >= 70) return "green";
  if (percentage >= 30) return "yellow";
  return "red";
}

function formatAmount(amount: number): string {
  return amount.toLocaleString("ja-JP");
}

export function Dashboard({ dashboard }: DashboardProps) {
  const colorClass = getRemainingColor(dashboard.remaining, dashboard.variable_budget);
  
  const remainingColorClasses = {
    green: "text-success",
    yellow: "text-warning", 
    red: "text-danger"
  }

  return (
    <div className="w-full max-w-5xl mx-auto space-y-6">
      {/* Remaining Section */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-6 text-center">
        <h2 className="text-lg font-semibold text-foreground mb-4">今月あといくら使える？</h2>
        <div className={`text-5xl font-bold my-4 ${remainingColorClasses[colorClass as keyof typeof remainingColorClasses]}`}>
          ¥{formatAmount(dashboard.remaining)}
        </div>
        <p className="text-sm text-muted-foreground">残り使える金額</p>
      </section>

      {/* Summary Grid */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4">月次サマリー</h2>
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="space-y-1">
            <div className="text-sm text-muted-foreground">収入</div>
            <div className="text-xl font-semibold text-foreground">¥{formatAmount(dashboard.income)}</div>
          </div>
          <div className="space-y-1">
            <div className="text-sm text-muted-foreground">貯金目標</div>
            <div className="text-xl font-semibold text-foreground">¥{formatAmount(dashboard.saving_goal)}</div>
          </div>
          <div className="space-y-1">
            <div className="text-sm text-muted-foreground">固定費</div>
            <div className="text-xl font-semibold text-foreground">¥{formatAmount(dashboard.fixed_costs)}</div>
          </div>
          <div className="space-y-1">
            <div className="text-sm text-muted-foreground">変動費</div>
            <div className="text-xl font-semibold text-foreground">¥{formatAmount(dashboard.variable_budget)}</div>
          </div>
        </div>
      </section>

      {/* Expenses Grid */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4">今月の支出</h2>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1">
            <div className="text-sm text-muted-foreground">確定支出</div>
            <div className="text-xl font-semibold text-foreground">¥{formatAmount(dashboard.confirmed_expenses)}</div>
          </div>
          <div className="space-y-1">
            <div className="text-sm text-muted-foreground">予定支出</div>
            <div className="text-xl font-semibold text-foreground">¥{formatAmount(dashboard.planned_expenses)}</div>
          </div>
        </div>
      </section>
    </div>
  );
}
