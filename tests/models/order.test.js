/**
 * Order Model Tests
 * Tests for the Order model class
 */

import { Order } from '../../src/models/order.js';
import { Product } from '../../src/models/product.js';

describe('Order Model', () => {
  let orderData;
  let product;

  beforeEach(() => {
    product = new Product({
      id: 1,
      name: 'Test Product',
      sku: 'TEST-001',
      price: 29.99,
      stockQuantity: 100
    });

    orderData = {
      id: 1,
      userId: 1,
      status: 'pending',
      totalAmount: 0,
      items: [],
      shippingAddress: {
        street: '123 Main St',
        city: 'Test City',
        zipCode: '12345'
      },
      createdAt: new Date(),
      updatedAt: new Date()
    };
  });

  test('should create an order instance', () => {
    const order = new Order(orderData);
    expect(order.id).toBe(1);
    expect(order.userId).toBe(1);
    expect(order.status).toBe('pending');
  });

  test('should add item to order', () => {
    const order = new Order(orderData);
    order.addItem(product, 5);
    
    expect(order.items.length).toBe(1);
    expect(order.items[0].quantity).toBe(5);
    expect(order.totalAmount).toBe(149.95);
  });

  test('should throw error when adding item with insufficient stock', () => {
    const order = new Order(orderData);
    expect(() => {
      order.addItem(product, 101);
    }).toThrow('Insufficient stock');
  });

  test('should check if order can be cancelled', () => {
    const order = new Order(orderData);
    expect(order.canBeCancelled()).toBe(true);
    
    order.status = 'shipped';
    expect(order.canBeCancelled()).toBe(false);
  });
});

