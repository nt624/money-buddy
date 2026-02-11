import { useState, useEffect } from "react";
import { FixedCost, FixedCostInput } from "@/lib/types/fixed-cost";
import { BUSINESS_MAX_AMOUNT, FIXED_COST_NAME_MAX_LENGTH } from "@/lib/constants";

type Props = {
  fixedCost: FixedCost | null;
  onSubmit: (input: FixedCostInput) => void;
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

    if (trimmedName.length > FIXED_COST_NAME_MAX_LENGTH) {
      setError("固定費名は100文字以内で入力してください");
      return;
    }

    if (isNaN(amountNum) || amountNum <= 0) {
      setError("金額は1円以上で入力してください");
      return;
    }

    if (amountNum > BUSINESS_MAX_AMOUNT) {
      setError("金額は10億円以下で入力してください");
      return;
    }

    onSubmit({
      name: trimmedName,
      amount: amountNum,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h3 className="text-lg font-semibold text-foreground">
        {fixedCost ? "固定費を編集" : "固定費を追加"}
      </h3>

      {error && (
        <div className="p-3 bg-danger/10 border border-danger rounded-lg text-danger text-sm">
          {error}
        </div>
      )}

      <div className="space-y-2">
        <label htmlFor="name" className="block text-sm font-medium text-foreground">
          固定費名 <span className="text-danger">*</span>
        </label>
        <input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={isSubmitting}
          className="w-full px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed"
          placeholder="例: 家賃、光熱費"
          maxLength={FIXED_COST_NAME_MAX_LENGTH}
        />
      </div>

      <div className="space-y-2">
        <label htmlFor="amount" className="block text-sm font-medium text-foreground">
          金額（円） <span className="text-danger">*</span>
        </label>
        <input
          id="amount"
          type="number"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          disabled={isSubmitting}
          className="w-full px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed"
          placeholder="10000"
          min="1"
          max={BUSINESS_MAX_AMOUNT}
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
