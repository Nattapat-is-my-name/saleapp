"use client";

import * as React from "react";
import SelectMUI from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import InputLabel from "@mui/material/InputLabel";
import ListItemIcon from "@mui/material/ListItemIcon";
import CheckIcon from "@mui/icons-material/Check";

interface SelectProps {
  value?: string;
  onChange?: (value: string) => void;
  children?: React.ReactNode;
  placeholder?: string;
  className?: string;
  size?: "small" | "default";
  [key: string]: unknown;
}

function Select({ value, onChange, children, placeholder, size = "small", className, ...props }: SelectProps) {
  const muiSize = size === "default" ? "medium" : size;
  return (
    <FormControl fullWidth size={muiSize}>
      {placeholder && <InputLabel>{placeholder}</InputLabel>}
      <SelectMUI
        value={value || ""}
        onChange={(e) => onChange?.(e.target.value)}
        label={placeholder}
        sx={{ borderRadius: "8px" }}
        renderValue={(selected) => {
          if (!selected) return <span style={{ color: "#64748b" }}>{placeholder}</span>;
          return selected;
        }}
        {...props}
      >
        {children}
      </SelectMUI>
    </FormControl>
  );
}

function SelectItem({
  children,
  value,
  ...props
}: React.ComponentProps<typeof MenuItem> & { children?: React.ReactNode }) {
  return <MenuItem value={value} {...props}>{children}</MenuItem>;
}

function SelectValue({ children }: { children?: React.ReactNode }) {
  return <>{children}</>;
}

export { Select, SelectItem, SelectValue };
