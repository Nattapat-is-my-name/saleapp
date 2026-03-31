import * as React from "react";
import CardMUI from "@mui/material/Card";
import CardContentMUI from "@mui/material/CardContent";
import CardActionsMUI from "@mui/material/CardActions";

interface CardProps {
  children: React.ReactNode;
  className?: string;
  size?: "default" | "sm";
  [key: string]: unknown;
}

function Card({ children, className, size = "default", ...props }: CardProps) {
  return (
    <CardMUI
      sx={{
        display: "flex",
        flexDirection: "column",
        gap: size === "sm" ? 1.5 : 2,
        py: size === "sm" ? 1.5 : 2,
        ...(className ? {} : {}),
      }}
      {...props}
    >
      {children}
    </CardMUI>
  );
}

function CardHeader({
  children,
  className,
  ...props
}: React.ComponentProps<typeof CardMUI> & { className?: string }) {
  return (
    <CardContentMUI
      sx={{
        px: 2,
        pt: 2,
        pb: 1,
        display: "grid",
        gap: 0.5,
        alignItems: "start",
        "&:last-child": { pb: 1 },
      }}
      {...props}
    >
      {children}
    </CardContentMUI>
  );
}

function CardTitle({
  children,
  className,
  ...props
}: React.ComponentProps<"div"> & { className?: string }) {
  return (
    <div
      style={{
        fontSize: "1rem",
        fontWeight: 600,
        lineHeight: 1.4,
      }}
      {...props}
    >
      {children}
    </div>
  );
}

function CardDescription({
  children,
  className,
  ...props
}: React.ComponentProps<"div"> & { className?: string }) {
  return (
    <div
      style={{
        fontSize: "0.875rem",
        color: "#64748b",
      }}
      {...props}
    >
      {children}
    </div>
  );
}

function CardContent({
  children,
  className,
  ...props
}: React.ComponentProps<typeof CardContentMUI> & { className?: string }) {
  return (
    <CardContentMUI sx={{ px: 2, pt: 0, pb: 2 }} {...props}>
      {children}
    </CardContentMUI>
  );
}

function CardFooter({
  children,
  className,
  ...props
}: React.ComponentProps<typeof CardActionsMUI> & { className?: string }) {
  return (
    <CardActionsMUI
      sx={{
        px: 2,
        py: 1.5,
        borderTop: "1px solid #e2e8f0",
        backgroundColor: "#f8fafc",
      }}
      {...props}
    >
      {children}
    </CardActionsMUI>
  );
}

export { Card, CardHeader, CardFooter, CardTitle, CardDescription, CardContent };
