/**
 * Product Model
 * Represents a product in the inventory system
 */

export class Product {
  constructor(data) {
    this.id = data.id;
    this.name = data.name;
    this.sku = data.sku;
    this.description = data.description;
    this.price = data.price;
    this.stockQuantity = data.stockQuantity;
    this.categoryId = data.categoryId;
    this.supplierId = data.supplierId;
    this.createdAt = data.createdAt;
    this.updatedAt = data.updatedAt;
  }

  toJSON() {
    return {
      id: this.id,
      name: this.name,
      sku: this.sku,
      description: this.description,
      price: this.price,
      stockQuantity: this.stockQuantity,
      categoryId: this.categoryId,
      supplierId: this.supplierId,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt
    };
  }

  isInStock() {
    return this.stockQuantity > 0;
  }

  canFulfillOrder(quantity) {
    return this.stockQuantity >= quantity;
  }
}

