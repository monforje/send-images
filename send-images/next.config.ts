const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '9999',
        pathname: '/uploads/**',
      },
    ],
  },
};

module.exports = nextConfig;
