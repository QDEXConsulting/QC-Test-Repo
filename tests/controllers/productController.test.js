/**
 * Product Controller Tests
 * Tests for the ProductController class
 */

import { ProductController } from '../../src/controllers/productController.js';
import { ProductService } from '../../src/services/productService.js';

jest.mock('../../src/services/productService.js');

describe('ProductController', () => {
  let productController;
  let mockProductService;
  let req, res, next;

  beforeEach(() => {
    productController = new ProductController();
    mockProductService = productController.productService;
    
    req = {
      params: {},
      query: {},
      body: {}
    };
    
    res = {
      json: jest.fn(),
      status: jest.fn().mockReturnThis()
    };
    
    next = jest.fn();
  });

  test('should get product by id', async () => {
    const mockProduct = {
      id: 1,
      name: 'Test Product',
      toJSON: jest.fn().mockReturnValue({ id: 1, name: 'Test Product' })
    };
    
    req.params.id = '1';
    mockProductService.getProductById.mockResolvedValue(mockProduct);
    
    await productController.getProduct(req, res, next);
    
    expect(mockProductService.getProductById).toHaveBeenCalledWith('1');
    expect(res.json).toHaveBeenCalled();
  });

  test('should handle errors', async () => {
    req.params.id = '1';
    const error = new Error('Product not found');
    mockProductService.getProductById.mockRejectedValue(error);
    
    await productController.getProduct(req, res, next);
    
    expect(next).toHaveBeenCalledWith(error);
  });
});

