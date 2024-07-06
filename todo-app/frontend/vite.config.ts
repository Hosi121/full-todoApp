import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0', // 外部アクセスを許可
    port: 3000,      // デフォルトポートを変更
    watch: {
      usePolling: true,
      interval: 1000
    }
  }
});