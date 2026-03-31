"use client";

import * as React from "react";
import FormLabelMUI from "@mui/material/FormLabel";

interface LabelProps {
  children?: React.ReactNode;
  className?: string;
  htmlFor?: string;
}

export function Label({ children, className, htmlFor }: LabelProps) {
  return (
    <FormLabelMUI
      htmlFor={htmlFor}
      className={className}
      sx={{
        display: "block",
        fontSize: "0.875rem",
        fontWeight: 500,
        color: "text.secondary",
        mb: 0.5,
      }}
    >
      {children}
    </FormLabelMUI>
  );
}
