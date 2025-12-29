/**
 * User Routes
 * API routes for user endpoints
 */

import express from 'express';
import { UserController } from '../controllers/userController.js';

const router = express.Router();
const userController = new UserController();

router.post('/login', userController.login.bind(userController));
router.get('/email/:email', userController.getUserByEmail.bind(userController));
router.get('/:id', userController.getUser.bind(userController));
router.get('/', userController.listUsers.bind(userController));
router.post('/', userController.createUser.bind(userController));
router.put('/:id', userController.updateUser.bind(userController));
router.delete('/:id', userController.deleteUser.bind(userController));

export default router;

