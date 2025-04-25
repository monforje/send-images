// send-images/src/components/ImageGallery.tsx

'use client'; // Компонент работает на клиенте (есть обработчики событий)

// Импорт стилей из CSS-модуля
import styles from '@/styles/home.module.css';
// next/image используется для оптимизации изображений (ленивая загрузка, ресайз и т.д.)
import Image from 'next/image';

// Пропсы: массив изображений + коллбэки на выбор и удаление
type Props = {
  images: { filename: string; url: string }[]; // Список картинок
  onSelect: (img: { filename: string; url: string }) => void; // Клик по картинке
  onDelete: (filename: string) => void; // Клик по кнопке удаления
};

// Компонент галереи изображений
export default function ImageGallery({ images, onSelect, onDelete }: Props) {
  // Базовый адрес API — используется как префикс к путям картинок
  const API = process.env.NEXT_PUBLIC_API_URL;

  return (
    <div className={styles.sidebar}> {/* Правая колонка со списком */}
      <h2>Мои картинки</h2>
      <div className={styles.thumbs}> {/* Сетка превьюшек */}
        {images.map((img) => (
          <div key={img.filename} className={styles.thumbContainer}> {/* Обёртка превью + кнопка удаления */}
            <Image
              src={API + img.url} // Полный путь до изображения
              alt={img.filename} // Альтернативный текст
              width={80} // Размеры превью
              height={80}
              className={styles.thumb} // Стилизация
              unoptimized // Отключаем оптимизацию (т.к. изображения с внешнего сервера)
              onClick={() => onSelect(img)} // Открыть в модалке
            />
            <button
              className={styles.deleteBtn}
              onClick={() => onDelete(img.filename)} // Удалить изображение
            >
              ×
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
