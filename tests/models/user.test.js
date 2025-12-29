/**
 * User Model Tests
 * Tests for the User model class
 */

import { User } from '../../src/models/user.js';

describe('User Model', () => {
  let userData;

  beforeEach(() => {
    userData = {
      id: 1,
      email: 'test@example.com',
      passwordHash: 'hashed_password',
      firstName: 'John',
      lastName: 'Doe',
      role: 'customer',
      isActive: true,
      createdAt: new Date(),
      updatedAt: new Date()
    };
  });

  test('should create a user instance', () => {
    const user = new User(userData);
    expect(user.id).toBe(1);
    expect(user.email).toBe('test@example.com');
    expect(user.role).toBe('customer');
  });

  test('should get full name', () => {
    const user = new User(userData);
    expect(user.getFullName()).toBe('John Doe');
  });

  test('should check if user is admin', () => {
    const user = new User(userData);
    expect(user.isAdmin()).toBe(false);
    
    user.role = 'admin';
    expect(user.isAdmin()).toBe(true);
  });

  test('should check if user can manage inventory', () => {
    const customer = new User(userData);
    expect(customer.canManageInventory()).toBe(false);
    
    const admin = new User({ ...userData, role: 'admin' });
    expect(admin.canManageInventory()).toBe(true);
    
    const staff = new User({ ...userData, role: 'staff' });
    expect(staff.canManageInventory()).toBe(true);
  });
});

