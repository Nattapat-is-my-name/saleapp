"use client";

import { useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Grid from "@mui/material/Grid";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import Stack from "@mui/material/Stack";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import { ArrowBack } from "@mui/icons-material";
import { useProduct, useUpdateProduct } from "@/hooks/use-products";
import type { UpdateProductInput } from "@/types/product";

export default function EditProductPage() {
  const router = useRouter();
  const params = useParams();
  const productId = params.id as string;

  const { data: product, isLoading, error } = useProduct(productId);
  const updateProduct = useUpdateProduct();

  const [formData, setFormData] = useState<UpdateProductInput & { sku: string }>({
    sku: "",
    name: "",
    description: "",
    price: 0,
    cost: 0,
    stock: 0,
    isActive: true,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [initialized, setInitialized] = useState(false);

  // Initialize form when product loads
  if (product && !initialized) {
    setFormData({
      sku: product.sku || "",
      name: product.name || "",
      description: product.description || "",
      price: product.price || 0,
      cost: product.cost || 0,
      stock: product.stock || 0,
      isActive: product.isActive ?? true,
    });
    setInitialized(true);
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.sku?.trim()) {
      newErrors.sku = "SKU is required";
    }
    if (!formData.name?.trim()) {
      newErrors.name = "Name is required";
    }
    if (!formData.price || formData.price <= 0) {
      newErrors.price = "Price must be greater than 0";
    }
    if (formData.stock !== undefined && formData.stock < 0) {
      newErrors.stock = "Stock cannot be negative";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    try {
      await updateProduct.mutateAsync({ id: productId, ...formData });
      router.push("/products");
    } catch (err) {
      console.error("Failed to update product:", err);
    }
  };

  const handleChange = (field: keyof typeof formData, value: string | number | boolean) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  if (isLoading) {
    return (
      <Box sx={{ maxWidth: 720 }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
          <Button
            variant="text"
            startIcon={<ArrowBack />}
            component={Link}
            href="/products"
            sx={{ color: "#64748b" }}
          >
            Back
          </Button>
        </Box>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
          Edit Product
        </Typography>
        <Typography sx={{ color: "#64748b", mb: 3 }}>
          Loading product details...
        </Typography>
        <Card>
          <CardContent>
            <Box sx={{ py: 4, textAlign: "center" }}>
              <CircularProgress />
            </Box>
          </CardContent>
        </Card>
      </Box>
    );
  }

  if (error || !product) {
    return (
      <Box sx={{ maxWidth: 720 }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
          <Button
            variant="text"
            startIcon={<ArrowBack />}
            component={Link}
            href="/products"
            sx={{ color: "#64748b" }}
          >
            Back
          </Button>
        </Box>
        <Alert severity="error">Product not found.</Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 720 }}>
      <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
        <Button
          variant="text"
          startIcon={<ArrowBack />}
          component={Link}
          href="/products"
          sx={{ color: "#64748b" }}
        >
          Back
        </Button>
      </Box>

      <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
        Edit Product
      </Typography>
      <Typography sx={{ color: "#64748b", mb: 3 }}>
        Update product information
      </Typography>

      <Card>
        <CardContent>
          <Typography variant="h6" sx={{ mb: 1 }}>
            Product Information
          </Typography>
          <Typography variant="body2" sx={{ color: "#64748b", mb: 3 }}>
            Update the details below
          </Typography>

          <Box component="form" onSubmit={handleSubmit}>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  label="SKU"
                  value={formData.sku}
                  onChange={(e) => handleChange("sku", e.target.value)}
                  error={!!errors.sku}
                  helperText={errors.sku}
                  fullWidth
                  size="small"
                  required
                  placeholder="e.g., PROD-001"
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  label="Product Name"
                  value={formData.name}
                  onChange={(e) => handleChange("name", e.target.value)}
                  error={!!errors.name}
                  helperText={errors.name}
                  fullWidth
                  size="small"
                  required
                  placeholder="e.g., Classic T-Shirt"
                />
              </Grid>

              <Grid size={{ xs: 12 }}>
                <TextField
                  label="Description"
                  value={formData.description}
                  onChange={(e) => handleChange("description", e.target.value)}
                  fullWidth
                  size="small"
                  placeholder="Product description (optional)"
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 4 }}>
                <TextField
                  label="Price"
                  type="number"
                  value={formData.price}
                  onChange={(e) => handleChange("price", parseFloat(e.target.value) || 0)}
                  error={!!errors.price}
                  helperText={errors.price}
                  fullWidth
                  size="small"
                  required
                  slotProps={{
                    input: {
                      startAdornment: <InputAdornment position="start">$</InputAdornment>,
                    },
                  }}
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 4 }}>
                <TextField
                  label="Cost"
                  type="number"
                  value={formData.cost}
                  onChange={(e) => handleChange("cost", parseFloat(e.target.value) || 0)}
                  fullWidth
                  size="small"
                  slotProps={{
                    input: {
                      startAdornment: <InputAdornment position="start">$</InputAdornment>,
                    },
                  }}
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 4 }}>
                <TextField
                  label="Stock Quantity"
                  type="number"
                  value={formData.stock}
                  onChange={(e) => handleChange("stock", parseInt(e.target.value) || 0)}
                  error={!!errors.stock}
                  helperText={errors.stock}
                  fullWidth
                  size="small"
                />
              </Grid>
            </Grid>

            <Stack direction="row" spacing={2} sx={{ mt: 4 }}>
              <Button
                type="submit"
                variant="contained"
                disabled={updateProduct.isPending}
                fullWidth
                startIcon={updateProduct.isPending ? <CircularProgress size={16} /> : null}
              >
                {updateProduct.isPending ? "Saving..." : "Save Changes"}
              </Button>
              <Button
                type="button"
                variant="outlined"
                onClick={() => router.push("/products")}
                fullWidth
              >
                Cancel
              </Button>
            </Stack>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
