import { createExpense } from "./expenses";
import { CreateExpenseInput, Expense } from "../types/expense";

// Mock global.fetch for all tests
beforeEach(() => {
  jest.resetAllMocks();
  (global as any).fetch = jest.fn();
  // Mock console.error to avoid cluttering test output
  jest.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  jest.restoreAllMocks();
});

describe("createExpense", () => {
  const API_URL = "http://localhost:8080/expenses";

  it("calls fetch with correct URL, method, headers, and body and returns expense on success", async () => {
    const input: CreateExpenseInput = {
      amount: 1200,
      category_id: 3,
      memo: "通勤",
      spent_at: "2025-12-31",
    };

    const expectedExpense: Expense = {
      id: 1,
      amount: 1200,
      category: { id: 3, name: "交通費" },
      memo: "通勤",
      spent_at: "2025-12-31",
    };

    (global as any).fetch.mockResolvedValue({
      ok: true,
      json: async () => ({ expense: expectedExpense }),
    });

    const result = await createExpense(input);

    expect((global as any).fetch).toHaveBeenCalledTimes(1);
    expect((global as any).fetch).toHaveBeenCalledWith(API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(input),
    });

    expect(result).toEqual(expectedExpense);
  });

  it("throws Error when response.ok is false", async () => {
    const input: CreateExpenseInput = {
      amount: 500,
      category_id: 2,
      memo: "snack",
      spent_at: "2025-12-30",
    };

    (global as any).fetch.mockResolvedValue({
      ok: false,
      status: 500,
      text: async () => "internal error",
    });

    await expect(createExpense(input)).rejects.toThrow('支出の作成に失敗しました');
  });
});