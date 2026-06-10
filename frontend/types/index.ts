// types/index.ts

export interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
  created_at: string;
  updated_at: string;
}

export interface Order {
  id: string;
  product_id: string;
  product_name: string; // Denormalized field
  product_price: number; // Denormalized field
  quantity: number;
  created_at: string;
}

// Payload types
export type CreateProductPayload = Omit<
  Product,
  "id" | "created_at" | "updated_at"
>;
export type UpdateProductPayload = Partial<
  Omit<Product, "id" | "created_at" | "updated_at">
>;

export interface CreateOrderPayload {
  productId: string;
  quantity: number;
}

export interface UpdateOrderPayload {
  quantity: number;
}

// Untuk response dari force delete/delete all
export interface DeleteAllResponse {
  message: string;
  deleted: number;
}

// Extended type untuk display di UI
export interface OrderWithTotal extends Order {
  total_price: number;
}

// Helper function untuk menghitung total
export function calculateOrderTotal(order: Order): number {
  return order.product_price * order.quantity;
}
