"use client";

import * as React from "react";
import Chip from "@mui/material/Chip";

interface BadgeProps {
  children?: React.ReactNode;
  variant?: "default" | "secondary" | "destructive" | "outline" | "success" | "warning";
  className?: string;
  [key: string]: unknown;
}

function Badge({ children, variant = "default", className, ...props }: BadgeProps) {
  const colorMap: Record<string, "primary" | "secondary" | "error" | "default" | "success" | "warning"> = {
    default: "primary",
    secondary: "secondary",
    destructive: "error",
    outline: "default",
    success: "success",
    warning: "warning",
  };

  const muiColor = colorMap[variant] || "primary";

  const sxMap: Record<string, object> = {
    default: { backgroundColor: "#3b82f6", color: "#fff" },
    secondary: { backgroundColor: "#f1f5f9", color: "#1e293b" },
    destructive: { backgroundColor: "#fef2f2", color: "#dc2626" },
    outline: { backgroundColor: "transparent", border: "1px solid #e2e8f0", color: "#1e293b" },
    success: { backgroundColor: "#f0fdf4", color: "#16a34a" },
    warning: { backgroundColor: "#fffbeb", color: "#d97706" },
  };

  return (
    <Chip
      label={children}
      color={muiColor}
      size="small"
      sx={{ borderRadius: "6px", fontWeight: 500, fontSize: "0.75rem", ...sxMap[variant] }}
      {...props}
    />
  );
}

export { Badge, Badge as Chip };
