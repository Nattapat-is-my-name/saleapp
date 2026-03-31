"use client";

import { useState } from "react";
import Link from "next/link";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import Stack from "@mui/material/Stack";
import Grid from "@mui/material/Grid";
import Chip from "@mui/material/Chip";
import { Add, Search, Visibility, Refresh } from "@mui/icons-material";
import { useOrders, useUpdateOrderStatus } from "@/hooks/use-orders";
import { formatCurrency, formatDateTime } from "@/lib/utils";
import type { Order, OrderStatus } from "@/types/order";

const statusColors: Record<OrderStatus, "default" | "success" | "warning" | "error" | "info" | "secondary"> = {
  pending: "warning",
  completed: "success",
  cancelled: "error",
  refunded: "default",
};

export default function OrdersPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");

  const { data, isLoading } = useOrders({
    page,
    limit: 20,
  });

  const updateStatus = useUpdateOrderStatus();

  const handleStatusUpdate = async (id: string, status: OrderStatus) => {
    await updateStatus.mutateAsync({ id, status });
  };

  return (
    <Box>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 3 }}>
        <Box>
          <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
            Orders
          </Typography>
          <Typography sx={{ color: "#64748b" }}>
            View and manage customer orders
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<Add />}
          component={Link}
          href="/orders/new"
        >
          New Order
        </Button>
      </Box>

      <Card>
        <Box sx={{ p: 2 }}>
          <TextField
            placeholder="Search orders..."
            value={search}
            size="small"
            slotProps={{
              input: {
                startAdornment: (
                  <InputAdornment position="start">
                    <Search sx={{ color: "#64748b" }} />
                  </InputAdornment>
                ),
              },
            }}
            sx={{ width: 300 }}
          />
        </Box>
        <Box>
          {isLoading ? (
            <Box sx={{ py: 8, textAlign: "center", color: "#64748b" }}>
              Loading orders...
            </Box>
          ) : !data?.data.length ? (
            <Box sx={{ py: 8, textAlign: "center", color: "#64748b" }}>
              No orders found. Create your first order to get started.
            </Box>
          ) : (
            <>
              {/* Table Header */}
              <Box
                sx={{
                  display: "grid",
                  gridTemplateColumns: "140px 1fr 180px 100px 120px 120px",
                  gap: 2,
                  px: 2,
                  py: 1.5,
                  backgroundColor: "#f8fafc",
                  borderBottom: "1px solid #e2e8f0",
                }}
              >
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Order #
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Customer
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Date
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Status
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Total
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b", textAlign: "right" }}>
                  Actions
                </Typography>
              </Box>

              {/* Table Rows */}
              {data.data.map((order: Order) => (
                <Box
                  key={order.id}
                  sx={{
                    display: "grid",
                    gridTemplateColumns: "140px 1fr 180px 100px 120px 120px",
                    gap: 2,
                    px: 2,
                    py: 1.5,
                    borderBottom: "1px solid #f1f5f9",
                    alignItems: "center",
                    "&:hover": { backgroundColor: "#f8fafc" },
                    "&:last-child": { borderBottom: 0 },
                  }}
                >
                  <Typography variant="body2" sx={{ fontWeight: 500 }}>
                    {order.orderNumber}
                  </Typography>
                  <Typography variant="body2">
                    {order.customer
                      ? `${order.customer.firstName} ${order.customer.lastName}`
                      : "Walk-in Customer"}
                  </Typography>
                  <Typography variant="body2" sx={{ color: "#64748b", fontSize: "0.8rem" }}>
                    {formatDateTime(order.createdAt)}
                  </Typography>
                  <Box>
                    <Chip
                      label={order.status}
                      color={statusColors[order.status]}
                      size="small"
                    />
                  </Box>
                  <Typography variant="body2" sx={{ fontWeight: 600 }}>
                    {formatCurrency(order.total)}
                  </Typography>
                  <Box sx={{ display: "flex", justifyContent: "flex-end", gap: 0.5 }}>
                    <IconButton
                      size="small"
                      component={Link}
                      href={`/orders/${order.id}`}
                      sx={{ color: "#64748b" }}
                    >
                      <Visibility sx={{ fontSize: 18 }} />
                    </IconButton>
                    {order.status === "pending" && (
                      <IconButton
                        size="small"
                        onClick={() => handleStatusUpdate(order.id, "completed")}
                        disabled={updateStatus.isPending}
                        sx={{ color: "#10b981" }}
                      >
                        <Refresh sx={{ fontSize: 18 }} />
                      </IconButton>
                    )}
                  </Box>
                </Box>
              ))}

              {data.meta && data.meta.totalPages > 1 && (
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    px: 2,
                    py: 1.5,
                    borderTop: "1px solid #e2e8f0",
                  }}
                >
                  <Typography variant="body2" sx={{ color: "#64748b" }}>
                    Showing {(page - 1) * 20 + 1} to{" "}
                    {Math.min(page * 20, data.meta.total)} of {data.meta.total} orders
                  </Typography>
                  <Stack direction="row" spacing={1}>
                    <Button
                      variant="outlined"
                      size="small"
                      onClick={() => setPage(page - 1)}
                      disabled={page === 1}
                    >
                      Previous
                    </Button>
                    <Button
                      variant="outlined"
                      size="small"
                      onClick={() => setPage(page + 1)}
                      disabled={page >= data.meta.totalPages}
                    >
                      Next
                    </Button>
                  </Stack>
                </Box>
              )}
            </>
          )}
        </Box>
      </Card>
    </Box>
  );
}
