'use client';

import React from 'react';
import { useSearchParams } from 'next/navigation';
import { Container, Typography, Box, Button, Paper } from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import Link from 'next/link';

export default function CheckoutSuccessPage() {
  const searchParams = useSearchParams();
  const paymentIntent = searchParams.get('payment_intent');
  const orderId = searchParams.get('order_id');

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Paper elevation={3} sx={{ p: 4, textAlign: 'center' }}>
        <Box sx={{ mb: 3 }}>
          <CheckCircleIcon sx={{ fontSize: 80, color: 'success.main' }} />
        </Box>
        <Typography variant="h4" gutterBottom color="success.main">
          Payment Successful!
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Thank you for your purchase. Your order has been confirmed.
        </Typography>
        {orderId && (
          <Typography variant="body2" sx={{ mb: 2 }}>
            Order ID: {orderId}
          </Typography>
        )}
        {paymentIntent && (
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Payment ID: {paymentIntent}
          </Typography>
        )}
        <Box sx={{ mt: 4 }}>
          <Button 
            variant="contained" 
            component={Link} 
            href="/dashboard"
            sx={{ mr: 2 }}
          >
            Go to Dashboard
          </Button>
          <Button variant="outlined" component={Link} href="/">
            Continue Shopping
          </Button>
        </Box>
      </Paper>
    </Container>
  );
}
