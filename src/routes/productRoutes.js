/**
 * Product Routes
 * API routes for product endpoints
 */

import express from 'express';
import { ProductController } from '../controllers/productController.js';

const router = express.Router();
const productController = new ProductController();

router.get('/', productController.listProducts.bind(productController));
router.get('/sku/:sku', productController.getProductBySku.bind(productController));
router.get('/:id', productController.getProduct.bind(productController));
router.post('/', productController.createProduct.bind(productController));
router.put('/:id', productController.updateProduct.bind(productController));
router.patch('/:id/stock', productController.updateProductStock.bind(productController));
router.delete('/:id', productController.deleteProduct.bind(productController));

export default router;

