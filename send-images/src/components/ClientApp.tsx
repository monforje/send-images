'use client';

import { useState, Suspense, lazy } from 'react';
import Dropzone from './Dropzone';
import ImageGallery from './ImageGallery';
import styles from '@/styles/home.module.css';
import { Image } from '@/lib/api';

const Modal = lazy(() => import('./Modal'));

export default function ClientApp({ initialImages }: { initialImages: Image[] }) {
  const [images, setImages] = useState<Image[]>(initialImages);
  const [selected, setSelected] = useState<Image | null>(null);
  const [saved, setSaved] = useState(false);
  const [deleted, setDeleted] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false); // состояние модалки

  const handleUpload = async (files: File[]) => {
    setLoading(true);
    setError('');

    const formData = new FormData();
    for (const file of files) {
      formData.append('file', file);
    }

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
      const newImages = (data.files ?? []).map((f: { filename: string; url: string }) => ({
        filename: f.filename,
        url: f.url,
      }));

      setImages((prev) => [...prev, ...newImages]);
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
      setDeleted(true);
      setTimeout(() => setDeleted(false), 3000);
      if (selected?.filename === filename) setSelected(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.wrapperColumn}>
      <div className={styles.topBlock}>
        <Dropzone onUpload={handleUpload} />

        {/* Уведомления показывать только если модалка НЕ открыта */}
        {!isModalOpen && (
          <>
            {saved && (
              <div className={styles.toast}>
                <span>Картинка сохранена</span>
              </div>
            )}
            {deleted && (
              <div className={`${styles.toast} ${styles.toastDanger}`}>
                <span>Картинка удалена</span>
              </div>
            )}
            {error && <p className={styles.error}>{error}</p>}
          </>
        )}
      </div>

      <ImageGallery images={images} onSelect={setSelected} onDelete={handleDelete} />

      {selected && (
        <Suspense fallback={null}>
          <Modal
            selected={selected}
            images={images}
            setSelected={setSelected}
            setIsModalOpen={setIsModalOpen} // передаем сюда!
          />
        </Suspense>
      )}
    </div>
  );
}
