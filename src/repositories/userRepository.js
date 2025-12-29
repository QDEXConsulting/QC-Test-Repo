/**
 * User Repository
 * Handles database operations for users
 */

import { User } from '../models/user.js';
import { db } from '../config/database.js';

export class UserRepository {
  async findById(id) {
    const query = 'SELECT * FROM users WHERE id = $1';
    const result = await db.query(query, [id]);
    
    if (result.rows.length === 0) {
      return null;
    }
    
    return new User(result.rows[0]);
  }

  async findByEmail(email) {
    const query = 'SELECT * FROM users WHERE email = $1';
    const result = await db.query(query, [email]);
    
    if (result.rows.length === 0) {
      return null;
    }
    
    return new User(result.rows[0]);
  }

  async findAll(filters = {}) {
    let query = 'SELECT * FROM users WHERE 1=1';
    const params = [];
    let paramIndex = 1;

    if (filters.role) {
      query += ` AND role = $${paramIndex++}`;
      params.push(filters.role);
    }

    if (filters.isActive !== undefined) {
      query += ` AND is_active = $${paramIndex++}`;
      params.push(filters.isActive);
    }

    query += ' ORDER BY created_at DESC';
    
    const result = await db.query(query, params);
    return result.rows.map(row => new User(row));
  }

  async create(userData) {
    const query = `
      INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
      VALUES ($1, $2, $3, $4, $5, $6)
      RETURNING *
    `;
    
    const params = [
      userData.email,
      userData.passwordHash,
      userData.firstName,
      userData.lastName,
      userData.role || 'customer',
      userData.isActive ?? true
    ];
    
    const result = await db.query(query, params);
    return new User(result.rows[0]);
  }

  async update(id, userData) {
    const updates = [];
    const params = [];
    let paramIndex = 1;

    if (userData.firstName !== undefined) {
      updates.push(`first_name = $${paramIndex++}`);
      params.push(userData.firstName);
    }

    if (userData.lastName !== undefined) {
      updates.push(`last_name = $${paramIndex++}`);
      params.push(userData.lastName);
    }

    if (userData.role !== undefined) {
      updates.push(`role = $${paramIndex++}`);
      params.push(userData.role);
    }

    if (userData.isActive !== undefined) {
      updates.push(`is_active = $${paramIndex++}`);
      params.push(userData.isActive);
    }

    if (updates.length === 0) {
      return await this.findById(id);
    }

    updates.push(`updated_at = NOW()`);
    params.push(id);

    const query = `
      UPDATE users
      SET ${updates.join(', ')}
      WHERE id = $${paramIndex}
      RETURNING *
    `;
    
    const result = await db.query(query, params);
    return new User(result.rows[0]);
  }

  async delete(id) {
    const query = 'DELETE FROM users WHERE id = $1';
    await db.query(query, [id]);
  }
}

