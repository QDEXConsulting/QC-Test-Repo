/**
 * User Service
 * Business logic for user operations
 */

import { UserRepository } from '../repositories/userRepository.js';
import bcrypt from 'bcrypt';

export class UserService {
  constructor() {
    this.userRepository = new UserRepository();
  }

  async getUserById(id) {
    if (!id) {
      throw new Error('User ID is required');
    }
    
    const user = await this.userRepository.findById(id);
    
    if (!user) {
      throw new Error(`User with ID ${id} not found`);
    }
    
    return user;
  }

  async getUserByEmail(email) {
    if (!email) {
      throw new Error('Email is required');
    }
    
    const user = await this.userRepository.findByEmail(email);
    
    if (!user) {
      throw new Error(`User with email ${email} not found`);
    }
    
    return user;
  }

  async listUsers(filters = {}) {
    return await this.userRepository.findAll(filters);
  }

  async createUser(userData) {
    this.validateUserData(userData);
    
    // Check if email already exists
    const existingUser = await this.userRepository.findByEmail(userData.email);
    if (existingUser) {
      throw new Error(`User with email ${userData.email} already exists`);
    }
    
    // Hash password
    const passwordHash = await bcrypt.hash(userData.password, 10);
    
    return await this.userRepository.create({
      ...userData,
      passwordHash
    });
  }

  async updateUser(id, userData) {
    const existingUser = await this.getUserById(id);
    
    // If email is being changed, check if new email exists
    if (userData.email && userData.email !== existingUser.email) {
      const emailUser = await this.userRepository.findByEmail(userData.email);
      if (emailUser) {
        throw new Error(`User with email ${userData.email} already exists`);
      }
    }
    
    this.validateUserData(userData, true);
    
    return await this.userRepository.update(id, userData);
  }

  async deleteUser(id) {
    const user = await this.getUserById(id);
    await this.userRepository.delete(id);
    return user;
  }

  async authenticateUser(email, password) {
    const user = await this.userRepository.findByEmail(email);
    
    if (!user) {
      throw new Error('Invalid email or password');
    }
    
    if (!user.isActive) {
      throw new Error('User account is inactive');
    }
    
    const isValidPassword = await bcrypt.compare(password, user.passwordHash);
    
    if (!isValidPassword) {
      throw new Error('Invalid email or password');
    }
    
    return user;
  }

  validateUserData(data, isUpdate = false) {
    if (!isUpdate && !data.email) {
      throw new Error('Email is required');
    }
    
    if (data.email && !this.isValidEmail(data.email)) {
      throw new Error('Invalid email format');
    }
    
    if (!isUpdate && !data.password) {
      throw new Error('Password is required');
    }
    
    if (data.password && data.password.length < 8) {
      throw new Error('Password must be at least 8 characters long');
    }
    
    if (!isUpdate && !data.firstName) {
      throw new Error('First name is required');
    }
    
    if (!isUpdate && !data.lastName) {
      throw new Error('Last name is required');
    }
    
    if (data.role && !['admin', 'customer', 'staff'].includes(data.role)) {
      throw new Error('Invalid role. Must be admin, customer, or staff');
    }
  }

  isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }
}

