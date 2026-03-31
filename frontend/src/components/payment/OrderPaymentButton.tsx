'use client';

import React from 'react';
import { Button } from '@mui/material';
import { useRouter } from 'next/navigation';

interface OrderPaymentButtonProps {
  orderId: string;
  disabled?: boolean;
}

export default function OrderPaymentButton({ orderId, disabled }: OrderPaymentButtonProps) {
  const router = useRouter();

  const handleClick = () => {
    router.push(`/checkout?order_id=${orderId}`);
  };

  return (
    <Button
      variant="contained"
      color="primary"
      onClick={handleClick}
      disabled={disabled}
    >
      Pay Now
    </Button>
  );
}
