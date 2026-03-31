"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";

export interface DashboardProduct {
  id: string;
  sku: string;
  name: string;
  description: string;
  price: number;
  cost: number;
  stock: number;
  categoryId: string | null;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface TopProduct {
  product: DashboardProduct;
  quantity_sold: number;
}

export interface DashboardResponse {
  today_sales: number;
  today_orders: number;
  today_revenue: number;
  top_products: TopProduct[];
  low_stock: DashboardProduct[];
}

export function useDashboard() {
  return useQuery({
    queryKey: ["dashboard"],
    queryFn: async () => {
      const response = await api.get<{ data: DashboardResponse }>("/reports/dashboard");
      return response.data;
    },
    retry: 1,
    staleTime: 30000,
  });
}
