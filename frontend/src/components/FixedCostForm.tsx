import { useState, useEffect } from "react";
import { FixedCost, UpdateFixedCostInput } from "@/lib/types/fixed-cost";
import styles from "./FixedCostForm.module.css";

type Props = {
  fixedCost: FixedCost | null;
  onSubmit: (input: UpdateFixedCostInput) => void;
  onCancel: () => void;
  isSubmitting: boolean;
};

export function FixedCostForm({
  fixedCost,
  onSubmit,
  onCancel,
  isSubmitting,
}: Props) {
  const [name, setName] = useState("");
  const [amount, setAmount] = useState("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (fixedCost) {
      setName(fixedCost.name);
      setAmount(fixedCost.amount.toString());
    } else {
      setName("");
      setAmount("");
    }
    setError(null);
  }, [fixedCost]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const trimmedName = name.trim();
    const amountNum = parseInt(amount, 10);

    // Validation
    if (!trimmedName) {
      setError("固定費名を入力してください");
      return;
    }

    if (trimmedName.length > 100) {
      setError("固定費名は100文字以内で入力してください");
      return;
    }

    if (isNaN(amountNum) || amountNum <= 0) {
      setError("金額は1円以上で入力してください");
      return;
    }

    if (amountNum > 1000000000) {
      setError("金額は10億円以下で入力してください");
      return;
    }

    onSubmit({
      name: trimmedName,
      amount: amountNum,
    });
  };

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <h3 className={styles.title}>
        {fixedCost ? "固定費を編集" : "固定費を追加"}
      </h3>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.field}>
        <label htmlFor="name" className={styles.label}>
          固定費名 <span className={styles.required}>*</span>
        </label>
        <input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={isSubmitting}
          className={styles.input}
          placeholder="例: 家賃、光熱費"
          maxLength={100}
        />
      </div>

      <div className={styles.field}>
        <label htmlFor="amount" className={styles.label}>
          金額（円） <span className={styles.required}>*</span>
        </label>
        <input
          id="amount"
          type="number"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          disabled={isSubmitting}
          className={styles.input}
          placeholder="10000"
          min="1"
          max="1000000000"
        />
      </div>

      <div className={styles.actions}>
        <button
          type="button"
          onClick={onCancel}
          disabled={isSubmitting}
          className={`${styles.button} ${styles.cancelButton}`}
        >
          キャンセル
        </button>
        <button
          type="submit"
          disabled={isSubmitting}
          className={`${styles.button} ${styles.submitButton}`}
        >
          {isSubmitting ? "保存中..." : "保存"}
        </button>
      </div>
    </form>
  );
}
