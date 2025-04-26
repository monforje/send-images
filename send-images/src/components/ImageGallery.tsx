'use client';

import styles from '@/styles/home.module.css';
import Image from 'next/image';
import { MdClose } from 'react-icons/md';
import { useInView } from 'react-intersection-observer';

type ImageType = {
  filename: string;
  url: string;
};

type Props = {
  images: ImageType[];
  onSelect: (img: ImageType) => void;
  onDelete: (filename: string) => void;
};

const API = process.env.NEXT_PUBLIC_API_URL;

export default function ImageGallery({ images, onSelect, onDelete }: Props) {
  return (
    <div className={styles.galleryWrapper}>
      <h2 className={styles.galleryTitle}>Мои картинки</h2>

      {images.length === 0 ? (
        <p className={styles.empty}>Здесь пока пусто 😢</p>
      ) : (
        <div className={styles.thumbs}>
          {images.map((img, index) => (
            <ImageItem
              key={img.filename}
              img={img}
              onSelect={onSelect}
              onDelete={onDelete}
              priority={index < 9} // <-- ВАЖНО
            />
          ))}
        </div>
      )}
    </div>
  );
}

type ImageItemProps = {
  img: ImageType;
  onSelect: (img: ImageType) => void;
  onDelete: (filename: string) => void;
  priority: boolean;
};

function ImageItem({ img, onSelect, onDelete, priority }: ImageItemProps) {
  const { ref, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
  });

  return (
    <div ref={ref} className={styles.thumbWrapper}>
      {inView ? (
        <Image
          src={API + img.url}
          alt={img.filename}
          width={100}
          height={100}
          className={styles.thumb}
          onClick={() => onSelect(img)}
          unoptimized
          priority={priority} // <-- передаем сюда
          style={{ width: '100%', height: 'auto', objectFit: 'cover', borderRadius: '8px' }}
        />
      ) : (
        <div 
          style={{ 
            width: '100%', 
            paddingTop: '100%',
            background: '#eee',
            borderRadius: '8px'
          }}
        />
      )}
      <button
        className={styles.deleteBtn}
        onClick={() => onDelete(img.filename)}
        aria-label="Удалить"
      >
        <MdClose size={16} />
      </button>
    </div>
  );
}
