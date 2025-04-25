import styles from '@/styles/home.module.css';
import Image from 'next/image';
import { MdClose } from 'react-icons/md';

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
  const handleDelete = (filename: string) => {
    if (confirm(`Удалить ${filename}?`)) {
      onDelete(filename);
    }
  };

  return (
    <div className={styles.galleryWrapper}>
      <h2 className={styles.galleryTitle}>Мои картинки</h2>

      {images.length === 0 ? (
        <p className={styles.empty}>Здесь пока пусто 😢</p>
      ) : (
        <div className={styles.thumbs}>
          {images.map((img) => (
            <div key={img.filename} className={styles.thumbWrapper}>
              <Image
                src={API + img.url}
                alt={img.filename}
                width={100}
                height={100}
                className={styles.thumb}
                onClick={() => onSelect(img)}
                unoptimized
                // lazy loading будет по умолчанию включен
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
      )}
    </div>
  );
}
