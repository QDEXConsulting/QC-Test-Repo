# File Relationship Mapping

This document maps out the relationships between files in this test repository. These relationships should be preserved when batching files for LLM analysis.

## Relationship Types

- **Strong (Weight: 3)**: Must stay together if possible
- **Medium (Weight: 2)**: Should stay together
- **Weak (Weight: 1)**: Optional context

## Strong Relationships

### Test Pairing
```
tests/models/product.test.js → src/models/product.js (Strong, "test pairing")
tests/models/order.test.js → src/models/order.js (Strong, "test pairing")
tests/models/user.test.js → src/models/user.js (Strong, "test pairing")
tests/services/productService.test.js → src/services/productService.js (Strong, "test pairing")
tests/services/orderService.test.js → src/services/orderService.js (Strong, "test pairing")
tests/controllers/productController.test.js → src/controllers/productController.js (Strong, "test pairing")
```

### Migration Links
```
db/migrations/001_create_users_table.sql → src/models/user.js (Strong, "schema definition")
db/migrations/002_create_products_table.sql → src/models/product.js (Strong, "schema definition")
db/migrations/003_create_orders_table.sql → src/models/order.js (Strong, "schema definition")
db/migrations/004_create_categories_table.sql → src/models/category.js (Strong, "schema definition")
```

## Medium Relationships

### Layer Dependencies (Controller → Service → Repository → Model)

#### Product Domain
```
src/controllers/productController.js → src/services/productService.js (Medium, "layer dependency")
src/services/productService.js → src/repositories/productRepository.js (Medium, "layer dependency")
src/repositories/productRepository.js → src/models/product.js (Medium, "layer dependency")
```

#### Order Domain
```
src/controllers/orderController.js → src/services/orderService.js (Medium, "layer dependency")
src/services/orderService.js → src/repositories/orderRepository.js (Medium, "layer dependency")
src/repositories/orderRepository.js → src/models/order.js (Medium, "layer dependency")
```

#### User Domain
```
src/controllers/userController.js → src/services/userService.js (Medium, "layer dependency")
src/services/userService.js → src/repositories/userRepository.js (Medium, "layer dependency")
src/repositories/userRepository.js → src/models/user.js (Medium, "layer dependency")
```

### Route Definitions
```
src/routes/productRoutes.js → src/controllers/productController.js (Medium, "route binding")
src/routes/orderRoutes.js → src/controllers/orderController.js (Medium, "route binding")
src/routes/userRoutes.js → src/controllers/userController.js (Medium, "route binding")
```

### Model References
```
src/models/order.js → src/models/product.js (Medium, "model reference - Order contains Products")
src/repositories/orderRepository.js → src/models/product.js (Medium, "model reference")
```

### Cross-Domain Dependencies
```
src/services/orderService.js → src/services/productService.js (Medium, "cross-domain dependency")
src/services/orderService.js → src/services/userService.js (Medium, "cross-domain dependency")
src/repositories/productRepository.js → src/models/user.js (Medium, "supplier reference")
src/repositories/orderRepository.js → src/models/user.js (Medium, "user reference")
```

### Middleware Dependencies
```
src/middleware/auth.js → src/services/userService.js (Medium, "authentication dependency")
```

## Weak Relationships

### Application Setup
```
src/index.js → src/routes/productRoutes.js (Weak, "application setup")
src/index.js → src/routes/orderRoutes.js (Weak, "application setup")
src/index.js → src/routes/userRoutes.js (Weak, "application setup")
src/index.js → src/middleware/errorHandler.js (Weak, "application setup")
```

### Configuration
```
src/repositories/productRepository.js → src/config/database.js (Weak, "configuration")
src/repositories/orderRepository.js → src/config/database.js (Weak, "configuration")
src/repositories/userRepository.js → src/config/database.js (Weak, "configuration")
```

## Domain Clusters

### Cluster A: Product Domain
**Files:**
- `src/models/product.js`
- `src/models/category.js` (related)
- `src/repositories/productRepository.js`
- `src/services/productService.js`
- `src/controllers/productController.js`
- `src/routes/productRoutes.js`
- `tests/models/product.test.js`
- `tests/services/productService.test.js`
- `tests/controllers/productController.test.js`
- `db/migrations/002_create_products_table.sql`
- `db/migrations/004_create_categories_table.sql`

**Total Files:** 11
**Has Changed Files:** Yes (all files in this cluster)

### Cluster B: Order Domain
**Files:**
- `src/models/order.js` (imports Product)
- `src/repositories/orderRepository.js`
- `src/services/orderService.js` (imports ProductService, UserService)
- `src/controllers/orderController.js`
- `src/routes/orderRoutes.js`
- `tests/models/order.test.js`
- `tests/services/orderService.test.js`
- `db/migrations/003_create_orders_table.sql`

**Dependencies:**
- Requires Product domain files (for Product model/service)
- Requires User domain files (for User service)

**Total Files:** 8 (plus dependencies)
**Has Changed Files:** Yes

### Cluster C: User Domain
**Files:**
- `src/models/user.js`
- `src/repositories/userRepository.js`
- `src/services/userService.js`
- `src/controllers/userController.js`
- `src/routes/userRoutes.js`
- `tests/models/user.test.js`
- `src/middleware/auth.js` (uses UserService)
- `db/migrations/001_create_users_table.sql`

**Total Files:** 8
**Has Changed Files:** Yes

## Example Batching Scenarios

### Scenario 1: Single Product File Change
**Changed File:** `src/models/product.js`

**Expected Batch:**
- `src/models/product.js` (changed)
- `src/repositories/productRepository.js` (depends on Product)
- `src/services/productService.js` (depends on ProductRepository)
- `src/controllers/productController.js` (depends on ProductService)
- `tests/models/product.test.js` (test pairing)

### Scenario 2: Order Service Change
**Changed File:** `src/services/orderService.js`

**Expected Batch:**
- `src/services/orderService.js` (changed)
- `src/repositories/orderRepository.js` (used by OrderService)
- `src/models/order.js` (used by OrderRepository)
- `src/services/productService.js` (imported by OrderService)
- `src/services/userService.js` (imported by OrderService)
- `src/models/product.js` (used by ProductService)
- `src/models/user.js` (used by UserService)
- `tests/services/orderService.test.js` (test pairing)

### Scenario 3: Migration Change
**Changed File:** `db/migrations/002_create_products_table.sql`

**Expected Batch:**
- `db/migrations/002_create_products_table.sql` (changed)
- `src/models/product.js` (schema definition)
- `src/repositories/productRepository.js` (uses Product model)

### Scenario 4: Cross-Domain Refactoring
**Changed Files:** 
- `src/models/product.js`
- `src/services/orderService.js`

**Expected Clusters:**
- **Cluster A (Product):** All Product domain files
- **Cluster B (Order):** All Order domain files + Product domain dependencies

**Batching Strategy:**
- Batch 1: Product domain core (models, repositories, services)
- Batch 2: Product domain controllers/routes + tests
- Batch 3: Order domain with Product context summary

