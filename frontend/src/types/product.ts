export interface Product {
  id: string;
  sku: string;
  name: string;
  description: string;
  price: number;
  cost: number;
  stock: number;
  categoryId: string | null;
  category?: Category | null;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Category {
  id: string;
  name: string;
  description: string;
  createdAt: string;
}

export interface CreateProductInput {
  sku: string;
  name: string;
  description?: string;
  price: number;
  cost?: number;
  stock?: number;
  categoryId?: string;
  isActive?: boolean;
}

export interface UpdateProductInput extends Partial<CreateProductInput> {}
