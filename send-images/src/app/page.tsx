import Header from '@/components/Header';
import ClientApp from '@/components/ClientApp';
import { getImages, type Image } from '@/lib/api';

export default async function Page() {
  let images: Image[] = [];

  try {
    images = await getImages();
  } catch (error) {
    console.error('Ошибка при получении изображений:', error);
  }

  return (
    <>
      <Header />
      <ClientApp initialImages={images} />
    </>
  );
}
