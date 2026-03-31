"use client";

import { useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Stack from "@mui/material/Stack";
import Grid from "@mui/material/Grid";
import Divider from "@mui/material/Divider";
import {
  ArrowBack,
  ShoppingCart,
  Delete,
  CreditCard,
  Money,
  QrCode2,
  Check,
} from "@mui/icons-material";
import { useProducts } from "@/hooks/use-products";
import { useCreateOrder } from "@/hooks/use-orders";
import { formatCurrency } from "@/lib/utils";
import type { Product } from "@/types/product";
import type { CreateOrderInput } from "@/types/order";

interface CartItem {
  product: Product;
  quantity: number;
}

type PaymentMethod = "cash" | "card" | "qr";

export default function NewOrderPage() {
  const router = useRouter();
  const createOrder = useCreateOrder();
  const { data: productsData, isLoading: productsLoading } = useProducts({
    limit: 100,
    isActive: true,
  });

  const [cart, setCart] = useState<CartItem[]>([]);
  const [search, setSearch] = useState("");
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>("cash");
  const [notes, setNotes] = useState("");
  const [customerId, setCustomerId] = useState("");

  const filteredProducts = useMemo(() => {
    if (!productsData?.data) return [];
    const searchLower = search.toLowerCase();
    return productsData.data.filter(
      (p) =>
        p.name.toLowerCase().includes(searchLower) ||
        p.sku.toLowerCase().includes(searchLower)
    );
  }, [productsData?.data, search]);

  const addToCart = (product: Product) => {
    setCart((prev) => {
      const existing = prev.find((item) => item.product.id === product.id);
      if (existing) {
        if (existing.quantity >= product.stock) return prev;
        return prev.map((item) =>
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item
        );
      }
      return [...prev, { product, quantity: 1 }];
    });
  };

  const updateQuantity = (productId: string, quantity: number) => {
    if (quantity < 1) {
      removeFromCart(productId);
      return;
    }
    setCart((prev) =>
      prev.map((item) =>
        item.product.id === productId ? { ...item, quantity } : item
      )
    );
  };

  const removeFromCart = (productId: string) => {
    setCart((prev) => prev.filter((item) => item.product.id !== productId));
  };

  const subtotal = useMemo(() => {
    return cart.reduce(
      (sum, item) => sum + item.product.price * item.quantity,
      0
    );
  }, [cart]);

  const handleCompleteSale = async () => {
    if (cart.length === 0) return;

    const orderInput: CreateOrderInput = {
      paymentMethod,
      notes,
      items: cart.map((item) => ({
        productId: item.product.id,
        quantity: item.quantity,
        discount: 0,
      })),
    };

    if (customerId.trim()) {
      orderInput.customerId = customerId;
    }

    try {
      await createOrder.mutateAsync(orderInput);
      router.push("/orders");
    } catch (error) {
      console.error("Failed to create order:", error);
    }
  };

  const paymentMethods = [
    { id: "cash" as PaymentMethod, label: "Cash", icon: Money },
    { id: "card" as PaymentMethod, label: "Card", icon: CreditCard },
    { id: "qr" as PaymentMethod, label: "QR Code", icon: QrCode2 },
  ];

  const totalItems = cart.reduce((sum, i) => sum + i.quantity, 0);

  return (
    <Box sx={{ display: "flex", gap: 3, height: "calc(100vh - 120px)" }}>
      {/* Product Selection */}
      <Box sx={{ flex: 1, display: "flex", flexDirection: "column" }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
          <Button
            variant="text"
            startIcon={<ArrowBack />}
            component={Link}
            href="/orders"
            sx={{ color: "#64748b" }}
          >
            Back
          </Button>
          <Typography variant="h5" sx={{ fontWeight: 700 }}>
            Point of Sale
          </Typography>
        </Box>

        <TextField
          placeholder="Search products by name or SKU..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          size="small"
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ShoppingCart sx={{ color: "#64748b" }} />
                </InputAdornment>
              ),
            },
          }}
          sx={{ mb: 2 }}
        />

        <Box sx={{ flex: 1, overflowY: "auto", pr: 1 }}>
          {productsLoading ? (
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                height: 200,
              }}
            >
              <CircularProgress size={24} />
            </Box>
          ) : filteredProducts.length === 0 ? (
            <Box sx={{ textAlign: "center", py: 8, color: "#64748b" }}>
              No products found
            </Box>
          ) : (
            <Grid container spacing={2}>
              {filteredProducts.map((product) => (
                <Grid key={product.id} size={{ xs: 6 }}>
                  <Card
                    sx={{
                      cursor: "pointer",
                      "&:hover": { borderColor: "#3b82f6" },
                      transition: "border-color 0.2s",
                    }}
                    onClick={() => addToCart(product)}
                  >
                    <CardContent sx={{ p: 2, "&:last-child": { pb: 2 } }}>
                      <Typography variant="body2" sx={{ fontWeight: 600, mb: 0.25 }}>
                        {product.name}
                      </Typography>
                      <Typography variant="caption" sx={{ color: "#64748b" }}>
                        {product.sku}
                      </Typography>
                      <Box
                        sx={{
                          display: "flex",
                          justifyContent: "space-between",
                          alignItems: "center",
                          mt: 1,
                        }}
                      >
                        <Typography variant="body2" sx={{ fontWeight: 700 }}>
                          {formatCurrency(product.price)}
                        </Typography>
                        <Chip
                          label={`${product.stock} in stock`}
                          size="small"
                          sx={{
                            backgroundColor:
                              product.stock < 10 ? "#fef2f2" : "#f1f5f9",
                            color: product.stock < 10 ? "#dc2626" : "#64748b",
                            fontSize: "0.65rem",
                            height: 20,
                          }}
                        />
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          )}
        </Box>
      </Box>

      {/* Cart & Checkout */}
      <Card sx={{ width: 400, display: "flex", flexDirection: "column" }}>
        <Box sx={{ p: 2, borderBottom: "1px solid #e2e8f0" }}>
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <ShoppingCart sx={{ color: "#3b82f6" }} />
            <Typography variant="h6" sx={{ fontWeight: 700 }}>
              Cart ({cart.length})
            </Typography>
          </Box>
          <Typography variant="caption" sx={{ color: "#64748b" }}>
            {totalItems > 0 ? `${totalItems} items` : "Add items to start a sale"}
          </Typography>
        </Box>

        <CardContent sx={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
          {cart.length === 0 ? (
            <Box
              sx={{
                flex: 1,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "#64748b",
              }}
            >
              Cart is empty
            </Box>
          ) : (
            <Box sx={{ flex: 1, overflowY: "auto", mb: 2, pr: 1 }}>
              <Stack spacing={1.5}>
                {cart.map((item) => (
                  <Box
                    key={item.product.id}
                    sx={{
                      display: "flex",
                      gap: 1.5,
                      p: 1.5,
                      borderRadius: "8px",
                      border: "1px solid #e2e8f0",
                    }}
                  >
                    <Box sx={{ flex: 1, minWidth: 0 }}>
                      <Typography
                        variant="body2"
                        sx={{ fontWeight: 600, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}
                      >
                        {item.product.name}
                      </Typography>
                      <Typography variant="caption" sx={{ color: "#64748b" }}>
                        {formatCurrency(item.product.price)} each
                      </Typography>
                    </Box>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
                      <Button size="small" variant="outlined" disableElevation
                        onClick={() =>
                          updateQuantity(item.product.id, item.quantity - 1)
                        }
                        sx={{
                          width: 24,
                          height: 24,
                          borderColor: "#e2e8f0",
                          color: "#64748b",
                        }}
                      >
                        <Typography variant="body2">−</Typography>
                      </Button>
                      <Typography
                        variant="body2"
                        sx={{ width: 24, textAlign: "center", fontWeight: 600 }}
                      >
                        {item.quantity}
                      </Typography>
                      <Button size="small" variant="outlined" disableElevation
                        onClick={() =>
                          updateQuantity(item.product.id, item.quantity + 1)
                        }
                        disabled={item.quantity >= item.product.stock}
                        sx={{
                          width: 24,
                          height: 24,
                          borderColor: "#e2e8f0",
                          color: "#64748b",
                        }}
                      >
                        <Typography variant="body2">+</Typography>
                      </Button>
                      <IconButton
                        size="small"
                        onClick={() => removeFromCart(item.product.id)}
                        sx={{
                          width: 24,
                          height: 24,
                          color: "#dc2626",
                          "&:hover": { backgroundColor: "#fef2f2" },
                        }}
                      >
                        <Delete sx={{ fontSize: 16 }} />
                      </IconButton>
                    </Box>
                  </Box>
                ))}
              </Stack>
            </Box>
          )}

          <Divider sx={{ my: 2 }} />

          <Box>
            <Box
              sx={{
                display: "flex",
                justifyContent: "space-between",
                mb: 2,
              }}
            >
              <Typography variant="h6" sx={{ fontWeight: 700 }}>
                Total
              </Typography>
              <Typography variant="h6" sx={{ fontWeight: 700, color: "#3b82f6" }}>
                {formatCurrency(subtotal)}
              </Typography>
            </Box>

            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" sx={{ fontWeight: 500, mb: 1 }}>
                Payment Method
              </Typography>
              <Grid container spacing={1}>
                {paymentMethods.map((method) => (
                  <Grid key={method.id} size={4}>
                    <Button
                      variant={paymentMethod === method.id ? "contained" : "outlined"}
                      onClick={() => setPaymentMethod(method.id)}
                      fullWidth
                      sx={{
                        display: "flex",
                        flexDirection: "column",
                        py: 1.5,
                        gap: 0.5,
                        fontSize: "0.75rem",
                      }}
                    >
                      <method.icon sx={{ fontSize: 20 }} />
                      {method.label}
                    </Button>
                  </Grid>
                ))}
              </Grid>
            </Box>

            <TextField
              label="Notes (optional)"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              size="small"
              fullWidth
              sx={{ mb: 2 }}
            />

            <Button
              variant="contained"
              fullWidth
              size="large"
              disabled={cart.length === 0 || createOrder.isPending}
              onClick={handleCompleteSale}
              startIcon={
                createOrder.isPending ? (
                  <CircularProgress size={18} color="inherit" />
                ) : (
                  <Check />
                )
              }
            >
              {createOrder.isPending ? "Processing..." : "Complete Sale"}
            </Button>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
