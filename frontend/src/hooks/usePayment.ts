import { useState } from 'react';
import { PaymentIntent, PaymentStatus, CreatePaymentIntentRequest } from '@/types/payment';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export function usePayment() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createPaymentIntent = async (data: CreatePaymentIntentRequest): Promise<PaymentIntent> => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/payments/intent`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        },
        body: JSON.stringify(data),
      });
      
      if (!response.ok) {
        throw new Error('Failed to create payment intent');
      }
      
      return await response.json();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Payment failed';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const getPaymentStatus = async (orderId: string): Promise<PaymentStatus> => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/payments/${orderId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        },
      });
      
      if (!response.ok) {
        throw new Error('Failed to get payment status');
      }
      
      return await response.json();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to get status';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { createPaymentIntent, getPaymentStatus, loading, error };
}
