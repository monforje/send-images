const API = process.env.NEXT_PUBLIC_API_URL;

export type Image = {
  filename: string;
  url: string;
};

export async function getImages(): Promise<Image[]> {
  try {
    const res = await fetch(`${API}/images`);
    const data = await res.json();
    return data.images || [];
  } catch (error) {
    console.error('Ошибка при загрузке изображений:', error);
    return [];
  }
}

export async function uploadImage(file: File): Promise<Image | null> {
  const formData = new FormData();
  formData.append('file', file);

  try {
    const res = await fetch(`${API}/upload`, {
      method: 'POST',
      body: formData,
    });

    if (!res.ok) throw new Error('Ошибка при загрузке');

    const data = await res.json();
    return { filename: data.filename, url: data.url };
  } catch (error) {
    console.error('Ошибка загрузки:', error);
    return null;
  }
}

export async function deleteImage(filename: string): Promise<boolean> {
  try {
    const res = await fetch(`${API}/images?filename=${encodeURIComponent(filename)}`, {
      method: 'DELETE',
    });

    return res.ok;
  } catch (error) {
    console.error('Ошибка удаления:', error);
    return false;
  }
}
