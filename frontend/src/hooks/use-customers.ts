"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { Customer, CreateCustomerInput, UpdateCustomerInput } from "@/types/customer";
import type { PaginatedResponse } from "@/types/api";

interface CustomersFilter {
  page?: number;
  limit?: number;
  search?: string;
}

export function useCustomers(filters: CustomersFilter = {}) {
  const { page = 1, limit = 20, search } = filters;

  return useQuery({
    queryKey: ["customers", { page, limit, search }],
    queryFn: async () => {
      const response = await api.get<PaginatedResponse<Customer>>("/customers", {
        page,
        limit,
        search,
      });
      return response;
    },
  });
}

export function useCustomer(id: string) {
  return useQuery({
    queryKey: ["customers", id],
    queryFn: async () => {
      const response = await api.get<{ data: Customer }>(`/customers/${id}`);
      return response.data;
    },
    enabled: !!id,
  });
}

export function useCreateCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateCustomerInput) => {
      const response = await api.post<{ data: Customer }>("/customers", input);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customers"] });
    },
  });
}

export function useUpdateCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...input }: UpdateCustomerInput & { id: string }) => {
      const response = await api.put<{ data: Customer }>(`/customers/${id}`, input);
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["customers"] });
      queryClient.setQueryData(["customers", data.id], data);
    },
  });
}

export function useDeleteCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/customers/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customers"] });
    },
  });
}
