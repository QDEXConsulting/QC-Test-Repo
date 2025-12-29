/**
 * Product Model Tests
 * Tests for the Product model class
 */

import { Product } from '../../src/models/product.js';

describe('Product Model', () => {
  let productData;

  beforeEach(() => {
    productData = {
      id: 1,
      name: 'Test Product',
      sku: 'TEST-001',
      description: 'A test product',
      price: 29.99,
      stockQuantity: 100,
      categoryId: 1,
      supplierId: 1,
      createdAt: new Date(),
      updatedAt: new Date()
    };
  });

  test('should create a product instance', () => {
    const product = new Product(productData);
    expect(product.id).toBe(1);
    expect(product.name).toBe('Test Product');
    expect(product.sku).toBe('TEST-001');
  });

  test('should convert to JSON', () => {
    const product = new Product(productData);
    const json = product.toJSON();
    
    expect(json).toHaveProperty('id');
    expect(json).toHaveProperty('name');
    expect(json).toHaveProperty('sku');
    expect(json).toHaveProperty('price');
  });

  test('should check if product is in stock', () => {
    const product = new Product(productData);
    expect(product.isInStock()).toBe(true);
    
    product.stockQuantity = 0;
    expect(product.isInStock()).toBe(false);
  });

  test('should check if product can fulfill order', () => {
    const product = new Product(productData);
    expect(product.canFulfillOrder(50)).toBe(true);
    expect(product.canFulfillOrder(100)).toBe(true);
    expect(product.canFulfillOrder(101)).toBe(false);
  });
});

