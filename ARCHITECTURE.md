# SaleApp Architecture

**Version:** 1.0.0
**Last Updated:** 2026-03-31
**Stack:** Go (backend) + Next.js 14+ (frontend)
**Target:** Sales/POS Application

---

## 1. System Overview

### Purpose
SaleApp is a point-of-sale (POS) and sales management application designed for retail businesses. It handles:
- Product catalog management
- Sales transactions (cash register operations)
- Inventory tracking
- Customer management
- Sales reporting and analytics

### High-Level Architecture
```
┌─────────────────┐     HTTP/JSON      ┌─────────────────┐
│   Next.js 14    │ ◄─────────────────► │   Go Backend    │
│  (Frontend)     │    REST API         │   (Gin)         │
│                 │                     │                 │
│  - App Router   │                     │  - Handlers     │
│  - React Query  │                     │  - Services     │
│  - TailwindCSS  │                     │  - Models       │
└─────────────────┘                     └────────┬────────┘
                                                  │
                                                  │ SQL
                                                  ▼
                                          ┌───────────────┐
                                          │   PostgreSQL  │
                                          └───────────────┘
```

---

## 2. Technology Stack

### Backend
| Component | Technology | Version | Rationale |
|-----------|------------|---------|-----------|
| Language | Go | 1.21+ | Type safety, concurrency, fast compile |
| Framework | Gin | v1.9+ | Performance, middleware ecosystem, Chi as alternative |
| ORM | GORM | v1.25+ | Developer experience, auto-migrations |
| Database | PostgreSQL | 15+ | ACID compliance, JSON support, relational strength |
| Auth | JWT (golang-jwt/jwt) | v5 | Stateless, industry standard |
| Config | viper | v1.18+ | YAML/TOML support, environment overrides |
| Logger | zerolog | v1.31+ | JSON output, performance |
| Migration | golang-migrate | v4 | Versioned SQL migrations |
| Testing | testify | v1.8+ | Assertions, mocking |

### Frontend
| Component | Technology | Version | Rationale |
|-----------|------------|---------|-----------|
| Framework | Next.js | 14+ | App Router, Server Components, RSC |
| Language | TypeScript | 5.x | Type safety across stack |
| Styling | TailwindCSS | 3.4+ | Utility-first, consistent design system |
| State/Fetching | TanStack Query (React Query) | v5 | Caching, optimistic updates, background refetch |
| Forms | React Hook Form + Zod | latest | Performance, schema validation |
| UI Components | shadcn/ui | latest | Radix primitives, customizable |
| Icons | Lucide React | latest | Consistent, tree-shakeable |
| Charts | Recharts | latest | React-native, composable |

---

## 3. Project Folder Structure

### Root
```
saleapp/
├── backend/                  # Go application
├── frontend/                 # Next.js application
├── docs/                     # Shared documentation
├── ARCHITECTURE.md           # This file
└── README.md
```

### Backend Structure
```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration loading (Viper)
│   ├── models/
│   │   ├── product.go
│   │   ├── customer.go
│   │   ├── order.go
│   │   ├── order_item.go
│   │   └── user.go
│   ├── repository/
│   │   ├── product_repo.go
│   │   ├── customer_repo.go
│   │   ├── order_repo.go
│   │   └── user_repo.go
│   ├── service/
│   │   ├── product_svc.go
│   │   ├── customer_svc.go
│   │   ├── order_svc.go
│   │   ├── auth_svc.go
│   │   └── reporting_svc.go
│   ├── handler/
│   │   ├── product_handler.go
│   │   ├── customer_handler.go
│   │   ├── order_handler.go
│   │   ├── auth_handler.go
│   │   └── reporting_handler.go
│   ├── middleware/
│   │   ├── auth.go           # JWT validation
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── ratelimit.go
│   │   └── recovery.go
│   └── dto/
│       ├── request/
│       │   ├── create_product.go
│       │   ├── create_order.go
│       │   └── login.go
│       └── response/
│           ├── product.go
│           ├── order.go
│           ├── token.go
│           └── error.go
├── migrations/
│   ├── 001_init_schema.up.sql
│   └── 001_init_schema.down.sql
├── pkg/
│   ├── response/
│   │   └── response.go       # Standardized API responses
│   ├── validator/
│   │   └── validator.go      # Custom validators
│   └── errors/
│       └── errors.go         # Custom error types
├── .env.example
├── go.mod
├── go.sum
└── Makefile
```

