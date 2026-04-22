import { Request, Response, NextFunction } from 'express';
import { ExampleService } from '../services/example.service';

const service = new ExampleService();

export class ExampleController {
  async findAll(_req: Request, res: Response, next: NextFunction): Promise<void> {
    try {
      const items = await service.findAll();
      res.json(items);
    } catch (err) {
      next(err);
    }
  }

  async findById(req: Request, res: Response, next: NextFunction): Promise<void> {
    try {
      const item = await service.findById(req.params.id);
      res.json(item);
    } catch (err) {
      next(err);
    }
  }

  async create(req: Request, res: Response, next: NextFunction): Promise<void> {
    try {
      const item = await service.create(req.body);
      res.status(201).json(item);
    } catch (err) {
      next(err);
    }
  }
}
