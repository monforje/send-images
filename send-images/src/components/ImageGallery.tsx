'use client';

import styles from '@/styles/home.module.css';
import Image from 'next/image';
import { MdClose } from 'react-icons/md';

type Props = {
  images: { filename: string; url: string }[];
  onSelect: (img: { filename: string; url: string }) => void;
  onDelete: (filename: string) => void;
};

const API = process.env.NEXT_PUBLIC_API_URL;

export default function ImageGallery({ images, onSelect, onDelete }: Props) {
  const handleDelete = (filename: string) => {
    if (confirm(`Удалить ${filename}?`)) {
      onDelete(filename);
    }
  };

  if (images.length === 0) {
    return (
      <div className={styles.sidebar}>
        <h2>Мои картинки</h2>
        <p>Здесь пока пусто.</p>
      </div>
    );
  }

  return (
    <div className={styles.sidebar}>
      <h2>Мои картинки</h2>
      <div className={styles.thumbs}>
        {images.map((img) => (
          <div key={img.filename} className={styles.thumbContainer}>
            <Image
              src={API + img.url}
              alt={img.filename}
              width={90}
              height={90}
              className={styles.thumb}
              onClick={() => onSelect(img)}
              unoptimized
            />
          <button
            className={styles.deleteBtn}
            onClick={() => handleDelete(img.filename)}
            aria-label="Удалить"
          >
            <MdClose size={16} />
          </button>
          </div>
        ))}
      </div>
    </div>
  );
}
