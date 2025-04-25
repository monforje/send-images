// send-images/src/lib/api.ts

// Базовый URL для общения с Go-бэкендом
const API = process.env.NEXT_PUBLIC_API_URL;

// Получение списка изображений с сервера
export async function getImages(): Promise<{ filename: string; url: string }[]> {
  try {
    const res = await fetch(`${API}/images`); // GET-запрос на эндпоинт /images
    const data = await res.json(); // Парсим JSON
    return data.images || []; // Возвращаем массив изображений или пустой массив
  } catch {
    return []; // В случае ошибки — возвращаем пустой список
  }
}

// Загрузка изображения на сервер
export async function uploadImage(file: File): Promise<boolean> {
  const formData = new FormData(); // Создаём тело запроса
  formData.append('file', file); // Добавляем файл под ключом "file"

  try {
    const res = await fetch(`${API}/upload`, {
      method: 'POST',
      body: formData, // Отправляем как multipart/form-data
    });
    return res.ok; // true, если загрузка успешна (HTTP 200)
  } catch {
    return false; // В случае ошибки — false
  }
}

// Удаление изображения по имени файла
export async function deleteImage(filename: string): Promise<boolean> {
  try {
    const res = await fetch(`${API}/images?filename=${encodeURIComponent(filename)}`, {
      method: 'DELETE', // DELETE-запрос с query-параметром filename
    });
    return res.status === 204; // true, если файл удалён успешно (No Content)
  } catch {
    return false; // Ошибка — возвращаем false
  }
}