### Frontend Structure
```
frontend/
├── src/
│   ├── app/
│   │   ├── (auth)/
│   │   │   ├── login/
│   │   │   │   └── page.tsx
│   │   │   └── layout.tsx
│   │   ├── (dashboard)/
│   │   │   ├── products/
│   │   │   │   ├── page.tsx
│   │   │   │   ├── [id]/
│   │   │   │   │   └── page.tsx
│   │   │   │   └── new/
│   │   │   │       └── page.tsx
│   │   │   ├── customers/
│   │   │   ├── orders/
│   │   │   ├── reports/
│   │   │   ├── settings/
│   │   │   ├── layout.tsx     # Dashboard layout with sidebar
│   │   │   └── page.tsx       # Dashboard home
│   │   ├── api/               # API route handlers (if needed)
│   │   ├── layout.tsx         # Root layout
│   │   ├── page.tsx           # Root redirect to dashboard
│   │   └── globals.css
│   ├── components/
│   │   ├── ui/                # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── input.tsx
│   │   │   ├── table.tsx
│   │   │   └── ...
│   │   ├── layout/
│   │   │   ├── sidebar.tsx
│   │   │   ├── header.tsx
│   │   │   └── shell.tsx
│   │   ├── products/
│   │   │   ├── product-table.tsx
│   │   │   ├── product-form.tsx
│   │   │   └── product-card.tsx
│   │   ├── orders/
│   │   │   ├── order-table.tsx
│   │   │   ├── order-form.tsx
│   │   │   └── receipt.tsx
│   │   ├── customers/
│   │   │   └── customer-form.tsx
│   │   └── shared/
│   │       ├── data-table.tsx  # Generic table wrapper
│   │       ├── modal.tsx
│   │       ├── toast.tsx
│   │       └── loading-spinner.tsx
│   ├── hooks/
│   │   ├── use-products.ts
│   │   ├── use-orders.ts
│   │   ├── use-customers.ts
│   │   └── use-auth.ts
│   ├── lib/
│   │   ├── api.ts             # Axios/fetch client
│   │   ├── utils.ts           # cn(), formatCurrency(), etc.
│   │   └── validators/        # Zod schemas
│   │       ├── product.ts
│   │       ├── order.ts
│   │       └── customer.ts
│   ├── types/
│   │   ├── product.ts
│   │   ├── order.ts
│   │   ├── customer.ts
│   │   └── api.ts             # Shared API types
│   └── store/
│       └── auth-store.ts      # Zustand auth state (if needed)
├── public/
├── .env.local
├── next.config.js
├── tailwind.config.ts
├── tsconfig.json
├── package.json
└── .gitignore
```

**Note on App Router:** Route groups `(auth)` and `(dashboard)` are for layout organization only. They don't affect URL paths.

---

## 4. Data Models

### Product
```go
type Product struct {
    ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    SKU         string     `gorm:"uniqueIndex;size:50;not null"`
    Name        string     `gorm:"size:255;not null"`
    Description string     `gorm:"type:text"`
    Price       decimal.Decimal `gorm:"type:decimal(10,2);not null"`
    Cost        decimal.Decimal `gorm:"type:decimal(10,2)"` // For margin calculation
    Stock       int        `gorm:"default:0"`
    CategoryID  *uuid.UUID `gorm:"type:uuid"`
    Category    *Category  `gorm:"foreignKey:CategoryID"`
    IsActive    bool       `gorm:"default:true"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Customer
```go
type Customer struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email     string    `gorm:"uniqueIndex;size:255"`
    Phone     string    `gorm:"size:20"`
    FirstName string    `gorm:"size:100"`
    LastName  string    `gorm:"size:100"`
    Address   string    `gorm:"type:text"`
    Notes     string    `gorm:"type:text"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Order
```go
type Order struct {
    ID           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    OrderNumber  string      `gorm:"uniqueIndex;size:50;not null"`
    CustomerID   *uuid.UUID  `gorm:"type:uuid"`
    Customer     *Customer   `gorm:"foreignKey:CustomerID"`
    UserID       uuid.UUID   `gorm:"type:uuid;not null"` // Cashier
    User         User        `gorm:"foreignKey:UserID"`
    Status       OrderStatus `gorm:"type:varchar(20);default:'pending'"`
    Subtotal     decimal.Decimal `gorm:"type:decimal(10,2);not null"`
    Tax          decimal.Decimal `gorm:"type:decimal(10,2);default:0"`
    Discount     decimal.Decimal `gorm:"type:decimal(10,2);default:0"`
    Total        decimal.Decimal `gorm:"type:decimal(10,2);not null"`
    PaymentMethod string     `gorm:"size:50"`
    Notes        string      `gorm:"type:text"`
    Items        []OrderItem `gorm:"foreignKey:OrderID"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type OrderItem struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    OrderID   uuid.UUID      `gorm:"type:uuid;not null"`
    ProductID uuid.UUID      `gorm:"type:uuid;not null"`
    Product   Product        `gorm:"foreignKey:ProductID"`
    Quantity  int            `gorm:"not null"`
    UnitPrice decimal.Decimal `gorm:"type:decimal(10,2);not null"`
    Discount  decimal.Decimal `gorm:"type:decimal(10,2);default:0"`
    Total     decimal.Decimal `gorm:"type:decimal(10,2);not null"`
    CreatedAt time.Time
}

