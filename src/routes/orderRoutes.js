/**
 * Order Routes
 * API routes for order endpoints
 */

import express from 'express';
import { OrderController } from '../controllers/orderController.js';

const router = express.Router();
const orderController = new OrderController();

router.get('/user/:userId', orderController.getUserOrders.bind(orderController));
router.get('/:id', orderController.getOrder.bind(orderController));
router.post('/', orderController.createOrder.bind(orderController));
router.patch('/:id/status', orderController.updateOrderStatus.bind(orderController));
router.post('/:id/cancel', orderController.cancelOrder.bind(orderController));
router.delete('/:id', orderController.deleteOrder.bind(orderController));

export default router;

