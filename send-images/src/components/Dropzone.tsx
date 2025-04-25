'use client';

import { useRef, useState } from 'react';
import styles from '@/styles/home.module.css';
import { FiUploadCloud } from 'react-icons/fi';

type Props = {
  onUpload: (files: File[]) => void;
  maxSizeMb?: number;
  accept?: string;
  multiple?: boolean;
};

const defaultMaxSize = 5;

export default function Dropzone({
  onUpload,
  maxSizeMb = defaultMaxSize,
  accept = 'image/*',
  multiple = true,
}: Props) {
  const [dragOver, setDragOver] = useState(false);
  const [error, setError] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFiles = (files: FileList) => {
    const validFiles: File[] = [];

    for (const file of Array.from(files)) {
      if (!file.type.startsWith('image/')) {
        setError('Можно загружать только изображения.');
        return;
      }
      if (file.size > maxSizeMb * 1024 * 1024) {
        setError(`Файл слишком большой. Макс ${maxSizeMb}MB.`);
        return;
      }
      validFiles.push(file);
    }

    if (validFiles.length) {
      setError('');
      onUpload(validFiles);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    if (e.dataTransfer.files) handleFiles(e.dataTransfer.files);
  };

  return (
    <>
      <div
        className={`${styles.square} ${dragOver ? styles.dragOver : ''}`}
        onClick={() => inputRef.current?.click()}
        onDragOver={(e) => {
          e.preventDefault();
          setDragOver(true);
        }}
        onDragLeave={() => setDragOver(false)}
        onDrop={handleDrop}
      >
        <FiUploadCloud size={60} color="#42a5f5" />
        <span className={styles.squareLabel}>Кликни или перетащи файлы</span>
      </div>

      <input
        type="file"
        ref={inputRef}
        accept={accept}
        multiple={multiple}
        onChange={(e) => {
          if (e.target.files) handleFiles(e.target.files);
        }}
        hidden
      />

      {error && <p className={styles.error}>{error}</p>}
    </>
  );
}
