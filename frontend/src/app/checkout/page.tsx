'use client';

import React, { useEffect, useState, Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import { Elements } from '@stripe/react-stripe-js';
import { Container, Typography, Paper, CircularProgress, Alert, Box } from '@mui/material';
import { stripePromise } from '@/lib/stripe';
import { usePayment } from '@/hooks/usePayment';
import CheckoutForm from '@/components/payment/CheckoutForm';
import { PaymentIntent } from '@/types/payment';

function CheckoutContent() {
  const searchParams = useSearchParams();
  const orderId = searchParams.get('order_id');
  const { createPaymentIntent, loading, error } = usePayment();
  const [clientSecret, setClientSecret] = useState<string | null>(null);
  const [paymentData, setPaymentData] = useState<PaymentIntent | null>(null);

  useEffect(() => {
    if (orderId) {
      createPaymentIntent({ order_id: orderId, currency: 'usd' })
        .then((data) => {
          setClientSecret(data.client_secret);
          setPaymentData(data);
        })
        .catch(console.error);
    }
  }, [orderId]);

  const handleSuccess = (paymentId: string) => {
    console.log('Payment succeeded:', paymentId);
  };

  const handleError = (error: string) => {
    console.error('Payment error:', error);
  };

  if (!orderId) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Alert severity="error">No order ID provided</Alert>
      </Container>
    );
  }

  if (loading || !clientSecret) {
    return (
      <Container maxWidth="sm" sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress />
        <Typography sx={{ mt: 2 }}>Preparing checkout...</Typography>
      </Container>
    );
  }

  const appearance = {
    theme: 'stripe' as const,
    variables: {
      colorPrimary: '#1976d2',
    },
  };

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Typography variant="h4" gutterBottom>
          Checkout
        </Typography>
        <Typography variant="body1" color="text.secondary" gutterBottom>
          Order ID: {orderId}
        </Typography>
        {paymentData && (
          <Typography variant="h5" sx={{ my: 2 }}>
            Total: ${(paymentData.amount / 100).toFixed(2)}
          </Typography>
        )}
        <Elements stripe={stripePromise} options={{ clientSecret, appearance }}>
          <CheckoutForm 
            orderId={orderId} 
            onSuccess={handleSuccess} 
            onError={handleError} 
          />
        </Elements>
      </Paper>
    </Container>
  );
}

function CheckoutLoading() {
  return (
    <Container maxWidth="sm" sx={{ py: 8, textAlign: 'center' }}>
      <CircularProgress />
      <Typography sx={{ mt: 2 }}>Loading checkout...</Typography>
    </Container>
  );
}

export default function CheckoutPage() {
  return (
    <Suspense fallback={<CheckoutLoading />}>
      <CheckoutContent />
    </Suspense>
  );
}
