// send-images/src/components/ClientApp.tsx

'use client'; // Обозначаем, что компонент выполняется на клиенте (нужен доступ к useState, эффектам и т.д.)

import { useState, Suspense, lazy } from 'react';
import Dropzone from './Dropzone'; // Компонент для drag & drop загрузки
import ImageGallery from './ImageGallery'; // Галерея изображений
import styles from '@/styles/home.module.css'; // Стили

// Ленивая загрузка модального окна — будет загружен только при необходимости
const Modal = lazy(() => import('./Modal'));

// Тип для изображения
type Image = { filename: string; url: string };

// Главный клиентский компонент, получает список изображений с сервера и управляет состоянием UI
export default function ClientApp({ initialImages }: { initialImages: Image[] }) {
  const [images, setImages] = useState(initialImages); // Список изображений
  const [selected, setSelected] = useState<Image | null>(null); // Выбранное изображение для модалки
  const [saved, setSaved] = useState(false); // Флаг, показывающий, что файл успешно загружен

  // Обработчик загрузки нового изображения
  const handleUpload = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file); // Добавляем файл в тело запроса
    
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/upload`, {
      method: 'POST',
      body: formData,
    });
    
    if (res.ok) {
      const data = await res.json(); // Получаем имя и URL сохранённого файла
      setImages((prev) => [...prev, { filename: data.filename, url: data.url }]); // Добавляем в список
      setSaved(true); // Показываем сообщение "сохранено"
    }
  };

  // Обработчик удаления изображения
  const handleDelete = async (filename: string) => {
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/images?filename=${encodeURIComponent(filename)}`,
      { method: 'DELETE' }
    );
    if (res.status === 204) {
      setImages((prev) => prev.filter((img) => img.filename !== filename)); // Удаляем из локального списка
      if (selected?.filename === filename) setSelected(null); // Закрываем модалку, если выбрано удалённое изображение
    }
  };

  return (
    <div className={styles.wrapper}>
      <div className={styles.left}>
        <p className={styles.title}>Кинь картинку в прямоугольник!</p>
        <Dropzone onUpload={handleUpload} /> {/* DnD загрузка */}
        {saved && <p className={styles.saved}>Картинка сохранена ✅</p>}
      </div>

      {/* Отображаем превьюшки и кнопки удаления */}
      <ImageGallery images={images} onSelect={setSelected} onDelete={handleDelete} />

      {/* Модальное окно — отображается только при выбранной картинке */}
      {selected && (
        <Suspense fallback={null}> {/* Lazy loading через Suspense */}
          <Modal selected={selected} onClose={() => setSelected(null)} />
        </Suspense>
      )}
    </div>
  );
}
