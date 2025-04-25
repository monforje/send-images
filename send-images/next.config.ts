/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: process.env.NEXT_PUBLIC_API_PROTOCOL ?? 'http',
        hostname: process.env.NEXT_PUBLIC_API_HOST ?? 'localhost',
        port: process.env.NEXT_PUBLIC_API_PORT ?? '9999',
        pathname: '/uploads/**',
      },
    ],
  },

  reactStrictMode: true,

  experimental: {
    // serverActions: true, // если ты пользуешься Server Actions
  },

  typescript: {
    ignoreBuildErrors: false,
  },

  eslint: {
    // Чтобы билд не падал при ошибке ESLint (можно true → false на проде)
    ignoreDuringBuilds: true,
  },

  // optional: basePath: '/your-subpath',
};

export default nextConfig;
