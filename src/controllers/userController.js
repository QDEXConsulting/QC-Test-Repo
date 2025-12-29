/**
 * User Controller
 * HTTP request handlers for user endpoints
 */

import { UserService } from '../services/userService.js';

export class UserController {
  constructor() {
    this.userService = new UserService();
  }

  async getUser(req, res, next) {
    try {
      const { id } = req.params;
      const user = await this.userService.getUserById(id);
      res.json(user.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async getUserByEmail(req, res, next) {
    try {
      const { email } = req.params;
      const user = await this.userService.getUserByEmail(email);
      res.json(user.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async listUsers(req, res, next) {
    try {
      const filters = {
        role: req.query.role,
        isActive: req.query.isActive === 'true' ? true : req.query.isActive === 'false' ? false : undefined
      };
      
      const users = await this.userService.listUsers(filters);
      res.json(users.map(u => u.toJSON()));
    } catch (error) {
      next(error);
    }
  }

  async createUser(req, res, next) {
    try {
      const userData = req.body;
      const user = await this.userService.createUser(userData);
      res.status(201).json(user.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async updateUser(req, res, next) {
    try {
      const { id } = req.params;
      const userData = req.body;
      const user = await this.userService.updateUser(id, userData);
      res.json(user.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async deleteUser(req, res, next) {
    try {
      const { id } = req.params;
      const user = await this.userService.deleteUser(id);
      res.json({ message: 'User deleted successfully', user: user.toJSON() });
    } catch (error) {
      next(error);
    }
  }

  async login(req, res, next) {
    try {
      const { email, password } = req.body;
      const user = await this.userService.authenticateUser(email, password);
      res.json({ message: 'Login successful', user: user.toJSON() });
    } catch (error) {
      next(error);
    }
  }
}

