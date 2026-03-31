"use client";

import Box from "@mui/material/Box";
import AppBarMUI from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import InputBase from "@mui/material/InputBase";
import Badge from "@mui/material/Badge";
import Avatar from "@mui/material/Avatar";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Divider from "@mui/material/Divider";
import ListItemIcon from "@mui/material/ListItemIcon";
import { Notifications, Search, Person, Settings, Logout } from "@mui/icons-material";
import { useState, useRef } from "react";
import { useAuth } from "@/hooks/use-auth";
import Link from "next/link";
import { useRouter } from "next/navigation";

interface HeaderProps {
  title?: string;
}

export function Header({ title }: HeaderProps) {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleClose();
    logout();
    router.push("/login");
  };

  return (
    <AppBarMUI
      position="static"
      color="default"
      sx={{
        backgroundColor: "#fff",
        borderBottom: "1px solid #e2e8f0",
        boxShadow: "none",
      }}
    >
      <Toolbar sx={{ gap: 2 }}>
        <Typography variant="h6" sx={{ fontWeight: 600, color: "#1e293b" }}>
          {title}
        </Typography>

        <Box sx={{ flexGrow: 1 }} />

        {/* Search */}
        <Box
          sx={{
            display: { xs: "none", md: "flex" },
            alignItems: "center",
            backgroundColor: "#f1f5f9",
            borderRadius: "8px",
            px: 1.5,
            py: 0.5,
            gap: 1,
            width: 260,
          }}
        >
          <Search sx={{ color: "#64748b", fontSize: 20 }} />
          <InputBase
            placeholder="Search..."
            sx={{
              flex: 1,
              fontSize: "0.875rem",
              "& input": { padding: 0 },
            }}
          />
        </Box>

        {/* Notifications */}
        <IconButton size="small" sx={{ color: "#64748b" }}>
          <Badge badgeContent={3} color="error">
            <Notifications sx={{ fontSize: 20 }} />
          </Badge>
        </IconButton>

        {/* User Menu */}
        <IconButton onClick={handleClick} size="small">
          <Avatar
            sx={{
              width: 32,
              height: 32,
              backgroundColor: "#3b82f6",
              fontSize: "0.875rem",
            }}
          >
            {user?.firstName?.[0]}{user?.lastName?.[0]}
          </Avatar>
        </IconButton>
        <Menu
          anchorEl={anchorEl}
          open={open}
          onClose={handleClose}
          transformOrigin={{ horizontal: "right", vertical: "top" }}
          anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
          PaperProps={{
            sx: {
              mt: 1,
              minWidth: 200,
              boxShadow: "0 4px 16px rgba(0,0,0,0.12)",
            },
          }}
        >
          <Box sx={{ px: 2, py: 1.5 }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
              {user?.firstName} {user?.lastName}
            </Typography>
            <Typography variant="caption" sx={{ color: "#64748b" }}>
              {user?.email}
            </Typography>
          </Box>
          <Divider />
          <MenuItem
            component={Link}
            href="/settings"
            onClick={handleClose}
            sx={{ gap: 1.5 }}
          >
            <ListItemIcon>
              <Settings fontSize="small" />
            </ListItemIcon>
            Settings
          </MenuItem>
          <Divider />
          <MenuItem onClick={handleLogout} sx={{ gap: 1.5, color: "#dc2626" }}>
            <ListItemIcon>
              <Logout fontSize="small" sx={{ color: "#dc2626" }} />
            </ListItemIcon>
            Sign out
          </MenuItem>
        </Menu>
      </Toolbar>
    </AppBarMUI>
  );
}
