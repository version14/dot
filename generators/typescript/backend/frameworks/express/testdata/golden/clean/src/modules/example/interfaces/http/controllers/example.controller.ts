import { Request, Response, NextFunction } from 'express';
import { GetExampleUseCase } from '../../../application/use-cases/get-example.use-case';

const getExample = new GetExampleUseCase();

export class ExampleController {
  async findById(req: Request, res: Response, next: NextFunction): Promise<void> {
    try {
      const item = await getExample.execute(req.params.id);
      res.json(item);
    } catch (err) {
      next(err);
    }
  }
}
