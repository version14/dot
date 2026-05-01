import { Router } from 'express';
import { login, logout, me, refresh, register } from '../modules/auth/application/controllers/auth.controller';
import { authMiddleware } from '../shared/middlewares/auth.middleware';

const router = Router();

router.post('/register', register);
router.post('/login', login);
router.get('/me', authMiddleware, me);
router.post('/refresh', refresh);
router.post('/logout', logout);

export default router;
