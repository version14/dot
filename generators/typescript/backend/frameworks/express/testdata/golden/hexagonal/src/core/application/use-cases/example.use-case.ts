import { ExampleInputPort } from '../ports/in/example.port';
import { Example } from '../../domain/entities/example.entity';

export class ExampleUseCase implements ExampleInputPort {
  async findById(id: string): Promise<Example | null> {
    // TODO: inject output port (repository) via constructor
    return { id, name: 'example' };
  }
}
