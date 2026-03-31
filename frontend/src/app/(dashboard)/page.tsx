"use client";

import { DollarSign, Package, ShoppingCart, AlertTriangle, TrendingUp, ArrowUpRight } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useDashboard } from "@/hooks/use-dashboard";
import { formatCurrency } from "@/lib/utils";
import { 
  BarChart, 
  Bar, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from "recharts";

const COLORS = ["#3b82f6", "#8b5cf6", "#ec4899", "#f59e0b", "#10b981"];

export default function DashboardPage() {
  const { data: dashboard, isLoading } = useDashboard();

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
          <p className="text-muted-foreground">Loading your store data...</p>
        </div>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[1, 2, 3, 4].map((i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader className="pb-2">
                <div className="h-4 bg-muted rounded w-20" />
              </CardHeader>
              <CardContent>
                <div className="h-8 bg-muted rounded w-24" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  const stats = [
    {
      title: "Today's Revenue",
      value: formatCurrency(dashboard?.today_revenue ?? 0),
      icon: DollarSign,
      description: "Completed orders today",
      change: "+12%",
      changeType: "positive",
    },
    {
      title: "Today's Orders",
      value: dashboard?.today_orders ?? 0,
      icon: ShoppingCart,
      description: "Orders completed",
    },
    {
      title: "Low Stock Items",
      value: dashboard?.low_stock?.length ?? 0,
      icon: AlertTriangle,
      description: "Products need restocking",
      variant: (dashboard?.low_stock?.length ?? 0) > 0 ? "destructive" : "default",
    },
    {
      title: "Top Products",
      value: dashboard?.top_products?.length ?? 0,
      icon: TrendingUp,
      description: "Best sellers today",
    },
  ];

  const topProductsChart = (dashboard?.top_products ?? []).map((tp) => ({
    name: tp.product.name.length > 15 ? tp.product.name.substring(0, 15) + "..." : tp.product.name,
    quantity: tp.quantity_sold,
    revenue: tp.product.price * tp.quantity_sold,
  }));

  const lowStockChart = (dashboard?.low_stock ?? []).slice(0, 5).map((p) => ({
    name: p.name.length > 12 ? p.name.substring(0, 12) + "..." : p.name,
    stock: p.stock,
  }));

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
        <p className="text-muted-foreground">
          Welcome back! Here&apos;s what&apos;s happening in your store today.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">{stat.description}</p>
              {stat.change && (
                <div className={`flex items-center mt-1 text-xs ${stat.changeType === "positive" ? "text-green-600" : "text-red-600"}`}>
                  <ArrowUpRight className="h-3 w-3 mr-1" />
                  {stat.change}
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Top Selling Products</CardTitle>
          </CardHeader>
          <CardContent>
            {topProductsChart.length > 0 ? (
              <div className="h-[300px]">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={topProductsChart}>
                    <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                    <XAxis 
                      dataKey="name" 
                      tick={{ fontSize: 12 }}
                      tickLine={false}
                      axisLine={false}
                    />
                    <YAxis 
                      tick={{ fontSize: 12 }}
                      tickLine={false}
                      axisLine={false}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: "hsl(var(--card))", 
                        border: "1px solid hsl(var(--border))",
                        borderRadius: "8px"
                      }}
                    />
                    <Bar dataKey="quantity" fill="#3b82f6" radius={[4, 4, 0, 0]} name="Units Sold" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="h-[300px] flex items-center justify-center text-muted-foreground">
                No sales data for today yet
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Low Stock Alert</CardTitle>
          </CardHeader>
          <CardContent>
            {lowStockChart.length > 0 ? (
              <div className="h-[300px]">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={lowStockChart}
                      cx="50%"
                      cy="50%"
                      innerRadius={60}
                      outerRadius={100}
                      paddingAngle={2}
                      dataKey="stock"
                      nameKey="name"
                    >
                      {lowStockChart.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: "hsl(var(--card))", 
                        border: "1px solid hsl(var(--border))",
                        borderRadius: "8px"
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="h-[300px] flex items-center justify-center text-muted-foreground">
                All products are well-stocked
              </div>
            )}
            <div className="mt-4 space-y-2">
              {(dashboard?.low_stock ?? []).slice(0, 5).map((product) => (
                <div key={product.id} className="flex items-center justify-between text-sm">
                  <span className="truncate max-w-[150px]">{product.name}</span>
                  <Badge variant="destructive">{product.stock} left</Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-1">
        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Orders and activity will appear here in real-time.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
