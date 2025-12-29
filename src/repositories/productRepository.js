/**
 * Product Repository
 * Handles database operations for products
 */

import { Product } from '../models/product.js';
import { db } from '../config/database.js';

export class ProductRepository {
  async findById(id) {
    const query = 'SELECT * FROM products WHERE id = $1';
    const result = await db.query(query, [id]);
    
    if (result.rows.length === 0) {
      return null;
    }
    
    return new Product(result.rows[0]);
  }

  async findBySku(sku) {
    const query = 'SELECT * FROM products WHERE sku = $1';
    const result = await db.query(query, [sku]);
    
    if (result.rows.length === 0) {
      return null;
    }
    
    return new Product(result.rows[0]);
  }

  async findAll(filters = {}) {
    let query = 'SELECT * FROM products WHERE 1=1';
    const params = [];
    let paramIndex = 1;

    if (filters.categoryId) {
      query += ` AND category_id = $${paramIndex++}`;
      params.push(filters.categoryId);
    }

    if (filters.supplierId) {
      query += ` AND supplier_id = $${paramIndex++}`;
      params.push(filters.supplierId);
    }

    if (filters.inStock !== undefined) {
      if (filters.inStock) {
        query += ` AND stock_quantity > 0`;
      } else {
        query += ` AND stock_quantity = 0`;
      }
    }

    query += ' ORDER BY created_at DESC';
    
    const result = await db.query(query, params);
    return result.rows.map(row => new Product(row));
  }

  async create(productData) {
    const query = `
      INSERT INTO products (name, sku, description, price, stock_quantity, category_id, supplier_id)
      VALUES ($1, $2, $3, $4, $5, $6, $7)
      RETURNING *
    `;
    
    const params = [
      productData.name,
      productData.sku,
      productData.description,
      productData.price,
      productData.stockQuantity,
      productData.categoryId,
      productData.supplierId
    ];
    
    const result = await db.query(query, params);
    return new Product(result.rows[0]);
  }

  async update(id, productData) {
    const query = `
      UPDATE products
      SET name = $1, description = $2, price = $3, stock_quantity = $4,
          category_id = $5, updated_at = NOW()
      WHERE id = $6
      RETURNING *
    `;
    
    const params = [
      productData.name,
      productData.description,
      productData.price,
      productData.stockQuantity,
      productData.categoryId,
      id
    ];
    
    const result = await db.query(query, params);
    return new Product(result.rows[0]);
  }

  async updateStock(id, quantity) {
    const query = `
      UPDATE products
      SET stock_quantity = stock_quantity + $1, updated_at = NOW()
      WHERE id = $2
      RETURNING *
    `;
    
    const result = await db.query(query, [quantity, id]);
    return new Product(result.rows[0]);
  }

  async delete(id) {
    const query = 'DELETE FROM products WHERE id = $1';
    await db.query(query, [id]);
  }
}

