export interface Order {
  id: string;
  orderNumber: string;
  customerId: string | null;
  customer?: Customer | null;
  userId: string;
  user?: User;
  status: OrderStatus;
  subtotal: number;
  tax: number;
  discount: number;
  total: number;
  paymentMethod: string;
  notes: string;
  items: OrderItem[];
  createdAt: string;
  updatedAt: string;
}

export interface OrderItem {
  id: string;
  orderId: string;
  productId: string;
  product?: Product;
  quantity: number;
  unitPrice: number;
  discount: number;
  total: number;
}

export interface Product {
  id: string;
  sku: string;
  name: string;
  price: number;
}

export interface Customer {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
}

export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
}

export type OrderStatus = "pending" | "completed" | "cancelled" | "refunded";

export interface CreateOrderInput {
  customerId?: string;
  paymentMethod: string;
  notes?: string;
  items: Array<{
    productId: string;
    quantity: number;
    discount?: number;
  }>;
}

export interface UpdateOrderStatusInput {
  status: OrderStatus;
}
