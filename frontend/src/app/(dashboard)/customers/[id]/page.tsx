"use client";

import { useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Grid from "@mui/material/Grid";
import TextField from "@mui/material/TextField";
import Stack from "@mui/material/Stack";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import { ArrowBack } from "@mui/icons-material";
import { useCustomer, useUpdateCustomer } from "@/hooks/use-customers";
import type { CreateCustomerInput } from "@/types/customer";

export default function EditCustomerPage() {
  const router = useRouter();
  const params = useParams();
  const customerId = params.id as string;

  const { data: customer, isLoading, error } = useCustomer(customerId);
  const updateCustomer = useUpdateCustomer();

  const [formData, setFormData] = useState<CreateCustomerInput>({
    firstName: "",
    lastName: "",
    email: "",
    phone: "",
    address: "",
    notes: "",
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [initialized, setInitialized] = useState(false);

  // Initialize form when customer loads
  if (customer && !initialized) {
    setFormData({
      firstName: customer.firstName || "",
      lastName: customer.lastName || "",
      email: customer.email || "",
      phone: customer.phone || "",
      address: customer.address || "",
      notes: customer.notes || "",
    });
    setInitialized(true);
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.firstName?.trim()) {
      newErrors.firstName = "First name is required";
    }
    if (!formData.lastName?.trim()) {
      newErrors.lastName = "Last name is required";
    }
    if (!formData.email?.trim()) {
      newErrors.email = "Email is required";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = "Invalid email format";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    try {
      await updateCustomer.mutateAsync({ id: customerId, ...formData });
      router.push("/customers");
    } catch (err) {
      console.error("Failed to update customer:", err);
    }
  };

  const handleChange = (field: keyof CreateCustomerInput, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  if (isLoading) {
    return (
      <Box sx={{ maxWidth: 720 }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
          <Button
            variant="text"
            startIcon={<ArrowBack />}
            component={Link}
            href="/customers"
            sx={{ color: "#64748b" }}
          >
            Back
          </Button>
        </Box>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
          Edit Customer
        </Typography>
        <Typography sx={{ color: "#64748b", mb: 3 }}>
          Loading customer details...
        </Typography>
        <Card>
          <CardContent>
            <Box sx={{ py: 4, textAlign: "center" }}>
              <CircularProgress />
            </Box>
          </CardContent>
        </Card>
      </Box>
    );
  }

  if (error || !customer) {
    return (
      <Box sx={{ maxWidth: 720 }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
          <Button
            variant="text"
            startIcon={<ArrowBack />}
            component={Link}
            href="/customers"
            sx={{ color: "#64748b" }}
          >
            Back
          </Button>
        </Box>
        <Alert severity="error">Customer not found.</Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 720 }}>
      <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
        <Button
          variant="text"
          startIcon={<ArrowBack />}
          component={Link}
          href="/customers"
          sx={{ color: "#64748b" }}
        >
          Back
        </Button>
      </Box>

      <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
        Edit Customer
      </Typography>
      <Typography sx={{ color: "#64748b", mb: 3 }}>
        Update customer information
      </Typography>

      <Card>
        <CardContent>
          <Typography variant="h6" sx={{ mb: 1 }}>
            Customer Information
          </Typography>
          <Typography variant="body2" sx={{ color: "#64748b", mb: 3 }}>
            Update the details below
          </Typography>

          <Box component="form" onSubmit={handleSubmit}>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  label="First Name"
                  value={formData.firstName}
                  onChange={(e) => handleChange("firstName", e.target.value)}
                  error={!!errors.firstName}
                  helperText={errors.firstName}
                  fullWidth
                  size="small"
                  required
                  placeholder="John"
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  label="Last Name"
                  value={formData.lastName}
                  onChange={(e) => handleChange("lastName", e.target.value)}
                  error={!!errors.lastName}
                  helperText={errors.lastName}
                  fullWidth
                  size="small"
                  required
                  placeholder="Doe"
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  label="Email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => handleChange("email", e.target.value)}
                  error={!!errors.email}
                  helperText={errors.email}
                  fullWidth
                  size="small"
                  required
                  placeholder="john.doe@example.com"
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  label="Phone"
                  value={formData.phone}
                  onChange={(e) => handleChange("phone", e.target.value)}
                  fullWidth
                  size="small"
                  placeholder="+1 234 567 8900"
                />
              </Grid>

              <Grid size={{ xs: 12 }}>
                <TextField
                  label="Address"
                  value={formData.address}
                  onChange={(e) => handleChange("address", e.target.value)}
                  fullWidth
                  size="small"
                  placeholder="123 Main St, City, Country"
                />
              </Grid>

              <Grid size={{ xs: 12 }}>
                <TextField
                  label="Notes"
                  value={formData.notes}
                  onChange={(e) => handleChange("notes", e.target.value)}
                  fullWidth
                  size="small"
                  multiline
                  rows={3}
                  placeholder="Additional notes about this customer..."
                />
              </Grid>
            </Grid>

            <Stack direction="row" spacing={2} sx={{ mt: 4 }}>
              <Button
                type="submit"
                variant="contained"
                disabled={updateCustomer.isPending}
                fullWidth
                startIcon={updateCustomer.isPending ? <CircularProgress size={16} /> : null}
              >
                {updateCustomer.isPending ? "Saving..." : "Save Changes"}
              </Button>
              <Button
                type="button"
                variant="outlined"
                onClick={() => router.push("/customers")}
                fullWidth
              >
                Cancel
              </Button>
            </Stack>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
