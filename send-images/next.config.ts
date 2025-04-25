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
    // serverActions: true, // если ты используешь Server Actions
  },

  eslint: {
    ignoreDuringBuilds: true,
  },

  // Убираем productionSourceMaps, так как это больше не поддерживается
};

export default nextConfig;
