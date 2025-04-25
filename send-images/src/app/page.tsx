// send-images/src/app/page.tsx

// Импорт функции получения изображений с Go-бэкенда
import { getImages } from '@/lib/api';
// Импорт client-side компонента, куда мы передаём изображения
import ClientApp from '@/components/ClientApp';

// Серверный компонент страницы (Next.js 14 App Router)
// Выполняется на сервере перед отрисовкой
export default async function Page() {
  // Получаем список изображений с бэкенда заранее (SSR)
  const images = await getImages(); // ✅ Выполняется до рендера клиента

  // Передаём изображения клиентскому компоненту (useState, UI и т.д.)
  return <ClientApp initialImages={images} />;
}
