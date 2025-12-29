/**
 * Product Service
 * Business logic for product operations
 */

import { ProductRepository } from '../repositories/productRepository.js';
import { Category } from '../models/category.js';

export class ProductService {
  constructor() {
    this.productRepository = new ProductRepository();
  }

  async getProductById(id) {
    if (!id) {
      throw new Error('Product ID is required');
    }
    
    const product = await this.productRepository.findById(id);
    
    if (!product) {
      throw new Error(`Product with ID ${id} not found`);
    }
    
    return product;
  }

  async getProductBySku(sku) {
    if (!sku) {
      throw new Error('Product SKU is required');
    }
    
    const product = await this.productRepository.findBySku(sku);
    
    if (!product) {
      throw new Error(`Product with SKU ${sku} not found`);
    }
    
    return product;
  }

  async listProducts(filters = {}) {
    return await this.productRepository.findAll(filters);
  }

  async createProduct(productData) {
    this.validateProductData(productData);
    
    // Check if SKU already exists
    const existingProduct = await this.productRepository.findBySku(productData.sku);
    if (existingProduct) {
      throw new Error(`Product with SKU ${productData.sku} already exists`);
    }
    
    return await this.productRepository.create(productData);
  }

  async updateProduct(id, productData) {
    const existingProduct = await this.getProductById(id);
    
    this.validateProductData(productData, true);
    
    return await this.productRepository.update(id, productData);
  }

  async updateProductStock(id, quantity) {
    const product = await this.getProductById(id);
    
    const newQuantity = product.stockQuantity + quantity;
    
    if (newQuantity < 0) {
      throw new Error(`Insufficient stock. Current stock: ${product.stockQuantity}`);
    }
    
    return await this.productRepository.updateStock(id, quantity);
  }

  async deleteProduct(id) {
    const product = await this.getProductById(id);
    await this.productRepository.delete(id);
    return product;
  }

  validateProductData(data, isUpdate = false) {
    if (!isUpdate && !data.name) {
      throw new Error('Product name is required');
    }
    
    if (!isUpdate && !data.sku) {
      throw new Error('Product SKU is required');
    }
    
    if (data.price !== undefined && data.price < 0) {
      throw new Error('Product price cannot be negative');
    }
    
    if (data.stockQuantity !== undefined && data.stockQuantity < 0) {
      throw new Error('Stock quantity cannot be negative');
    }
  }
}

