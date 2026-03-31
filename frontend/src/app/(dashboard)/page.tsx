"use client";

import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Skeleton from "@mui/material/Skeleton";
import Stack from "@mui/material/Stack";
import {
  Money,
  ShoppingBag,
  Warning,
  ShowChart,
} from "@mui/icons-material";
import Chip from "@mui/material/Chip";
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
  Cell,
} from "recharts";

const COLORS = ["#3b82f6", "#8b5cf6", "#ec4899", "#f59e0b", "#10b981"];

export default function DashboardPage() {
  const { data: dashboard, isLoading } = useDashboard();

  if (isLoading) {
    return (
      <Box>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
          Dashboard
        </Typography>
        <Typography sx={{ color: "#64748b", mb: 4 }}>
          Loading your store data...
        </Typography>
        <Grid container spacing={3}>
          {[1, 2, 3, 4].map((i) => (
            <Grid key={i} size={{ xs: 12, sm: 6, lg: 3 }}>
              <Card>
                <CardContent>
                  <Skeleton width={100} height={20} />
                  <Skeleton width={80} height={36} sx={{ mt: 1 }} />
                  <Skeleton width={120} height={16} sx={{ mt: 1 }} />
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Box>
    );
  }

  const stats = [
    {
      title: "Today's Revenue",
      value: formatCurrency(dashboard?.today_revenue ?? 0),
      icon: Money,
      description: "Completed orders today",
      change: "+12%",
    },
    {
      title: "Today's Orders",
      value: dashboard?.today_orders ?? 0,
      icon: ShoppingBag,
      description: "Orders completed",
    },
    {
      title: "Low Stock Items",
      value: dashboard?.low_stock?.length ?? 0,
      icon: Warning,
      description: "Products need restocking",
      warning: (dashboard?.low_stock?.length ?? 0) > 0,
    },
    {
      title: "Top Products",
      value: dashboard?.top_products?.length ?? 0,
      icon: ShowChart,
      description: "Best sellers today",
    },
  ];

  const topProductsChart = (dashboard?.top_products ?? []).map((tp: { product: { name: string; price: number }; quantity_sold: number }) => ({
    name: tp.product.name.length > 15 ? tp.product.name.substring(0, 15) + "..." : tp.product.name,
    quantity: tp.quantity_sold,
    revenue: tp.product.price * tp.quantity_sold,
  }));

  const lowStockChart = (dashboard?.low_stock ?? []).slice(0, 5).map((p: { name: string; stock: number }) => ({
    name: p.name.length > 12 ? p.name.substring(0, 12) + "..." : p.name,
    stock: p.stock,
  }));

  return (
    <Box>
      <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
        Dashboard
      </Typography>
      <Typography sx={{ color: "#64748b", mb: 4 }}>
        Welcome back! Here&apos;s what&apos;s happening in your store today.
      </Typography>

      {/* Stat Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {stats.map((stat) => (
          <Grid key={stat.title} size={{ xs: 12, sm: 6, lg: 3 }}>
            <Card>
              <CardContent>
                <Box sx={{ display: "flex", alignItems: "center", justifyContent: "space-between", mb: 1 }}>
                  <Typography variant="body2" sx={{ fontWeight: 500, color: "#64748b" }}>
                    {stat.title}
                  </Typography>
                    <Box
                    sx={{
                      width: 32,
                      height: 32,
                      borderRadius: "8px",
                      backgroundColor: stat.warning ? "#fef3c7" : "#f1f5f9",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                    }}
                  >
                    <stat.icon
                      sx={{
                        fontSize: 18,
                        color: stat.warning ? "#d97706" : "#64748b",
                      }}
                    />
                  </Box>
                </Box>
                <Typography variant="h5" sx={{ fontWeight: 700 }}>
                  {stat.value}
                </Typography>
                <Typography variant="caption" sx={{ color: "#64748b" }}>
                  {stat.description}
                </Typography>
                {stat.change && (
                  <Box sx={{ display: "flex", alignItems: "center", mt: 0.5 }}>
                    <Chip
                      label={stat.change}
                      sx={{
                        backgroundColor: "#f0fdf4",
                        color: "#16a34a",
                        fontSize: "0.7rem",
                        height: 20,
                      }}
                    />
                  </Box>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Charts */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid size={{ xs: 12, lg: 8 }}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 2 }}>
                Top Selling Products
              </Typography>
              {topProductsChart.length > 0 ? (
                <Box sx={{ height: 300 }}>
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={topProductsChart}>
                      <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
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
                          backgroundColor: "#fff",
                          border: "1px solid #e2e8f0",
                          borderRadius: 8,
                        }}
                      />
                      <Bar
                        dataKey="quantity"
                        fill="#3b82f6"
                        radius={[4, 4, 0, 0]}
                        name="Units Sold"
                      />
                    </BarChart>
                  </ResponsiveContainer>
                </Box>
              ) : (
                <Box
                  sx={{
                    height: 300,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "#64748b",
                  }}
                >
                  No sales data for today yet
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>

        <Grid size={{ xs: 12, lg: 4 }}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 2 }}>
                Low Stock Alert
              </Typography>
              {lowStockChart.length > 0 ? (
                <>
                  <Box sx={{ height: 200 }}>
                    <ResponsiveContainer width="100%" height="100%">
                      <PieChart>
                        <Pie
                          data={lowStockChart}
                          cx="50%"
                          cy="50%"
                          innerRadius={50}
                          outerRadius={80}
                          paddingAngle={2}
                          dataKey="stock"
                          nameKey="name"
                        >
                          {lowStockChart.map((_: unknown, index: number) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                          ))}
                        </Pie>
                        <Tooltip
                          contentStyle={{
                            backgroundColor: "#fff",
                            border: "1px solid #e2e8f0",
                            borderRadius: 8,
                          }}
                        />
                      </PieChart>
                    </ResponsiveContainer>
                  </Box>
                  <Stack spacing={1} sx={{ mt: 2 }}>
                    {(dashboard?.low_stock ?? []).slice(0, 5).map((product: { id: string; name: string; stock: number }) => (
                      <Box
                        key={product.id}
                        sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}
                      >
                        <Typography
                          variant="body2"
                          sx={{
                            overflow: "hidden",
                            textOverflow: "ellipsis",
                            whiteSpace: "nowrap",
                            maxWidth: 150,
                          }}
                        >
                          {product.name}
                        </Typography>
                        <Chip
                          label={`${product.stock} left`}
                          sx={{
                            backgroundColor: "#fef2f2",
                            color: "#dc2626",
                            fontSize: "0.7rem",
                            height: 20,
                          }}
                        />
                      </Box>
                    ))}
                  </Stack>
                </>
              ) : (
                <Box
                  sx={{
                    height: 300,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "#64748b",
                  }}
                >
                  All products are well-stocked
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Recent Activity */}
      <Card>
        <CardContent>
          <Typography variant="h6" sx={{ mb: 1 }}>
            Recent Activity
          </Typography>
          <Typography variant="body2" sx={{ color: "#64748b" }}>
            Orders and activity will appear here in real-time.
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}
