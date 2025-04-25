'use client';

import { useRef, useState } from 'react';
import styles from '@/styles/home.module.css';

type Props = {
  onUpload: (file: File) => void;
};

const MAX_SIZE_MB = 5;

export default function Dropzone({ onUpload }: Props) {
  const [dragOver, setDragOver] = useState(false);
  const [error, setError] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFiles = (files: FileList) => {
    Array.from(files).forEach((file) => {
      if (!file.type.startsWith('image/')) {
        setError('Можно загружать только изображения.');
        return;
      }
      if (file.size > MAX_SIZE_MB * 1024 * 1024) {
        setError(`Файл слишком большой. Макс ${MAX_SIZE_MB}MB.`);
        return;
      }
      setError('');
      onUpload(file);
    });
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    if (e.dataTransfer.files) {
      handleFiles(e.dataTransfer.files);
    }
  };

  const handleClick = () => inputRef.current?.click();

  return (
    <>
      <div
        className={`${styles.square} ${dragOver ? styles.dragOver : ''}`}
        onClick={handleClick}
        onDragOver={(e) => {
          e.preventDefault();
          setDragOver(true);
        }}
        onDragLeave={() => setDragOver(false)}
        onDrop={handleDrop}
      >
        <span>Кликни или перетащи сюда</span>
      </div>

      <input
        type="file"
        multiple
        accept="image/*"
        ref={inputRef}
        onChange={(e) => {
          if (e.target.files) handleFiles(e.target.files);
        }}
        hidden
      />

      {error && <p className={styles.error}>{error}</p>}
    </>
  );
}
