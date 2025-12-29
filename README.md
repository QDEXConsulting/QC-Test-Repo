# Inventory System Test Repository

This is a comprehensive Node.js test repository designed to test relationship-aware file batching systems for LLM analysis. The repository contains a complete inventory management system with clear architectural layers and file relationships.

## Project Structure

```
inventory-system-test-repo/
├── src/
│   ├── models/              # Domain models
│   │   ├── product.js      # Product model
│   │   ├── order.js        # Order model (imports Product)
│   │   ├── user.js         # User model
│   │   └── category.js     # Category model
│   │
│   ├── repositories/       # Data access layer
│   │   ├── productRepository.js   # Product DB operations (imports Product model)
│   │   ├── orderRepository.js     # Order DB operations (imports Order model)
│   │   └── userRepository.js      # User DB operations (imports User model)
│   │
│   ├── services/           # Business logic layer
│   │   ├── productService.js      # Product business logic (imports ProductRepository)
│   │   ├── orderService.js        # Order business logic (imports OrderRepository, ProductService, UserService)
│   │   └── userService.js         # User business logic (imports UserRepository)
│   │
│   ├── controllers/        # HTTP request handlers
│   │   ├── productController.js   # Product endpoints (imports ProductService)
│   │   ├── orderController.js     # Order endpoints (imports OrderService)
│   │   └── userController.js      # User endpoints (imports UserService)
│   │
│   ├── routes/             # Express route definitions
│   │   ├── productRoutes.js      # Product routes (imports ProductController)
│   │   ├── orderRoutes.js        # Order routes (imports OrderController)
│   │   └── userRoutes.js         # User routes (imports UserController)
│   │
│   ├── middleware/         # Express middleware
│   │   ├── auth.js        # Authentication (imports UserService)
│   │   └── errorHandler.js
│   │
│   ├── config/            # Configuration
│   │   └── database.js    # Database connection
│   │
│   └── index.js           # Application entry point (imports routes)
│
├── db/
│   └── migrations/        # Database migrations
│       ├── 001_create_users_table.sql
│       ├── 002_create_products_table.sql
│       ├── 003_create_orders_table.sql
│       └── 004_create_categories_table.sql
│
├── tests/                 # Test files
│   ├── models/
│   │   ├── product.test.js      # Tests Product model
│   │   ├── order.test.js        # Tests Order model (imports Product)
│   │   └── user.test.js         # Tests User model
│   ├── services/
│   │   ├── productService.test.js   # Tests ProductService (imports ProductRepository)
│   │   └── orderService.test.js     # Tests OrderService (imports multiple services)
│   └── controllers/
│       └── productController.test.js # Tests ProductController (imports ProductService)
│
├── package.json
├── jest.config.js
└── README.md
```

## File Relationships

### Strong Relationships (Must Stay Together)

1. **Test Pairing**
   - `tests/models/product.test.js` ↔ `src/models/product.js`
   - `tests/models/order.test.js` ↔ `src/models/order.js`
   - `tests/services/productService.test.js` ↔ `src/services/productService.js`
   - `tests/controllers/productController.test.js` ↔ `src/controllers/productController.js`

2. **Migration Links**
   - `db/migrations/002_create_products_table.sql` → `src/models/product.js`
   - `db/migrations/003_create_orders_table.sql` → `src/models/order.js`
   - `db/migrations/001_create_users_table.sql` → `src/models/user.js`

### Medium Relationships (Should Stay Together)

1. **Layer Dependencies**
   - Controller → Service → Repository → Model
   - Example: `productController.js` → `productService.js` → `productRepository.js` → `product.js`

2. **Model References**
   - `order.js` imports `product.js` (Order contains Product items)
   - `orderRepository.js` references both Order and Product models

3. **Cross-Domain Dependencies**
   - `orderService.js` imports `productService.js` and `userService.js`
   - `orderRepository.js` references users table

### Weak Relationships (Optional Context)

- Shared utilities and configuration files
- Middleware that can be analyzed independently

## Domain Clusters

### Cluster A: Product Domain
- `src/models/product.js`
- `src/repositories/productRepository.js`
- `src/services/productService.js`
- `src/controllers/productController.js`
- `src/routes/productRoutes.js`
- `tests/models/product.test.js`
- `tests/services/productService.test.js`
- `tests/controllers/productController.test.js`
- `db/migrations/002_create_products_table.sql`
- `db/migrations/004_create_categories_table.sql` (related)

### Cluster B: Order Domain
- `src/models/order.js` (imports Product)
- `src/repositories/orderRepository.js`
- `src/services/orderService.js` (imports ProductService, UserService)
- `src/controllers/orderController.js`
- `src/routes/orderRoutes.js`
- `tests/models/order.test.js`
- `tests/services/orderService.test.js`
- `db/migrations/003_create_orders_table.sql`

### Cluster C: User Domain
- `src/models/user.js`
- `src/repositories/userRepository.js`
- `src/services/userService.js`
- `src/controllers/userController.js`
- `src/routes/userRoutes.js`
- `tests/models/user.test.js`
- `db/migrations/001_create_users_table.sql`
- `src/middleware/auth.js` (uses UserService)

## Testing the Batching System

This repository is designed to test relationship-aware batching with the following scenarios:

1. **Import Relationships**: Files that import each other should be batched together
2. **Layer Relationships**: Controller-Service-Repository-Model chains should stay together
3. **Test Pairing**: Test files should be batched with their corresponding source files
4. **Migration Links**: Migrations should be grouped with related model files
5. **Cross-Domain Dependencies**: Order domain depends on Product and User domains
6. **Cluster Boundaries**: Each domain should form a distinct cluster

## Setup

```bash
# Install dependencies
npm install

# Run tests
npm test

# Start server (requires PostgreSQL)
npm start
```

## Expected Batching Behavior

When analyzing changes to this repository, the batching system should:

1. **Identify Clusters**: Group files by domain (Product, Order, User)
2. **Preserve Relationships**: Keep related files together within batches
3. **Respect Limits**: Split clusters into batches when exceeding token/file limits
4. **Maintain Context**: Include anchor files (controllers/models) when splitting
5. **Handle Dependencies**: Include Product and User files when analyzing Order changes

## Use Cases for Testing

1. **Single File Change**: Change to `product.js` should include related files
2. **Cross-Domain Change**: Change to `orderService.js` should include Product and User files
3. **Test Addition**: Adding a test should be batched with the source file
4. **Migration Change**: Migration changes should include related models
5. **Refactoring**: Changes across layers should maintain relationships

This structure provides a realistic codebase with clear relationships that can be used to validate the effectiveness of relationship-aware batching algorithms.

