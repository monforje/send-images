import { getImages, type Image } from '@/lib/api';
import ClientApp from '@/components/ClientApp';

export default async function Page() {
  let images: Image[] = [];

  try {
    images = await getImages();
  } catch (error) {
    console.error('Ошибка при получении изображений:', error);
  }

  return <ClientApp initialImages={images} />;
}
