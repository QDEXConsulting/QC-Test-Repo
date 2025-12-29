/**
 * Order Controller
 * HTTP request handlers for order endpoints
 */

import { OrderService } from '../services/orderService.js';

export class OrderController {
  constructor() {
    this.orderService = new OrderService();
  }

  async getOrder(req, res, next) {
    try {
      const { id } = req.params;
      const order = await this.orderService.getOrderById(id);
      res.json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async getUserOrders(req, res, next) {
    try {
      const { userId } = req.params;
      const orders = await this.orderService.getOrdersByUserId(userId);
      res.json(orders.map(o => o.toJSON()));
    } catch (error) {
      next(error);
    }
  }

  async createOrder(req, res, next) {
    try {
      const orderData = req.body;
      const order = await this.orderService.createOrder(orderData);
      res.status(201).json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async updateOrderStatus(req, res, next) {
    try {
      const { id } = req.params;
      const { status } = req.body;
      const order = await this.orderService.updateOrderStatus(id, status);
      res.json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async cancelOrder(req, res, next) {
    try {
      const { id } = req.params;
      const order = await this.orderService.cancelOrder(id);
      res.json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async deleteOrder(req, res, next) {
    try {
      const { id } = req.params;
      const order = await this.orderService.deleteOrder(id);
      res.json({ message: 'Order deleted successfully', order: order.toJSON() });
    } catch (error) {
      next(error);
    }
  }
}

