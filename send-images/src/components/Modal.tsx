'use client';

import { useEffect, useState } from 'react';
import styles from '@/styles/home.module.css';
import Image from 'next/image';
import { MdNavigateBefore, MdNavigateNext } from 'react-icons/md';

type ImageType = { filename: string; url: string };

type ModalProps = {
  selected: ImageType;
  images: ImageType[];
  setSelected: (img: ImageType | null) => void;
  setIsModalOpen: (isOpen: boolean) => void;
};

export default function Modal({ selected, images, setSelected, setIsModalOpen }: ModalProps) {
  const [visible, setVisible] = useState(true);
  const currentIndex = images.findIndex((img) => img.filename === selected.filename);

  useEffect(() => {
    setIsModalOpen(true);
    return () => setIsModalOpen(false);
  }, [setIsModalOpen]);

  const showPrev = () => {
    const prevIndex = currentIndex === 0 ? images.length - 1 : currentIndex - 1;
    setSelected(images[prevIndex]);
  };

  const showNext = () => {
    const nextIndex = currentIndex === images.length - 1 ? 0 : currentIndex + 1;
    setSelected(images[nextIndex]);
  };

  const closeModal = () => {
    setVisible(false);
    setTimeout(() => setSelected(null), 200); // Совпадает с анимацией fadeOut
  };

  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') closeModal();
      if (e.key === 'ArrowLeft') showPrev();
      if (e.key === 'ArrowRight') showNext();
    };
    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  }, [currentIndex]);

  return (
    <div
      className={`${styles.modalOverlay} ${!visible ? styles.modalFadeOut : ''}`}
      onClick={closeModal}
      role="dialog"
      aria-modal="true"
    >
      {/* Контент модалки тоже плавно исчезает */}
      <div
        className={`${styles.modalContent} ${!visible ? styles.modalContentFadeOut : ''}`}
        onClick={(e) => e.stopPropagation()}
      >
        <Image
          src={process.env.NEXT_PUBLIC_API_URL + selected.url}
          alt={selected.filename}
          width={1200}
          height={800}
          priority
          unoptimized
          style={{ height: 'auto' }}
          className={styles.modalImage}
        />

        <button
          className={styles.navLeft}
          onClick={(e) => { e.stopPropagation(); showPrev(); }}
        >
          <MdNavigateBefore size={28} />
        </button>

        <button
          className={styles.navRight}
          onClick={(e) => { e.stopPropagation(); showNext(); }}
        >
          <MdNavigateNext size={28} />
        </button>

        <div className={styles.modalCounter}>
          {currentIndex + 1} / {images.length}
        </div>
      </div>
    </div>
  );
}
