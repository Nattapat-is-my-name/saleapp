#!/bin/bash
# Stripe webhook forwarding script for local development
# Usage: ./scripts/stripe-webhook.sh

# The port where Stripe CLI will forward webhooks
LOCAL_PORT=8080

# Check if stripe CLI is installed
if ! command -v stripe &> /dev/null; then
    echo "Stripe CLI not found. Install from: https://stripe.com/docs/stripe-cli"
    exit 1
fi

echo "Listening for Stripe webhooks..."
stripe listen --forward-to localhost:${LOCAL_PORT}/api/v1/payments/webhook
