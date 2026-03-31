"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import Stack from "@mui/material/Stack";
import Chip from "@mui/material/Chip";
import { Add, Search, Edit, Delete } from "@mui/icons-material";
import { useCustomers, useDeleteCustomer } from "@/hooks/use-customers";
import type { Customer } from "@/types/customer";

export default function CustomersPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");

  const { data, isLoading } = useCustomers({
    page,
    limit: 20,
    search: debouncedSearch || undefined,
  });

  const deleteCustomer = useDeleteCustomer();

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
      setPage(1);
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  const handleDelete = async (id: string) => {
    if (window.confirm("Are you sure you want to delete this customer?")) {
      await deleteCustomer.mutateAsync(id);
    }
  };

  return (
    <Box>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 3 }}>
        <Box>
          <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
            Customers
          </Typography>
          <Typography sx={{ color: "#64748b" }}>
            Manage your customer database
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<Add />}
          component={Link}
          href="/customers/new"
        >
          Add Customer
        </Button>
      </Box>

      <Card>
        <Box sx={{ p: 2 }}>
          <TextField
            placeholder="Search customers..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
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
        <Box sx={{ borderTop: "1px solid #e2e8f0" }}>
          {isLoading ? (
            <Box sx={{ py: 8, textAlign: "center", color: "#64748b" }}>
              Loading customers...
            </Box>
          ) : !data?.data.length ? (
            <Box sx={{ py: 8, textAlign: "center", color: "#64748b" }}>
              No customers found. Add your first customer to get started.
            </Box>
          ) : (
            <>
              {/* Table Header */}
              <Box
                sx={{
                  display: "grid",
                  gridTemplateColumns: "1fr 1fr 150px 1fr 120px",
                  gap: 2,
                  px: 2,
                  py: 1.5,
                  backgroundColor: "#f8fafc",
                  borderBottom: "1px solid #e2e8f0",
                }}
              >
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Name
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Email
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Phone
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Address
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b", textAlign: "right" }}>
                  Actions
                </Typography>
              </Box>

              {/* Table Rows */}
              {data.data.map((customer: Customer) => (
                <Box
                  key={customer.id}
                  sx={{
                    display: "grid",
                    gridTemplateColumns: "1fr 1fr 150px 1fr 120px",
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
                    {customer.firstName} {customer.lastName}
                  </Typography>
                  <Typography variant="body2">{customer.email}</Typography>
                  <Typography variant="body2" sx={{ color: "#64748b" }}>
                    {customer.phone || "—"}
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{
                      color: "#64748b",
                      overflow: "hidden",
                      textOverflow: "ellipsis",
                      whiteSpace: "nowrap",
                    }}
                  >
                    {customer.address || "—"}
                  </Typography>
                  <Box sx={{ display: "flex", justifyContent: "flex-end", gap: 0.5 }}>
                    <IconButton
                      size="small"
                      component={Link}
                      href={`/customers/${customer.id}`}
                      sx={{ color: "#64748b" }}
                    >
                      <Edit sx={{ fontSize: 18 }} />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => handleDelete(customer.id)}
                      disabled={deleteCustomer.isPending}
                      sx={{ color: "#dc2626" }}
                    >
                      <Delete sx={{ fontSize: 18 }} />
                    </IconButton>
                  </Box>
                </Box>
              ))}

              {/* Pagination */}
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
                    {Math.min(page * 20, data.meta.total)} of {data.meta.total} customers
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
