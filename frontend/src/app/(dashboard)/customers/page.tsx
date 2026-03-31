"use client";

import Link from "next/link";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import Card from "@mui/material/Card";
import { Add, Search } from "@mui/icons-material";

export default function CustomersPage() {
  return (
    <Box>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 3 }}>
        <Box>
          <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
            Customers
          </Typography>
          <Typography sx={{ color: "#64748b" }}>
            Manage your customer database
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<Add />}
          component={Link}
          href="/customers/new"
        >
          Add Customer
        </Button>
      </Box>

      <Card>
        <Box sx={{ p: 2 }}>
          <TextField
            placeholder="Search customers..."
            size="small"
            slotProps={{
              input: {
                startAdornment: (
                  <InputAdornment position="start">
                    <Search sx={{ color: "#64748b" }} />
                  </InputAdornment>
                ),
              },
            }}
            sx={{ width: 300 }}
          />
        </Box>
        <Box sx={{ py: 8, textAlign: "center", color: "#64748b" }}>
          No customers found. Add your first customer to get started.
        </Box>
      </Card>
    </Box>
  );
}
