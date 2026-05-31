import {
  getCategories,
  createCategory,
  updateCategory,
  deleteCategory,
  reorderCategories,
} from "./categories";
import { Category, CreateCategoryInput, UpdateCategoryInput, ReorderCategoryItem } from "@pace/core/types/category";

type GlobalWithFetch = typeof globalThis & { fetch: jest.Mock };

beforeEach(() => {
  jest.resetAllMocks();
  (global as GlobalWithFetch).fetch = jest.fn();
  jest.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  jest.restoreAllMocks();
});

const API_BASE = "http://localhost:8080";

const sampleCategory: Category = {
  id: 1,
  name: "食費",
  category_type: "default",
  sort_order: 1,
};

// ---------------------------------------------------------------------------
// getCategories
// ---------------------------------------------------------------------------
describe("getCategories", () => {
  it("calls fetch with correct URL and method and returns categories on success", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: true,
      json: async () => ({ categories: [sampleCategory] }),
    });

    const result = await getCategories();

    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledTimes(1);
    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledWith(
      `${API_BASE}/categories`,
      expect.objectContaining({ method: "GET" })
    );
    expect(result).toEqual([sampleCategory]);
  });

  it("throws Error when response.ok is false", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => { throw new Error("not json"); },
    });

    await expect(getCategories()).rejects.toThrow("カテゴリの取得に失敗しました");
  });

  it("throws Error when response body is not a valid categories array", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: true,
      json: async () => ({ categories: null }),
    });

    await expect(getCategories()).rejects.toThrow("カテゴリのレスポンスが正しくありません");
  });
});

// ---------------------------------------------------------------------------
// createCategory
// ---------------------------------------------------------------------------
describe("createCategory", () => {
  const input: CreateCategoryInput = { name: "外食" };

  it("calls fetch with correct URL, method, headers, and body and returns category on success", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: true,
      json: async () => ({ category: { ...sampleCategory, name: "外食" } }),
    });

    const result = await createCategory(input);

    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledTimes(1);
    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledWith(
      `${API_BASE}/categories`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(input),
      }
    );
    expect(result.name).toBe("外食");
  });

  it("throws Error when response.ok is false", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => { throw new Error("not json"); },
    });

    await expect(createCategory(input)).rejects.toThrow("カテゴリの作成に失敗しました");
  });

  it("surfaces server error message from JSON body", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({ error: "同じ名前のカテゴリがすでに存在します" }),
    });

    await expect(createCategory(input)).rejects.toThrow("同じ名前のカテゴリがすでに存在します");
  });
});

// ---------------------------------------------------------------------------
// updateCategory
// ---------------------------------------------------------------------------
describe("updateCategory", () => {
  const input: UpdateCategoryInput = { name: "外食費" };

  it("calls fetch with correct URL, method, headers, and body and returns category on success", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: true,
      json: async () => ({ category: { ...sampleCategory, name: "外食費" } }),
    });

    const result = await updateCategory(1, input);

    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledTimes(1);
    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledWith(
      `${API_BASE}/categories/1`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(input),
      }
    );
    expect(result.name).toBe("外食費");
  });

  it("throws Error when response.ok is false", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 404,
      json: async () => { throw new Error("not json"); },
    });

    await expect(updateCategory(1, input)).rejects.toThrow("カテゴリの更新に失敗しました");
  });

  it("surfaces server error message from JSON body", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({ error: "同じ名前のカテゴリがすでに存在します" }),
    });

    await expect(updateCategory(1, input)).rejects.toThrow("同じ名前のカテゴリがすでに存在します");
  });
});

// ---------------------------------------------------------------------------
// deleteCategory
// ---------------------------------------------------------------------------
describe("deleteCategory", () => {
  it("calls fetch with correct URL and method and resolves on success", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({ ok: true });

    await expect(deleteCategory(1)).resolves.toBeUndefined();

    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledTimes(1);
    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledWith(
      `${API_BASE}/categories/1`,
      expect.objectContaining({ method: "DELETE" })
    );
  });

  it("throws Error when response.ok is false", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 404,
      json: async () => { throw new Error("not json"); },
    });

    await expect(deleteCategory(1)).rejects.toThrow("カテゴリの削除に失敗しました");
  });

  it("surfaces server error message from JSON body", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({ error: "このカテゴリは支出で使用されているため削除できません" }),
    });

    await expect(deleteCategory(1)).rejects.toThrow("このカテゴリは支出で使用されているため削除できません");
  });
});

// ---------------------------------------------------------------------------
// reorderCategories
// ---------------------------------------------------------------------------
describe("reorderCategories", () => {
  const items: ReorderCategoryItem[] = [
    { id: 1, sort_order: 2 },
    { id: 2, sort_order: 1 },
  ];

  it("calls fetch with correct URL, method, headers, and body and resolves on success", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({ ok: true });

    await expect(reorderCategories(items)).resolves.toBeUndefined();

    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledTimes(1);
    expect((global as GlobalWithFetch).fetch).toHaveBeenCalledWith(
      `${API_BASE}/categories/reorder`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ items }),
      }
    );
  });

  it("throws Error when response.ok is false", async () => {
    (global as GlobalWithFetch).fetch.mockResolvedValue({
      ok: false,
      status: 404,
      json: async () => { throw new Error("not json"); },
    });

    await expect(reorderCategories(items)).rejects.toThrow("カテゴリの並び替えに失敗しました");
  });
});
