import { Router } from 'express';
import { ExampleController } from '../controllers/example.controller';

const router: Router = Router();
const controller = new ExampleController();

router.get('/:id', controller.findById.bind(controller));

export default router;
