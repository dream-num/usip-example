import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ConsoleLogger } from '@nestjs/common';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  await app.listen(8080, () => {
    new ConsoleLogger().log('Server is running on http://localhost:8080');
  });
}
bootstrap();
