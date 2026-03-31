"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { Product, CreateProductInput, UpdateProductInput } from "@/types/product";
import type { PaginatedResponse } from "@/types/api";

interface ProductsFilter {
  page?: number;
  limit?: number;
  search?: string;
  categoryId?: string;
  isActive?: boolean;
}

export function useProducts(filters: ProductsFilter = {}) {
  const { page = 1, limit = 20, search, categoryId, isActive } = filters;

  return useQuery({
    queryKey: ["products", { page, limit, search, categoryId, isActive }],
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<Product>>("/products", {
        page,
        limit,
        search,
        categoryId,
        isActive,
      });
      return response;
    },
  });
}

export function useProduct(id: string) {
  return useQuery({
    queryKey: ["products", id],
    queryFn: async () => {
      const response = await api.get<{ data: Product }>(`/products/${id}`);
      return response.data;
    },
    enabled: !!id,
  });
}

export function useCreateProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateProductInput) => {
      const response = await api.post<{ data: Product }>("/products", input);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });
}

export function useUpdateProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...input }: UpdateProductInput & { id: string }) => {
      const response = await api.put<{ data: Product }>(`/products/${id}`, input);
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
      queryClient.setQueryData(["products", data.id], data);
    },
  });
}

export function useDeleteProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/products/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });
}
