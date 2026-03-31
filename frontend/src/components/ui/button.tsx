"use client";

import * as React from "react";
import ButtonMUI from "@mui/material/Button";

type ButtonVariant = "contained" | "outlined" | "text";
type ButtonSize = "small" | "medium" | "large";

interface ButtonProps {
  variant?: ButtonVariant;
  size?: ButtonSize;
  className?: string;
  children?: React.ReactNode;
  disabled?: boolean;
  onClick?: () => void;
  type?: "button" | "submit" | "reset";
  startIcon?: React.ReactNode;
}

export function Button({
  variant = "contained",
  size = "medium",
  className,
  children,
  disabled,
  onClick,
  type = "button",
  startIcon,
}: ButtonProps) {
  return (
    <ButtonMUI
      variant={variant}
      size={size}
      className={className}
      disabled={disabled}
      onClick={onClick}
      type={type}
      startIcon={startIcon}
    >
      {children}
    </ButtonMUI>
  );
}

export type { ButtonProps };