type OrderStatus string
const (
    StatusPending   OrderStatus = "pending"
    StatusCompleted OrderStatus = "completed"
    StatusCancelled OrderStatus = "cancelled"
    StatusRefunded  OrderStatus = "refunded"
)
```

### User (for authentication)
```go
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email        string    `gorm:"uniqueIndex;size:255;not null"`
    PasswordHash string    `gorm:"size:255;not null"`
    FirstName    string    `gorm:"size:100"`
    LastName     string    `gorm:"size:100"`
    Role         UserRole  `gorm:"type:varchar(20);default:'cashier'"`
    IsActive     bool      `gorm:"default:true"`
    LastLoginAt  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type UserRole string
const (
    RoleAdmin   UserRole = "admin"
    RoleManager UserRole = "manager"
    RoleCashier UserRole = "cashier"
)
```

---

## 5. API Design

### Base URL
```
/api/v1
```

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login with email/password |
| POST | `/auth/logout` | Invalidate session |
| GET | `/auth/me` | Get current user |

**Login Request:**
```json
POST /api/v1/auth/login
{
  "email": "admin@example.com",
  "password": "secret"
}
```

**Login Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expiresAt": "2026-04-01T06:27:00Z",
  "user": {
    "id": "uuid",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/products` | List products (paginated) |
| GET | `/products/:id` | Get single product |
| POST | `/products` | Create product |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Soft delete (set inactive) |

**List Response:**
```json
GET /api/v1/products?page=1&limit=20&search=shirt
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8
  }
}
```

### Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/orders` | List orders (paginated, filterable) |
| GET | `/orders/:id` | Get single order with items |
| POST | `/orders` | Create new order (sale) |
| PUT | `/orders/:id/status` | Update order status |
| DELETE | `/orders/:id` | Cancel order |

**Create Order Request:**
```json
POST /api/v1/orders
{
  "customerId": "uuid (optional)",
  "paymentMethod": "cash",
  "notes": "optional",
  "items": [
    {
      "productId": "uuid",
      "quantity": 2,
      "discount": 0.00
    }
  ]
}
```

### Customers
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/customers` | List customers |
| GET | `/customers/:id` | Get single customer |
| POST | `/customers` | Create customer |
| PUT | `/customers/:id` | Update customer |
| DELETE | `/customers/:id` | Delete customer |

### Reports
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/reports/sales` | Sales summary (daily/weekly/monthly) |
| GET | `/reports/products/top` | Top selling products |
| GET | `/reports/inventory/low` | Low stock alerts |

### Standard Response Envelope
```json
// Success
{
  "data": { ... },
  "meta": { ... }  // Optional: pagination, etc.
}

// Error
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human readable message",
    "details": [
      { "field": "email", "message": "invalid format" }
    ]
  }
}
```

### HTTP Status Codes
| Code | Usage |
|------|-------|
| 200 | Success (GET, PUT) |
| 201 | Created (POST) |
| 204 | No Content (DELETE) |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient role) |
| 404 | Not Found |
| 409 | Conflict (duplicate SKU, etc.) |
| 422 | Unprocessable Entity (business rule violation) |
| 500 | Internal Server Error |

---

## 6. Naming Conventions

### Go (Backend)

| Type | Convention | Example |
|------|------------|---------|
| Packages | lowercase, single word (avoid underscores) | `models`, `repository`, `handler` |
| Files | snake_case | `product_repo.go`, `order_handler.go` |
| Types | PascalCase | `Product`, `OrderStatus` |
| Interfaces | PascalCase with `-er` suffix | `Repository`, `Service` | `Validator` |
| Functions (exported) | PascalCase | `CreateProduct`, `GetOrderByID` |
| Functions (unexported) | camelCase | `validateProduct`, `buildQuery` |
| Variables | camelCase (exported = PascalCase) | `productID`, `OrderTotal` |
| Constants | PascalCase | `StatusPending`, `MaxRetries` |
| Database tables | snake_case, plural | `products`, `order_items` |
| JSON fields | snake_case | `order_number`, `created_at` |

