import '@/styles/globals.css';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Image Uploader',
  description: 'Загружай и просматривай изображения легко и быстро.',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ru">
      <head />
      <body>{children}</body>
    </html>
  );
}
