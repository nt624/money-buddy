import { FixedCost } from "@/lib/types/fixed-cost";
import styles from "./FixedCostList.module.css";

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
      <div className={styles.emptyState}>
        固定費が登録されていません
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.totalRow}>
        <span className={styles.totalLabel}>固定費合計:</span>
        <span className={styles.totalAmount}>¥{totalAmount.toLocaleString()}</span>
      </div>
      <ul className={styles.list}>
        {fixedCosts.map((fc) => (
          <li key={fc.id} className={styles.listItem}>
            <div className={styles.itemContent}>
              <div className={styles.itemName}>{fc.name}</div>
              <div className={styles.itemAmount}>
                ¥{fc.amount.toLocaleString()}
              </div>
            </div>
            <div className={styles.itemActions}>
              <button
                onClick={() => onEdit(fc)}
                disabled={isSubmitting}
                className={`${styles.button} ${styles.editButton}`}
              >
                編集
              </button>
              <button
                onClick={() => onDelete(fc.id)}
                disabled={isSubmitting}
                className={`${styles.button} ${styles.deleteButton}`}
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
