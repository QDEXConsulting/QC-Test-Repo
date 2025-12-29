/**
 * Order Service Tests
 * Tests for the OrderService class
 */

import { OrderService } from '../../src/services/orderService.js';
import { OrderRepository } from '../../src/repositories/orderRepository.js';
import { ProductService } from '../../src/services/productService.js';
import { UserService } from '../../src/services/userService.js';

// Mock dependencies
jest.mock('../../src/repositories/orderRepository.js');
jest.mock('../../src/services/productService.js');
jest.mock('../../src/services/userService.js');

describe('OrderService', () => {
  let orderService;
  let mockOrderRepository;
  let mockProductService;
  let mockUserService;

  beforeEach(() => {
    orderService = new OrderService();
    mockOrderRepository = orderService.orderRepository;
    mockProductService = orderService.productService;
    mockUserService = orderService.userService;
  });

  test('should create order with valid items', async () => {
    const mockUser = { id: 1, email: 'test@example.com' };
    const mockProduct = {
      id: 1,
      sku: 'TEST-001',
      name: 'Test Product',
      price: 29.99,
      stockQuantity: 100,
      canFulfillOrder: jest.fn().mockReturnValue(true)
    };
    const mockOrder = {
      id: 1,
      userId: 1,
      status: 'pending',
      totalAmount: 59.98,
      items: []
    };

    mockUserService.getUserById.mockResolvedValue(mockUser);
    mockProductService.getProductById.mockResolvedValue(mockProduct);
    mockProductService.updateProductStock.mockResolvedValue(mockProduct);
    mockOrderRepository.create.mockResolvedValue(mockOrder);

    const orderData = {
      userId: 1,
      items: [{ productId: 1, quantity: 2 }],
      shippingAddress: { street: '123 Main St' }
    };

    const order = await orderService.createOrder(orderData);
    expect(order).toEqual(mockOrder);
    expect(mockProductService.updateProductStock).toHaveBeenCalled();
  });

  test('should throw error when product stock insufficient', async () => {
    const mockUser = { id: 1 };
    const mockProduct = {
      id: 1,
      sku: 'TEST-001',
      stockQuantity: 5,
      canFulfillOrder: jest.fn().mockReturnValue(false)
    };

    mockUserService.getUserById.mockResolvedValue(mockUser);
    mockProductService.getProductById.mockResolvedValue(mockProduct);

    const orderData = {
      userId: 1,
      items: [{ productId: 1, quantity: 10 }],
      shippingAddress: {}
    };

    await expect(orderService.createOrder(orderData)).rejects.toThrow('Insufficient stock');
  });
});

