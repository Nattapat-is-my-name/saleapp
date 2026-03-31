"use client";

import * as React from "react";
import TextField from "@mui/material/TextField";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  className?: string;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      <TextField
        inputRef={ref}
        type={type}
        size="small"
        fullWidth
        slotProps={{
          input: {
            sx: {
              borderRadius: "8px",
              border: "1px solid #e2e8f0",
              "&:hover": { borderColor: "#94a3b8" },
              "&.Mui-focused": { borderColor: "#3b82f6" },
            },
          },
        }}
        {...(props as React.ComponentProps<typeof TextField>)}
      />
    );
  }
);
Input.displayName = "Input";

export { Input };
