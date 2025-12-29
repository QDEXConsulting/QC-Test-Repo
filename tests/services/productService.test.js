/**
 * Product Service Tests
 * Tests for the ProductService class
 */

import { ProductService } from '../../src/services/productService.js';
import { ProductRepository } from '../../src/repositories/productRepository.js';

// Mock the repository
jest.mock('../../src/repositories/productRepository.js');

describe('ProductService', () => {
  let productService;
  let mockProductRepository;

  beforeEach(() => {
    productService = new ProductService();
    mockProductRepository = productService.productRepository;
  });

  test('should get product by id', async () => {
    const mockProduct = {
      id: 1,
      name: 'Test Product',
      sku: 'TEST-001',
      price: 29.99,
      stockQuantity: 100
    };
    
    mockProductRepository.findById.mockResolvedValue(mockProduct);
    
    const product = await productService.getProductById(1);
    expect(product).toEqual(mockProduct);
    expect(mockProductRepository.findById).toHaveBeenCalledWith(1);
  });

  test('should throw error when product not found', async () => {
    mockProductRepository.findById.mockResolvedValue(null);
    
    await expect(productService.getProductById(999)).rejects.toThrow('Product with ID 999 not found');
  });

  test('should validate product data on create', async () => {
    await expect(productService.createProduct({})).rejects.toThrow('Product name is required');
    await expect(productService.createProduct({ name: 'Test' })).rejects.toThrow('Product SKU is required');
  });
});