**Special Cases:**
- UUID fields: suffixed with `_id` → `product_id`, `customer_id`
- Timestamps: `_at` suffix → `created_at`, `updated_at`
- Boolean: `is_`, `has_`, `can_` prefix → `is_active`, `has_discount`

### Database Columns
- Tables: plural snake_case (`customers`, `order_items`)
- Columns: snake_case (`order_number`, `postal_code`)
- Indexes: `idx_<table>_<columns>` → `idx_orders_user_id`
- Foreign keys: `<table>_id` → `customer_id`, `product_id`

### Next.js (Frontend)

| Type | Convention | Example |
|------|------------|---------|
| Directories | kebab-case | `order-form/`, `product-card/` |
| Components | PascalCase | `ProductTable.tsx`, `OrderForm.tsx` |
| Pages | kebab-case | `page.tsx`, `new/page.tsx` |
| Hooks | camelCase, `use` prefix | `useProducts.ts`, `useAuth.ts` |
| Utils | camelCase | `formatCurrency.ts`, `cn.ts` |
| Types | PascalCase | `Product`, `OrderStatus` |
| Files | kebab-case | `api-client.ts`, `auth-store.ts` |
| CSS/Tailwind classes | utility-first | `bg-primary text-white` |

### React Component Patterns
```
ComponentName/
├── index.tsx              # Re-export default
├── ComponentName.tsx      # Main component
├── ComponentName.types.ts # Type definitions
└── ComponentName.test.tsx # Tests (if applicable)
```

---

## 7. Layer Responsibilities

### Repository Layer
- Database access only
- No business logic
- Returns domain models
- Transaction management when needed

### Service Layer
- Business logic
- Validation (beyond DTO validation)
- Coordinates multiple repositories
- Returns domain models or DTOs

### Handler Layer (Controllers)
- HTTP request/response handling
- Input extraction and DTO creation
- Calls service layer
- Returns standardized responses
- **No business logic**

### Frontend Patterns
- **Server Components**: Data fetching, layouts, pages
- **Client Components**: `use client`, forms, interactivity
- **React Query**: All API calls, caching, mutations
- **Zustand**: Global UI state only (auth, modals, toasts)

---

## 8. Security Considerations

### Backend
- Passwords: bcrypt with cost 12
- JWT: HS256, 24-hour expiry, refresh token rotation
- Rate limiting: 100 req/min per IP, 1000 req/min per user
- CORS: Configured for specific frontend origin
- SQL injection: GORM parameterized queries
- Input validation: Zod schemas duplicated in Go (go-playground/validator)

### Frontend
- Token storage: httpOnly cookie (preferred) or memory
- XSS: React's default escaping, no `dangerouslySetInnerHTML`
- CSRF: SameSite cookies + CORS

---

## 9. Configuration

### Backend (.env)
```env
# Server
HOST=0.0.0.0
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=saleapp
DB_USER=postgres
DB_PASSWORD=secret
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key-min-32-chars
JWT_EXPIRY_HOURS=24

# Optional
REDIS_URL=redis://localhost:6379
LOG_LEVEL=debug
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

---

## 10. Git Workflow

### Branch Naming
```
feature/<ticket>-short-description
bugfix/<ticket>-short-description
hotfix/<ticket>-short-description
chore/task-description
```

### Commit Messages
```
<type>(<scope>): <subject>

types: feat, fix, docs, style, refactor, test, chore
scope: api, auth, products, orders, customers, ui

Examples:
feat(orders): add order cancellation
fix(auth): resolve token refresh race condition
docs(api): document pagination parameters
```

### PR Requirements
- Title follows commit convention
- Description links to ticket
- All tests pass
- Linting passes
- At least 1 review approval

---

## 11. Testing Strategy

### Backend
- **Unit**: Services with mocked repositories
- **Integration**: Handlers with test database
- **Coverage target**: 70%+ for services

### Frontend
- **Unit**: Utilities, hooks (Vitest)
- **Component**: Critical UI components
- **E2E**: Critical flows (Playwright) - login, create order, void order

---

## 12. Deployment Notes

### Backend
- Single binary deployment
- Environment variables for config
- Graceful shutdown handling
- Health check endpoint: `GET /health`

### Frontend
- Vercel (recommended) or any Node.js host
- ISR for product pages (60-second revalidation)
- Edge runtime for middleware

---

*This document is the source of truth. Any changes must be reviewed by the Architect and updated here first.*
