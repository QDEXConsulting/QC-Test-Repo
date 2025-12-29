/**
 * Order Model
 * Represents an order in the system
 */

import { Product } from './product.js';

export class Order {
  constructor(data) {
    this.id = data.id;
    this.userId = data.userId;
    this.status = data.status; // 'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    this.totalAmount = data.totalAmount;
    this.items = data.items || [];
    this.shippingAddress = data.shippingAddress;
    this.createdAt = data.createdAt;
    this.updatedAt = data.updatedAt;
  }

  toJSON() {
    return {
      id: this.id,
      userId: this.userId,
      status: this.status,
      totalAmount: this.totalAmount,
      items: this.items,
      shippingAddress: this.shippingAddress,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt
    };
  }

  addItem(product, quantity) {
    if (!(product instanceof Product)) {
      throw new Error('Item must be a Product instance');
    }
    
    if (!product.canFulfillOrder(quantity)) {
      throw new Error(`Insufficient stock for product ${product.sku}`);
    }

    this.items.push({
      productId: product.id,
      productSku: product.sku,
      productName: product.name,
      quantity: quantity,
      unitPrice: product.price,
      subtotal: product.price * quantity
    });

    this.calculateTotal();
  }

  calculateTotal() {
    this.totalAmount = this.items.reduce((sum, item) => sum + item.subtotal, 0);
  }

  canBeCancelled() {
    return this.status === 'pending' || this.status === 'processing';
  }
}

