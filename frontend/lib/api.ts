import {
  CreateOrderPayload,
  CreateProductPayload,
  Order,
  Product,
  UpdateProductPayload,
  UpdateOrderPayload,
} from "@/types";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body?.error ?? `Request failed: ${res.status}`);
  }

  // 204 No Content
  if (res.status === 204) return undefined as T;

  return res.json();
}

// ─── Products ────────────────────────────────────────────────────────────────

export const productApi = {
  getAll: (): Promise<Product[]> => request<Product[]>("/products"),

  create: (data: CreateProductPayload): Promise<Product> =>
    request<Product>("/products", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  update: (id: string, data: UpdateProductPayload): Promise<void> =>
    request<void>(`/products/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  delete: (id: string): Promise<void> =>
    request<void>(`/products/${id}`, { method: "DELETE" }),
};

// ─── Orders ──────────────────────────────────────────────────────────────────

export const orderApi = {
  getAll: (): Promise<Order[]> => request<Order[]>("/orders"),

  getById: (id: string): Promise<Order> => request<Order>(`/orders/${id}`),

  create: (data: CreateOrderPayload): Promise<Order> =>
    request<Order>("/orders", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  update: (id: string, data: UpdateOrderPayload): Promise<Order> =>
    request<Order>(`/orders/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  delete: (id: string): Promise<void> =>
    request<void>(`/orders/${id}`, { method: "DELETE" }),

  // Force delete untuk testing/cleanup (biasanya butuh admin key)
  forceDelete: (id: string, adminKey?: string): Promise<void> =>
    request<void>(`/orders/force/${id}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        ...(adminKey && { "X-Admin-Key": adminKey }),
      },
    }),

  deleteAll: (
    adminKey?: string,
  ): Promise<{ message: string; deleted: number }> =>
    request<{ message: string; deleted: number }>(`/orders/admin/delete-all`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        ...(adminKey && { "X-Admin-Key": adminKey }),
      },
    }),
};
