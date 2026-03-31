"use client";

import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { useState } from "react";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Alert from "@mui/material/Alert";
import CircularProgress from "@mui/material/CircularProgress";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import { Visibility, VisibilityOff, ShoppingCart } from "@mui/icons-material";
import { api } from "@/lib/api";

const loginSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
  password: z.string().min(1, "Password is required"),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginForm) => {
    setError(null);
    setIsLoading(true);

    try {
      const response = await api.post<{
        data: {
          token: string;
          expiresAt: string;
          user: {
            id: string;
            email: string;
            role: string;
          };
        };
      }>("/auth/login", data);

      localStorage.setItem("auth_token", response.data.token);
      localStorage.setItem(
        "auth_user",
        JSON.stringify({
          id: response.data.user.id,
          email: response.data.user.email,
          role: response.data.user.role,
        })
      );

      router.push("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "linear-gradient(135deg, #f0f9ff 0%, #e0e7ff 100%)",
        p: 2,
      }}
    >
      <Card sx={{ width: "100%", maxWidth: 420 }}>
        <CardContent sx={{ p: 4 }}>
          <Box sx={{ textAlign: "center", mb: 4 }}>
            <Box
              sx={{
                width: 56,
                height: 56,
                borderRadius: "50%",
                backgroundColor: "#eff6ff",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                mx: "auto",
                mb: 2,
              }}
            >
              <ShoppingCart sx={{ color: "#3b82f6", fontSize: 28 }} />
            </Box>
            <Typography variant="h5" sx={{ fontWeight: 700, mb: 0.5 }}>
              Welcome back
            </Typography>
            <Typography variant="body2" sx={{ color: "#64748b" }}>
              Sign in to your SaleApp account
            </Typography>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ display: "flex", flexDirection: "column", gap: 2.5 }}>
            <TextField
              label="Email"
              type="email"
              placeholder="admin@example.com"
              {...register("email")}
              error={!!errors.email}
              helperText={errors.email?.message}
              fullWidth
              size="small"
            />

            <TextField
              label="Password"
              type={showPassword ? "text" : "password"}
              placeholder="••••••••"
              {...register("password")}
              error={!!errors.password}
              helperText={errors.password?.message}
              fullWidth
              size="small"
              slotProps={{
                input: {
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() => setShowPassword(!showPassword)}
                        edge="end"
                        size="small"
                      >
                        {showPassword ? (
                          <VisibilityOff sx={{ fontSize: 18 }} />
                        ) : (
                          <Visibility sx={{ fontSize: 18 }} />
                        )}
                      </IconButton>
                    </InputAdornment>
                  ),
                },
              }}
            />

            <Button
              type="submit"
              variant="contained"
              disabled={isLoading}
              fullWidth
              sx={{ mt: 1 }}
            >
              {isLoading ? (
                <>
                  <CircularProgress size={16} color="inherit" sx={{ mr: 1 }} />
                  Signing in...
                </>
              ) : (
                "Sign in"
              )}
            </Button>
          </Box>

          <Box sx={{ mt: 3, textAlign: "center" }}>
            <Typography variant="caption" sx={{ color: "#64748b" }}>
              Demo credentials: admin@saleapp.com / password
            </Typography>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
