-- SaleApp Database Initialization Script
-- Version: 1.0.0
-- Description: Initial schema setup for PostgreSQL

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create enum types
DO $$ BEGIN
    CREATE TYPE order_status AS ENUM ('pending', 'completed', 'cancelled', 'refunded');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('admin', 'manager', 'cashier');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- ===================
-- Users Table
-- ===================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role user_role NOT NULL DEFAULT 'cashier',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- ===================
-- Categories Table
-- ===================
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ===================
-- Products Table
-- ===================
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    cost DECIMAL(10, 2),
    stock INTEGER NOT NULL DEFAULT 0,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);

-- ===================
-- Customers Table
-- ===================
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    address TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);

-- ===================
-- Orders Table
-- ===================
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status order_status NOT NULL DEFAULT 'pending',
    subtotal DECIMAL(10, 2) NOT NULL,
    tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_customer ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);

-- ===================
-- Order Items Table
-- ===================
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product ON order_items(product_id);

-- ===================
-- Function: Update updated_at timestamp
-- ===================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ===================
-- Triggers for updated_at
-- ===================
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ===================
-- Seed Data (optional admin user)
-- ===================
-- Default admin password: "admin123" (bcrypt hash)
-- IMPORTANT: Change this password in production!
INSERT INTO users (email, password_hash, first_name, last_name, role)
VALUES (
    'admin@example.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4pWGD5eOYkTql/TO',  -- admin123
    'System',
    'Admin',
    'admin'
) ON CONFLICT (email) DO NOTHING;

-- ===================
-- Sample Categories
-- ===================
INSERT INTO categories (name, description) VALUES
    ('Electronics', 'Electronic devices and accessories'),
    ('Clothing', 'Apparel and fashion items'),
    ('Food & Beverages', 'Food and drink products'),
    ('Home & Garden', 'Home improvement and garden items')
ON CONFLICT DO NOTHING;

-- ===================
-- Grant permissions (adjust as needed)
-- ===================
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO saleapp_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO saleapp_user;

