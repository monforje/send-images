'use client';

import styles from '@/styles/home.module.css';

export default function Header() {
  return (
    <header className={styles.header}>
      <div className={styles.headerContent}>
        <div className={styles.avatar}></div>
      </div>
    </header>
  );
}
