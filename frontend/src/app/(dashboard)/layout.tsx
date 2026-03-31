"use client";

import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import { Sidebar } from "@/components/layout/sidebar";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Box sx={{ display: "flex", minHeight: "100vh" }}>
      <Sidebar />
      <Box
        sx={{
          flexGrow: 1,
          backgroundColor: "#f8fafc",
        }}
      >
        <Container maxWidth={false} sx={{ py: 3 }}>
          {children}
        </Container>
      </Box>
    </Box>
  );
}
