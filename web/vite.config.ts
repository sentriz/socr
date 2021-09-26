import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { defineConfig } from 'vite'
import path from 'path'

const listenHost = process.env['VITE_DEV_SERVER_LISTEN_HOST'] || '0.0.0.0'
const listenPort = parseInt(process.env['VITE_DEV_SERVER_LISTEN_PORT'] || '') || 8080

const externalHost = process.env['VITE_DEV_SERVER_EXTERNAL_HOST'] || listenHost
const externalPort = parseInt(process.env['VITE_DEV_SERVER_EXTERNAL_PORT'] || '') || listenPort

const backendURL = process.env['VITE_DEV_SERVER_BACKEND_URL'] || 'http://localhost:8081'

const pluginVue = vue()
const pluginPWA = VitePWA({
  includeAssets: [
    'apple-icon-180.png',
    'apple-splash-1125-2436.jpg',
    'apple-splash-1136-640.jpg',
    'apple-splash-1170-2532.jpg',
    'apple-splash-1242-2208.jpg',
    'apple-splash-1242-2688.jpg',
    'apple-splash-1284-2778.jpg',
    'apple-splash-1334-750.jpg',
    'apple-splash-1536-2048.jpg',
    'apple-splash-1620-2160.jpg',
    'apple-splash-1668-2224.jpg',
    'apple-splash-1668-2388.jpg',
    'apple-splash-1792-828.jpg',
    'apple-splash-2048-1536.jpg',
    'apple-splash-2048-2732.jpg',
    'apple-splash-2160-1620.jpg',
    'apple-splash-2208-1242.jpg',
    'apple-splash-2224-1668.jpg',
    'apple-splash-2388-1668.jpg',
    'apple-splash-2436-1125.jpg',
    'apple-splash-2532-1170.jpg',
    'apple-splash-2688-1242.jpg',
    'apple-splash-2732-2048.jpg',
    'apple-splash-2778-1284.jpg',
    'apple-splash-640-1136.jpg',
    'apple-splash-750-1334.jpg',
    'apple-splash-828-1792.jpg',
    'favicon.ico',
    'favicon.svg',
    'manifest-icon-192.png',
    'manifest-icon-512.png',
  ],
  manifest: {
    name: 'socr',
    short_name: 'socr',
    description: 'screenshot ocr server',
    theme_color: '#90cdf4',
    icons: [
      {
        src: 'manifest-icon-192.png',
        sizes: '192x192',
        type: 'image/png',
        purpose: 'maskable any',
      },
      {
        src: 'manifest-icon-512.png',
        sizes: '512x512',
        type: 'image/png',
        purpose: 'maskable any',
      },
    ],
    // TODO: figure out if share_target is working or not
    // @ts-ignore
    share_target: {
      action: '/api/upload',
      method: 'POST',
      enctype: 'multipart/form-data',
      params: {
        files: [
          {
            name: 'i',
            accept: ['image/*'],
          },
        ],
      },
    },
  },
})

export default defineConfig({
  plugins: [pluginVue, pluginPWA],
  server: {
    host: listenHost,
    port: listenPort,
    strictPort: true,
    proxy: {
      '/api': backendURL,
    },
    hmr: {
      host: externalHost,
      port: externalPort,
    },
  },
  resolve: {
    alias: {
      '~': path.resolve(__dirname, './'),
    },
  },
})
