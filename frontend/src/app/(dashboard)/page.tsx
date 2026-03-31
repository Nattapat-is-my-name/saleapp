"use client";

import { DollarSign, Package, ShoppingCart, Users } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useOrders } from "@/hooks/use-orders";
import { useProducts } from "@/hooks/use-products";
import { useCustomers } from "@/hooks/use-customers";
import { formatCurrency } from "@/lib/utils";

export default function DashboardPage() {
  const { data: ordersData } = useOrders({ limit: 1 });
  const { data: productsData } = useProducts({ limit: 1 });
  const { data: customersData } = useCustomers({ limit: 1 });

  const stats = [
    {
      title: "Total Revenue",
      value: formatCurrency(ordersData?.meta?.total ?? 0),
      icon: DollarSign,
      description: "All time revenue",
    },
    {
      title: "Products",
      value: productsData?.meta?.total ?? 0,
      icon: Package,
      description: "Active products",
    },
    {
      title: "Orders",
      value: ordersData?.meta?.total ?? 0,
      icon: ShoppingCart,
      description: "Total orders",
    },
    {
      title: "Customers",
      value: customersData?.meta?.total ?? 0,
      icon: Users,
      description: "Registered customers",
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
        <p className="text-muted-foreground">
          Welcome back! Here&apos;s an overview of your store.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">
                {stat.description}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Recent Orders</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Recent orders will appear here once you have sales data.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Low Stock Alert</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Products with low stock will appear here.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
