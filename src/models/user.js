/**
 * User Model
 * Represents a user in the system
 */

export class User {
  constructor(data) {
    this.id = data.id;
    this.email = data.email;
    this.passwordHash = data.passwordHash;
    this.firstName = data.firstName;
    this.lastName = data.lastName;
    this.role = data.role; // 'admin', 'customer', 'staff'
    this.isActive = data.isActive ?? true;
    this.createdAt = data.createdAt;
    this.updatedAt = data.updatedAt;
  }

  toJSON() {
    return {
      id: this.id,
      email: this.email,
      firstName: this.firstName,
      lastName: this.lastName,
      role: this.role,
      isActive: this.isActive,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt
    };
  }

  getFullName() {
    return `${this.firstName} ${this.lastName}`;
  }

  isAdmin() {
    return this.role === 'admin';
  }

  canManageInventory() {
    return this.role === 'admin' || this.role === 'staff';
  }
}

