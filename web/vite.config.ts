import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import path from 'path'

const listenHost = process.env['VITE_DEV_SERVER_LISTEN_HOST'] || '0.0.0.0'
const listenPort = parseInt(process.env['VITE_DEV_SERVER_LISTEN_PORT'] || '') || 8080

export default defineConfig({
  plugins: [vue()],
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
