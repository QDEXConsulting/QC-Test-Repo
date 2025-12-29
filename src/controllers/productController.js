/**
 * Product Controller
 * HTTP request handlers for product endpoints
 */

import { ProductService } from '../services/productService.js';

export class ProductController {
  constructor() {
    this.productService = new ProductService();
  }

  async getProduct(req, res, next) {
    try {
      const { id } = req.params;
      const product = await this.productService.getProductById(id);
      res.json(product.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async getProductBySku(req, res, next) {
    try {
      const { sku } = req.params;
      const product = await this.productService.getProductBySku(sku);
      res.json(product.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async listProducts(req, res, next) {
    try {
      const filters = {
        categoryId: req.query.categoryId,
        supplierId: req.query.supplierId,
        inStock: req.query.inStock === 'true' ? true : req.query.inStock === 'false' ? false : undefined
      };
      
      const products = await this.productService.listProducts(filters);
      res.json(products.map(p => p.toJSON()));
    } catch (error) {
      next(error);
    }
  }

  async createProduct(req, res, next) {
    try {
      const productData = req.body;
      const product = await this.productService.createProduct(productData);
      res.status(201).json(product.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async updateProduct(req, res, next) {
    try {
      const { id } = req.params;
      const productData = req.body;
      const product = await this.productService.updateProduct(id, productData);
      res.json(product.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async updateProductStock(req, res, next) {
    try {
      const { id } = req.params;
      const { quantity } = req.body;
      const product = await this.productService.updateProductStock(id, quantity);
      res.json(product.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async deleteProduct(req, res, next) {
    try {
      const { id } = req.params;
      const product = await this.productService.deleteProduct(id);
      res.json({ message: 'Product deleted successfully', product: product.toJSON() });
    } catch (error) {
      next(error);
    }
  }
}

