export interface PaymentIntent {
  id: string;
  client_secret: string;
  payment_id: string;
  amount: number;
  currency: string;
}

export interface PaymentStatus {
  id: string;
  order_id: string;
  stripe_payment_id: string;
  amount: number;
  currency: string;
  status: 'pending' | 'succeeded' | 'failed' | 'refunded';
  payment_method?: string;
  error_message?: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePaymentIntentRequest {
  order_id: string;
  currency?: string;
}
