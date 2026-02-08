import { Dashboard as DashboardType } from "@/lib/types/dashboard";
import styles from "./Dashboard.module.css";

type DashboardProps = {
  dashboard: DashboardType;
};

function getRemainingColor(remaining: number, variableBudget: number): string {
  if (variableBudget === 0) return "green";
  
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

  return (
    <div className={styles.container}>
      {/* Remaining Section */}
      <section className={`${styles.section} ${styles.remainingSection}`}>
        <h2 className={styles.sectionTitle}>今月あといくら使える？</h2>
        <div className={`${styles.remainingAmount} ${styles[colorClass]}`}>
          ¥{formatAmount(dashboard.remaining)}
        </div>
        <p className={styles.remainingLabel}>残り使える金額</p>
      </section>

      {/* Summary Grid */}
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>月次サマリー</h2>
        <div className={styles.summaryGrid}>
          <div className={styles.summaryItem}>
            <div className={styles.summaryLabel}>収入</div>
            <div className={styles.summaryValue}>¥{formatAmount(dashboard.income)}</div>
          </div>
          <div className={styles.summaryItem}>
            <div className={styles.summaryLabel}>貯金目標</div>
            <div className={styles.summaryValue}>¥{formatAmount(dashboard.saving_goal)}</div>
          </div>
          <div className={styles.summaryItem}>
            <div className={styles.summaryLabel}>固定費</div>
            <div className={styles.summaryValue}>¥{formatAmount(dashboard.fixed_costs)}</div>
          </div>
          <div className={styles.summaryItem}>
            <div className={styles.summaryLabel}>変動費</div>
            <div className={styles.summaryValue}>¥{formatAmount(dashboard.variable_budget)}</div>
          </div>
        </div>
      </section>

      {/* Expenses Grid */}
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>今月の支出</h2>
        <div className={styles.expensesGrid}>
          <div className={styles.expenseItem}>
            <div className={styles.expenseLabel}>確定支出</div>
            <div className={styles.expenseValue}>¥{formatAmount(dashboard.confirmed_expenses)}</div>
          </div>
          <div className={styles.expenseItem}>
            <div className={styles.expenseLabel}>予定支出</div>
            <div className={styles.expenseValue}>¥{formatAmount(dashboard.planned_expenses)}</div>
          </div>
        </div>
      </section>
    </div>
  );
}
