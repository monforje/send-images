import { dirname } from 'path';
import { fileURLToPath } from 'url';
import { FlatCompat } from '@eslint/eslintrc';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const eslintConfig = [
  // Подключаем конфиги от Next.js и TypeScript
  ...compat.extends('next/core-web-vitals', 'next/typescript'),

  // Базовые настройки
  {
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      parserOptions: {
        project: './tsconfig.json',
        sourceType: 'module',
      },
    },
    rules: {
      // ❗ Строгая проверка неиспользуемых переменных, с игнором для `_`
      '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],

      // 🔁 Проверка на зависимости хуков (включена по умолчанию в Next, но на всякий случай)
      'react-hooks/exhaustive-deps': 'warn',

      // 👌 Предупреждение при использовании `any`
      '@typescript-eslint/no-explicit-any': 'warn',

      // 🧼 Стиль: запрещаем `console.log`, но разрешаем `console.error` и `warn`
      'no-console': ['warn', { allow: ['warn', 'error'] }],

      // 🎯 Ошибка при пропущенном `key` в списках
      'react/jsx-key': 'error',
    },
  },

  // Опционально: настройки для конфигов, если используешь *.config.ts
  {
    files: ['*.config.ts', '*.config.js', '*.mjs'],
    rules: {
      '@typescript-eslint/no-var-requires': 'off',
    },
  },
];

export default eslintConfig;
