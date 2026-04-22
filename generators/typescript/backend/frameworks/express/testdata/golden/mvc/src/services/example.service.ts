import { Example } from '../models/example.model';

export class ExampleService {
  async findAll(): Promise<Example[]> {
    return [];
  }

  async findById(id: string): Promise<Example | null> {
    return { id, name: 'example' };
  }

  async create(data: Partial<Example>): Promise<Example> {
    return { id: crypto.randomUUID(), ...data } as Example;
  }
}
