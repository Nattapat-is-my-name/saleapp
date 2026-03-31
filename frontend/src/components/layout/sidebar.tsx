"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Divider from "@mui/material/Divider";
import Typography from "@mui/material/Typography";
import {
  Dashboard,
  Inventory2,
  ShoppingBag,
  People,
  BarChart,
  Settings,
  Logout,
} from "@mui/icons-material";
import Button from "@mui/material/Button";

const DRAWER_WIDTH = 260;

const navigation = [
  { name: "Dashboard", href: "/", icon: Dashboard },
  { name: "Products", href: "/products", icon: Inventory2 },
  { name: "Orders", href: "/orders", icon: ShoppingBag },
  { name: "Customers", href: "/customers", icon: People },
  { name: "Reports", href: "/reports", icon: BarChart },
  { name: "Settings", href: "/settings", icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();

  const handleLogout = () => {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_user");
    router.push("/login");
  };

  const drawerContent = (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
      {/* Logo */}
      <Box
        sx={{
          height: 64,
          display: "flex",
          alignItems: "center",
          px: 2.5,
          borderBottom: "1px solid #e2e8f0",
        }}
      >
        <Link href="/" style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <Box
            sx={{
              width: 36,
              height: 36,
              borderRadius: "8px",
              backgroundColor: "#3b82f6",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <ShoppingBag sx={{ color: "#fff", fontSize: 20 }} />
          </Box>
          <Typography variant="h6" sx={{ fontWeight: 700, color: "#1e293b" }}>
            SaleApp
          </Typography>
        </Link>
      </Box>

      {/* Navigation */}
      <List sx={{ flex: 1, px: 1, py: 2 }}>
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          return (
            <ListItem key={item.name} disablePadding sx={{ mb: 0.5 }}>
              <ListItemButton
                component={Link}
                href={item.href}
                sx={{
                  borderRadius: "8px",
                  py: 1,
                  px: 1.5,
                  backgroundColor: isActive ? "#3b82f6" : "transparent",
                  color: isActive ? "#fff" : "#64748b",
                  "&:hover": {
                    backgroundColor: isActive ? "#2563eb" : "#f1f5f9",
                    color: isActive ? "#fff" : "#1e293b",
                  },
                }}
              >
                <ListItemIcon
                  sx={{
                    minWidth: 36,
                    color: "inherit",
                  }}
                >
                  <item.icon sx={{ fontSize: 20 }} />
                </ListItemIcon>
                <ListItemText
                  primary={item.name}
                  primaryTypographyProps={{ fontSize: "0.875rem", fontWeight: 500 }}
                />
              </ListItemButton>
            </ListItem>
          );
        })}
      </List>

      {/* Logout */}
      <Divider />
      <Box sx={{ p: 1.5 }}>
        <Button
          variant="text"
          onClick={handleLogout}
          startIcon={<Logout sx={{ fontSize: 18 }} />}
          sx={{
            justifyContent: "flex-start",
            color: "#64748b",
            px: 1.5,
            "&:hover": { color: "#dc2626", backgroundColor: "#fef2f2" },
          }}
        >
          Sign out
        </Button>
      </Box>
    </Box>
  );

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: DRAWER_WIDTH,
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          width: DRAWER_WIDTH,
          boxSizing: "border-box",
          backgroundColor: "#fff",
        },
      }}
    >
      {drawerContent}
    </Drawer>
  );
}
