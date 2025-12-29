/**
 * Order Service
 * Business logic for order operations
 */

import { Order } from '../models/order.js';
import { OrderRepository } from '../repositories/orderRepository.js';
import { ProductService } from './productService.js';
import { UserService } from './userService.js';

export class OrderService {
  constructor() {
    this.orderRepository = new OrderRepository();
    this.productService = new ProductService();
    this.userService = new UserService();
  }

  async getOrderById(id) {
    if (!id) {
      throw new Error('Order ID is required');
    }
    
    const order = await this.orderRepository.findById(id);
    
    if (!order) {
      throw new Error(`Order with ID ${id} not found`);
    }
    
    return order;
  }

  async getOrdersByUserId(userId) {
    if (!userId) {
      throw new Error('User ID is required');
    }
    
    // Verify user exists
    await this.userService.getUserById(userId);
    
    return await this.orderRepository.findByUserId(userId);
  }

  async createOrder(orderData) {
    // Validate user exists
    await this.userService.getUserById(orderData.userId);
    
    // Validate and reserve inventory for each item
    const validatedItems = [];
    
    for (const item of orderData.items) {
      const product = await this.productService.getProductById(item.productId);
      
      if (!product.canFulfillOrder(item.quantity)) {
        throw new Error(
          `Insufficient stock for product ${product.sku}. ` +
          `Requested: ${item.quantity}, Available: ${product.stockQuantity}`
        );
      }
      
      validatedItems.push({
        productId: product.id,
        productSku: product.sku,
        productName: product.name,
        quantity: item.quantity,
        unitPrice: product.price,
        subtotal: product.price * item.quantity
      });
      
      // Reserve inventory
      await this.productService.updateProductStock(product.id, -item.quantity);
    }
    
    // Calculate total
    const totalAmount = validatedItems.reduce((sum, item) => sum + item.subtotal, 0);
    
    // Create order
    const order = await this.orderRepository.create({
      userId: orderData.userId,
      status: 'pending',
      totalAmount: totalAmount,
      items: validatedItems,
      shippingAddress: orderData.shippingAddress
    });
    
    return order;
  }

  async updateOrderStatus(id, status) {
    const validStatuses = ['pending', 'processing', 'shipped', 'delivered', 'cancelled'];
    
    if (!validStatuses.includes(status)) {
      throw new Error(`Invalid order status: ${status}`);
    }
    
    const order = await this.getOrderById(id);
    
    if (status === 'cancelled' && !order.canBeCancelled()) {
      throw new Error(`Order cannot be cancelled. Current status: ${order.status}`);
    }
    
    // If cancelling, restore inventory
    if (status === 'cancelled' && order.status !== 'cancelled') {
      for (const item of order.items) {
        await this.productService.updateProductStock(item.productId, item.quantity);
      }
    }
    
    return await this.orderRepository.updateStatus(id, status);
  }

  async cancelOrder(id) {
    return await this.updateOrderStatus(id, 'cancelled');
  }

  async deleteOrder(id) {
    const order = await this.getOrderById(id);
    
    // Restore inventory if order was not cancelled
    if (order.status !== 'cancelled') {
      for (const item of order.items) {
        await this.productService.updateProductStock(item.productId, item.quantity);
      }
    }
    
    await this.orderRepository.delete(id);
    return order;
  }
}

