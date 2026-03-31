"use client";

import { useState } from "react";
import Link from "next/link";
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import Stack from "@mui/material/Stack";
import { Add, Search, Edit, Delete } from "@mui/icons-material";
import Chip from "@mui/material/Chip";
import { useProducts, useDeleteProduct } from "@/hooks/use-products";
import { formatCurrency } from "@/lib/utils";
import type { Product } from "@/types/product";

export default function ProductsPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");

  const { data, isLoading } = useProducts({
    page,
    limit: 20,
    search: debouncedSearch || undefined,
  });

  const deleteProduct = useDeleteProduct();

  const handleSearchChange = (value: string) => {
    setSearch(value);
    const timer = setTimeout(() => {
      setDebouncedSearch(value);
      setPage(1);
    }, 300);
    return () => clearTimeout(timer);
  };

  const handleDelete = async (id: string) => {
    if (window.confirm("Are you sure you want to delete this product?")) {
      await deleteProduct.mutateAsync(id);
    }
  };

  return (
    <Box>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 3 }}>
        <Box>
          <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
            Products
          </Typography>
          <Typography sx={{ color: "#64748b" }}>
            Manage your product catalog
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<Add />}
          component={Link}
          href="/products/new"
        >
          Add Product
        </Button>
      </Box>

      <Card>
        <Box sx={{ p: 2 }}>
          <TextField
            placeholder="Search products..."
            value={search}
            onChange={(e) => handleSearchChange(e.target.value)}
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
              Loading products...
            </Box>
          ) : !data?.data.length ? (
            <Box sx={{ py: 8, textAlign: "center", color: "#64748b" }}>
              No products found. Create your first product to get started.
            </Box>
          ) : (
            <>
              {/* Table Header */}
              <Box
                sx={{
                  display: "grid",
                  gridTemplateColumns: "100px 1fr 120px 100px 100px 140px",
                  gap: 2,
                  px: 2,
                  py: 1.5,
                  backgroundColor: "#f8fafc",
                  borderBottom: "1px solid #e2e8f0",
                }}
              >
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  SKU
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Name
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Price
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Stock
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b" }}>
                  Status
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: "#64748b", textAlign: "right" }}>
                  Actions
                </Typography>
              </Box>

              {/* Table Rows */}
              {data.data.map((product: Product) => (
                <Box
                  key={product.id}
                  sx={{
                    display: "grid",
                    gridTemplateColumns: "100px 1fr 120px 100px 100px 140px",
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
                    {product.sku}
                  </Typography>
                  <Typography variant="body2">{product.name}</Typography>
                  <Typography variant="body2">{formatCurrency(product.price)}</Typography>
                  <Box>
                    <Chip
                      label={product.stock}
                      color={product.stock < 10 ? "warning" : "secondary"}
                      size="small"
                    />
                  </Box>
                  <Box>
                    <Chip
                      label={product.isActive ? "Active" : "Inactive"}
                      color={product.isActive ? "success" : "secondary"}
                      size="small"
                    />
                  </Box>
                  <Box sx={{ display: "flex", justifyContent: "flex-end", gap: 0.5 }}>
                    <IconButton
                      size="small"
                      component={Link}
                      href={`/products/${product.id}`}
                      sx={{ color: "#64748b" }}
                    >
                      <Edit sx={{ fontSize: 18 }} />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => handleDelete(product.id)}
                      disabled={deleteProduct.isPending}
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
                    {Math.min(page * 20, data.meta.total)} of {data.meta.total} products
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
