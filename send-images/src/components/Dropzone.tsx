// send-images/src/components/Dropzone.tsx

'use client'; // Компонент работает на клиенте (используется drag'n'drop, взаимодействие с DOM)

// Импорт стилей из CSS-модуля
import styles from '@/styles/home.module.css';

// Тип пропсов: при загрузке вызывается функция onUpload с выбранным файлом
type Props = {
  onUpload: (file: File) => void;
};

// Компонент Dropzone — область, в которую можно перетаскивать изображение для загрузки
export default function Dropzone({ onUpload }: Props) {
  // Обработчик события drop
  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault(); // Отменяем стандартное поведение (иначе браузер может открыть файл)

    // Получаем первый файл из DataTransfer
    const file = e.dataTransfer.files?.[0];

    // Проверяем, что это изображение
    if (file && file.type.startsWith('image/')) {
      onUpload(file); // Передаём файл в родительский компонент
    }
  };

  // Область с рамкой, принимает события drag-over и drop
  return (
    <div
      className={styles.square}
      onDragOver={(e) => e.preventDefault()} // Не даём браузеру блокировать drop
      onDrop={handleDrop} // Обрабатываем drop
    >
      <span>Перетащи сюда</span>
    </div>
  );
}
