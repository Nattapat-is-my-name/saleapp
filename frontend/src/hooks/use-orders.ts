"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { Order, CreateOrderInput, UpdateOrderStatusInput } from "@/types/order";
import type { PaginatedResponse } from "@/types/api";

interface OrdersFilter {
  page?: number;
  limit?: number;
  status?: string;
  customerId?: string;
  startDate?: string;
  endDate?: string;
}

export function useOrders(filters: OrdersFilter = {}) {
  const { page = 1, limit = 20, status, customerId, startDate, endDate } = filters;

  return useQuery({
    queryKey: ["orders", { page, limit, status, customerId, startDate, endDate }],
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<Order>>("/orders", {
        page,
        limit,
        status,
        customerId,
        startDate,
        endDate,
      });
      return response;
    },
  });
}

export function useOrder(id: string) {
  return useQuery({
    queryKey: ["orders", id],
    queryFn: async () => {
      const response = await api.get<{ data: Order }>(`/orders/${id}`);
      return response.data;
    },
    enabled: !!id,
  });
}

export function useCreateOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateOrderInput) => {
      const response = await api.post<{ data: Order }>("/orders", input);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      queryClient.invalidateQueries({ queryKey: ["products"] }); // Update stock
    },
  });
}

export function useUpdateOrderStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...input }: UpdateOrderStatusInput & { id: string }) => {
      const response = await api.put<{ data: Order }>(`/orders/${id}/status`, input);
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      queryClient.setQueryData(["orders", data.id], data);
    },
  });
}

export function useCancelOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/orders/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    },
  });
}
