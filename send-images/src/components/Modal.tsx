// send-images/src/components/Modal.tsx

'use client'; // Компонент работает на клиенте (взаимодействие с событиями, DOM)

// Импорт стилей модального окна
import styles from '@/styles/home.module.css';
import Image from 'next/image';

// Пропсы модального окна: выбранное изображение и функция закрытия
type Props = {
  selected: { filename: string; url: string }; // Данные выбранной картинки
  onClose: () => void; // Функция закрытия модалки
};

// Компонент модального окна
export default function Modal({ selected, onClose }: Props) {
  return (
    // Оверлей, закрывает модалку при клике вне контента
    <div className={styles.modalOverlay} onClick={onClose}>
      {/* Контейнер с изображением — клики по нему не закрывают окно */}
      <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
        <Image
          src={process.env.NEXT_PUBLIC_API_URL + selected.url} // Полный путь до изображения с сервера
          alt={selected.filename} // Alt для доступности
          className={styles.modalImage} // Стили
          loading="lazy" // Ленивая загрузка (для производительности)
        />
        {/* Кнопка закрытия модалки */}
        <button className={styles.modalClose} onClick={onClose}>×</button>
      </div>
    </div>
  );
}
