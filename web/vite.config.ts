import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import path from 'path'

const listenHost = process.env['VITE_DEV_SERVER_LISTEN_HOST'] || '0.0.0.0'
const listenPort = parseInt(process.env['VITE_DEV_SERVER_LISTEN_PORT'] || '') || 8080

const externalHost = process.env['VITE_DEV_SERVER_EXTERNAL_HOST'] || listenHost
const externalPort = parseInt(process.env['VITE_DEV_SERVER_EXTERNAL_PORT'] || '') || listenPort

const backendURL = process.env['VITE_DEV_SERVER_BACKEND_URL'] || 'http://localhost:8081'

export default defineConfig({
  plugins: [vue()],
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
