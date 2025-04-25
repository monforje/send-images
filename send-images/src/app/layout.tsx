// send-images/src/app/layout.tsx

// Импорт глобальных стилей (шрифт, базовая разметка, цвета и т.д.)
import '@/styles/globals.css';
// Типы метаданных для страницы (используется Next.js для <head>)
import type { Metadata } from 'next';

// SEO-метаданные страницы (устанавливаются автоматически в <head>)
export const metadata: Metadata = {
  title: 'Image Uploader', // Заголовок вкладки браузера
};

// Корневой layout — оборачивает все страницы (аналог _app.tsx + _document.tsx в pages/)
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ru"><body>{children}</body></html>
  );
}