-- ===================
-- Seed Data: Products
-- ===================
INSERT INTO products (sku, name, description, price, cost, stock, category_id) VALUES
    ('ELEC-001', 'Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 29.99, 12.00, 45, (SELECT id FROM categories WHERE name = 'Electronics')),
    ('ELEC-002', 'USB-C Hub 7-Port', 'Multi-port USB-C hub with HDMI output', 49.99, 22.00, 28, (SELECT id FROM categories WHERE name = 'Electronics')),
    ('ELEC-003', 'Bluetooth Headset', 'Over-ear wireless headphones', 89.99, 45.00, 3, (SELECT id FROM categories WHERE name = 'Electronics')),
    ('ELEC-004', 'Mechanical Keyboard', 'RGB backlit mechanical keyboard', 79.99, 38.00, 15, (SELECT id FROM categories WHERE name = 'Electronics')),
    ('ELEC-005', 'Webcam 1080p', 'Full HD webcam with microphone', 59.99, 28.00, 22, (SELECT id FROM categories WHERE name = 'Electronics')),
    ('CLTH-001', 'Cotton T-Shirt', 'Premium cotton casual t-shirt', 19.99, 6.00, 120, (SELECT id FROM categories WHERE name = 'Clothing')),
    ('CLTH-002', 'Denim Jeans', 'Classic fit denim jeans', 49.99, 18.00, 65, (SELECT id FROM categories WHERE name = 'Clothing')),
    ('CLTH-003', 'Hoodie', 'Warm fleece hoodie', 39.99, 14.00, 8, (SELECT id FROM categories WHERE name = 'Clothing')),
    ('FOOD-001', 'Organic Coffee Beans 1kg', 'Premium roasted arabica beans', 24.99, 10.00, 55, (SELECT id FROM categories WHERE name = 'Food & Beverages')),
    ('FOOD-002', 'Green Tea Box', '30 bags of organic green tea', 12.99, 4.00, 90, (SELECT id FROM categories WHERE name = 'Food & Beverages')),
    ('HOME-001', 'LED Desk Lamp', 'Adjustable LED desk lamp with USB charging', 34.99, 15.00, 38, (SELECT id FROM categories WHERE name = 'Home & Garden')),
    ('HOME-002', 'Plant Pot Set', 'Set of 3 ceramic plant pots', 29.99, 11.00, 5, (SELECT id FROM categories WHERE name = 'Home & Garden'))
ON CONFLICT (sku) DO NOTHING;

-- ===================
-- Seed Data: Customers
-- ===================
INSERT INTO customers (email, phone, first_name, last_name, address) VALUES
    ('sara.phet@example.com', '+66-81-234-5678', 'Sara', 'Phet', '123 Sukhumvit Rd, Bangkok'),
    ('peeravit.nap@example.com', '+66-82-345-6789', 'Peeravit', 'Napawan', '456 Rama IV Rd, Bangkok'),
    ('wijitra.tan@example.com', '+66-83-456-7890', 'Wijitra', 'Tanakul', '789 Silom Rd, Bangkok'),
    ('chalida.rai@example.com', '+66-84-567-8901', 'Chalida', 'Raimaitre', '321 Ploenchit Rd, Bangkok'),
    ('krittin.pet@example.com', '+66-85-678-9012', 'Krittin', 'Petcharat', '654 Sathorn Rd, Bangkok')
ON CONFLICT (email) DO NOTHING;

-- ===================
-- Seed Data: Today's Orders (using CTE to capture inserted IDs)
-- ===================
DO $$
DECLARE
    admin_user_id UUID;
BEGIN
    SELECT id INTO admin_user_id FROM users WHERE email = 'admin@example.com' LIMIT 1;

    -- Only seed orders if we have the admin user
    IF admin_user_id IS NOT NULL THEN
        -- Order 1: Completed - Electronics bundle
        WITH ord1 AS (
            INSERT INTO orders (order_number, customer_id, user_id, status, subtotal, tax, discount, total, payment_method, created_at)
            SELECT
                'ORD-2026-0301-001',
                (SELECT id FROM customers WHERE email = 'sara.phet@example.com'),
                admin_user_id,
                'completed',
                119.98,
                8.40,
                0.00,
                128.38,
                'cash',
                CURRENT_TIMESTAMP - INTERVAL '4 hours'
            WHERE (SELECT id FROM customers WHERE email = 'sara.phet@example.com') IS NOT NULL
            RETURNING id
        )
        INSERT INTO order_items (order_id, product_id, quantity, unit_price, discount, total)
        SELECT ord1.id, p.id, 2, 29.99, 0.00, 59.98
        FROM ord1, (SELECT id FROM products WHERE sku = 'ELEC-001') p;

        -- Order 2: Completed - Clothing
        WITH ord2 AS (
            INSERT INTO orders (order_number, customer_id, user_id, status, subtotal, tax, discount, total, payment_method, created_at)
            SELECT
                'ORD-2026-0301-002',
                (SELECT id FROM customers WHERE email = 'peeravit.nap@example.com'),
                admin_user_id,
                'completed',
                59.97,
                4.20,
                5.00,
                59.17,
                'qr_code',
                CURRENT_TIMESTAMP - INTERVAL '3 hours'
            WHERE (SELECT id FROM customers WHERE email = 'peeravit.nap@example.com') IS NOT NULL
            RETURNING id
        )
        INSERT INTO order_items (order_id, product_id, quantity, unit_price, discount, total)
        SELECT ord2.id, p.id, 3, 19.99, 5.00, 54.97
        FROM ord2, (SELECT id FROM products WHERE sku = 'CLTH-001') p;

        -- Order 3: Completed - Mixed items
        WITH ord3 AS (
            INSERT INTO orders (order_number, customer_id, user_id, status, subtotal, tax, discount, total, payment_method, created_at)
            SELECT
                'ORD-2026-0301-003',
                (SELECT id FROM customers WHERE email = 'wijitra.tan@example.com'),
                admin_user_id,
                'completed',
                89.98,
                6.30,
                0.00,
                96.28,
                'cash',
                CURRENT_TIMESTAMP - INTERVAL '1 hour'
            WHERE (SELECT id FROM customers WHERE email = 'wijitra.tan@example.com') IS NOT NULL
            RETURNING id
        )
        INSERT INTO order_items (order_id, product_id, quantity, unit_price, discount, total)
        SELECT ord3.id, p.id, 2, 19.99, 0.00, 39.98
        FROM ord3, (SELECT id FROM products WHERE sku = 'CLTH-001') p;

        -- Order 4: Pending (should not count in completed revenue)
        WITH ord4 AS (
            INSERT INTO orders (order_number, customer_id, user_id, status, subtotal, tax, discount, total, payment_method, created_at)
            SELECT
                'ORD-2026-0301-004',
                (SELECT id FROM customers WHERE email = 'sara.phet@example.com'),
                admin_user_id,
                'pending',
                29.99,
                2.10,
                0.00,
                32.09,
                'cash',
                CURRENT_TIMESTAMP - INTERVAL '30 minutes'
            WHERE (SELECT id FROM customers WHERE email = 'sara.phet@example.com') IS NOT NULL
            RETURNING id
        )
        INSERT INTO order_items (order_id, product_id, quantity, unit_price, discount, total)
        SELECT ord4.id, p.id, 1, 29.99, 0.00, 29.99
        FROM ord4, (SELECT id FROM products WHERE sku = 'ELEC-001') p;
    END IF;
END $$;
