import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import { defineConfig } from 'vite'

const listenHost = process.env['VITE_DEV_SERVER_LISTEN_HOST'] || '0.0.0.0'
const listenPort = parseInt(process.env['VITE_DEV_SERVER_LISTEN_PORT'] || '') || 8080

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    host: listenHost,
    port: listenPort,
    strictPort: true,
  },
  resolve: {
    alias: {
      '~': path.resolve(__dirname, './'),
    },
  },
})
