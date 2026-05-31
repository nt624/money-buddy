import { Dashboard as DashboardType } from "@pace/core/types/dashboard";

type DashboardProps = {
  dashboard: DashboardType;
};

type RemainingColor = "green" | "yellow" | "red";

function getRemainingColor(remaining: number, variableBudget: number): RemainingColor {
  if (remaining < 0) return "red";
  if (variableBudget <= 0) return "red";
  const percentage = (remaining / variableBudget) * 100;
  if (percentage >= 70) return "green";
  if (percentage >= 30) return "yellow";
  return "red";
}

function formatAmount(amount: number): string {
  return new Intl.NumberFormat("ja-JP", { style: "currency", currency: "JPY" }).format(amount);
}

const remainingColorClasses: Record<RemainingColor, string> = {
  green: "text-success",
  yellow: "text-warning",
  red: "text-danger",
};

type SummaryRowProps = {
  operator?: string;
  label: string;
  amount: number;
  isResult?: boolean;
};

function SummaryRow({ operator, label, amount, isResult = false }: SummaryRowProps) {
  return (
    <div className={`flex items-center justify-between gap-4 py-2 ${isResult ? "pt-3" : ""}`}>
      <div className="flex items-center gap-2 min-w-0">
        <span
          className="w-4 shrink-0 text-center text-sm text-muted-foreground select-none font-medium"
          aria-hidden="true"
        >
          {operator ?? ""}
        </span>
        <span className={`text-sm ${isResult ? "font-semibold text-foreground" : "text-muted-foreground"}`}>
          {label}
        </span>
      </div>
      <span className={`tabular-nums shrink-0 ${isResult ? "text-base font-bold text-foreground" : "text-sm font-medium text-foreground"}`}>
        {formatAmount(amount)}
      </span>
    </div>
  );
}

export function Dashboard({ dashboard }: DashboardProps) {
  const colorClass = getRemainingColor(dashboard.remaining, dashboard.variable_budget);
  const totalExpenses = dashboard.confirmed_expenses + dashboard.planned_expenses;

  return (
    <div className="w-full max-w-2xl mx-auto space-y-4">

      {/* ── 残額カード ── */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-5 sm:p-6">
        <h2 className="text-sm font-medium text-muted-foreground mb-1">今月あといくら使える？</h2>
        <div
          className={`text-4xl sm:text-5xl font-bold tracking-tight mb-1 ${remainingColorClasses[colorClass]}`}
        >
          {formatAmount(dashboard.remaining)}
        </div>
        <p className="text-xs text-muted-foreground mb-4">残り使える金額</p>

        {/* Surface: 支出内訳 */}
        <div className="bg-muted rounded-lg px-4 py-3 space-y-1">
          <p className="text-xs text-muted-foreground mb-2">
            使えるお金 − 支出合計 ＝ 残り使える金額
          </p>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">確定支出</span>
            <span className="text-sm font-medium tabular-nums text-foreground">
              {formatAmount(dashboard.confirmed_expenses)}
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">予定支出</span>
            <span className="text-sm font-medium tabular-nums text-foreground">
              {formatAmount(dashboard.planned_expenses)}
            </span>
          </div>
          <div className="border-t border-border/60 mt-2 pt-2 flex justify-between items-center">
            <span className="text-sm text-muted-foreground">支出合計</span>
            <span className="text-sm font-semibold tabular-nums text-foreground">
              {formatAmount(totalExpenses)}
            </span>
          </div>
        </div>
      </section>

      {/* ── 月次サマリー ── */}
      <section className="bg-card border border-border rounded-lg shadow-sm p-5 sm:p-6">
        <h2 className="text-sm font-medium text-muted-foreground mb-1">月次サマリー</h2>
        <div
          role="math"
          aria-label="使えるお金の計算式"
          className="divide-y divide-border"
        >
          <SummaryRow label="収入" amount={dashboard.income} />
          <SummaryRow operator="−" label="貯金目標" amount={dashboard.saving_goal} />
          <SummaryRow operator="−" label="固定費" amount={dashboard.fixed_costs} />
        </div>

        {/* 計算結果の区切り線 */}
        <div className="border-t-2 border-foreground/20 mt-1" />

        <div className="flex items-center justify-between gap-4 pt-3">
          <div className="flex items-center gap-2">
            <span className="w-4 shrink-0 text-center text-sm text-muted-foreground select-none font-medium" aria-hidden="true">
              =
            </span>
            <span className="text-sm font-semibold text-foreground">使えるお金</span>
          </div>
          <span className="text-base font-bold tabular-nums text-foreground">
            {formatAmount(dashboard.variable_budget)}
          </span>
        </div>
      </section>

    </div>
  );
}
