'use client';

import { useEffect, useState, Suspense, lazy } from 'react';
import Dropzone from './Dropzone';
import ImageGallery from './ImageGallery';
import styles from '@/styles/home.module.css';

const Modal = lazy(() => import('./Modal'));

type Image = { filename: string; url: string };

export default function ClientApp({ initialImages }: { initialImages: Image[] }) {
  const [images, setImages] = useState(initialImages);
  const [selected, setSelected] = useState<Image | null>(null);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleUpload = async (file: File) => {
    setLoading(true);
    setError('');
    const formData = new FormData();
    formData.append('file', file);

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/upload`, {
        method: 'POST',
        body: formData,
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || 'Ошибка при загрузке');
      }

      const data = await res.json();
      setImages((prev) => [...prev, { filename: data.filename, url: data.url }]);
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (filename: string) => {
    setLoading(true);
    setError('');
    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/images?filename=${encodeURIComponent(filename)}`,
        { method: 'DELETE' }
      );
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || 'Не удалось удалить');
      }

      setImages((prev) => prev.filter((img) => img.filename !== filename));
      if (selected?.filename === filename) setSelected(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.wrapper}>
      <div className={styles.left}>
        <p className={styles.title}>Кинь картинку в прямоугольник!</p>
        <Dropzone onUpload={handleUpload} />
        {saved && <p className={styles.saved}>Картинка сохранена ✅</p>}
        {error && <p className={styles.error}>{error}</p>}
        {/* Удалена строка загрузки */}
      </div>

      <ImageGallery images={images} onSelect={setSelected} onDelete={handleDelete} />

      {selected && (
        <Suspense fallback={null}>
          <Modal
            selected={selected}
            images={images}
            setSelected={setSelected}
          />
        </Suspense>
      )}
    </div>
  );
}
