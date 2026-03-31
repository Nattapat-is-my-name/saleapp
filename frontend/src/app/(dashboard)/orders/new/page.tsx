"use client";

import { useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, ShoppingCart, Trash2, CreditCard, Banknote, QrCode, Check, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
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
  const { data: productsData, isLoading: productsLoading } = useProducts({ limit: 100, isActive: true });

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
    return cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0);
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
    { id: "cash", label: "Cash", icon: Banknote },
    { id: "card", label: "Card", icon: CreditCard },
    { id: "qr", label: "QR Code", icon: QrCode },
  ] as const;

  return (
    <div className="flex h-[calc(100vh-8rem)] gap-6">
      {/* Product Selection */}
      <div className="flex-1 flex flex-col">
        <div className="flex items-center gap-4 mb-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/orders">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div className="flex-1">
            <h2 className="text-xl font-bold">Point of Sale</h2>
          </div>
        </div>

        <div className="relative mb-4">
          <Input
            placeholder="Search products by name or SKU..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-8"
          />
        </div>

        <div className="flex-1 overflow-y-auto pr-4">
          {productsLoading ? (
            <div className="flex items-center justify-center h-32">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            </div>
          ) : filteredProducts.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              No products found
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-3">
              {filteredProducts.map((product) => (
                <Card
                  key={product.id}
                  className="cursor-pointer hover:border-primary transition-colors"
                  onClick={() => addToCart(product)}
                >
                  <CardContent className="p-3">
                    <div className="font-medium text-sm truncate">{product.name}</div>
                    <div className="text-xs text-muted-foreground">{product.sku}</div>
                    <div className="mt-1 flex items-center justify-between">
                      <span className="font-bold text-sm">
                        {formatCurrency(product.price)}
                      </span>
                      <Badge variant={product.stock < 10 ? "destructive" : "secondary"}>
                        {product.stock} in stock
                      </Badge>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Cart & Checkout */}
      <Card className="w-[400px] flex flex-col">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ShoppingCart className="h-5 w-5" />
            Cart ({cart.length})
          </CardTitle>
          <CardDescription>
            {cart.length > 0
              ? `${cart.reduce((sum, i) => sum + i.quantity, 0)} items`
              : "Add items to start a sale"}
          </CardDescription>
        </CardHeader>

        <CardContent className="flex-1 flex flex-col">
          {cart.length === 0 ? (
            <div className="flex-1 flex items-center justify-center text-muted-foreground">
              Cart is empty
            </div>
          ) : (
            <div className="flex-1 overflow-y-auto mb-4 space-y-3 pr-2">
              {cart.map((item) => (
                <div key={item.product.id} className="flex gap-2 p-2 border rounded-lg">
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-sm truncate">{item.product.name}</div>
                    <div className="text-xs text-muted-foreground">
                      {formatCurrency(item.product.price)} each
                    </div>
                  </div>
                  <div className="flex items-center gap-1">
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-6 w-6"
                      onClick={() => updateQuantity(item.product.id, item.quantity - 1)}
                    >
                      -
                    </Button>
                    <span className="w-8 text-center text-sm">{item.quantity}</span>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-6 w-6"
                      onClick={() => updateQuantity(item.product.id, item.quantity + 1)}
                      disabled={item.quantity >= item.product.stock}
                    >
                      +
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-6 w-6 text-destructive"
                      onClick={() => removeFromCart(item.product.id)}
                    >
                      <Trash2 className="h-3 w-3" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}

          <div className="border-t pt-4 mt-auto space-y-4">
            <div className="flex justify-between text-lg font-bold">
              <span>Total</span>
              <span>{formatCurrency(subtotal)}</span>
            </div>

            <div className="space-y-2">
              <Label>Payment Method</Label>
              <div className="grid grid-cols-3 gap-2">
                {paymentMethods.map((method) => (
                  <Button
                    key={method.id}
                    variant={paymentMethod === method.id ? "default" : "outline"}
                    onClick={() => setPaymentMethod(method.id)}
                    className="flex flex-col h-auto py-2"
                  >
                    <method.icon className="h-4 w-4 mb-1" />
                    <span className="text-xs">{method.label}</span>
                  </Button>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="notes">Notes (optional)</Label>
              <Input
                id="notes"
                placeholder="Add notes..."
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
              />
            </div>

            <Button
              className="w-full"
              size="lg"
              disabled={cart.length === 0 || createOrder.isPending}
              onClick={handleCompleteSale}
            >
              {createOrder.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Processing...
                </>
              ) : (
                <>
                  <Check className="mr-2 h-4 w-4" />
                  Complete Sale
                </>
              )}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
