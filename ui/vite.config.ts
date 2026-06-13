import tailwindcss from "@tailwindcss/vite";
import react from '@vitejs/plugin-react';
import path from 'path';
import mantineTheme from "tailwind-preset-mantine/vite";
import { defineConfig, loadEnv } from 'vite';
import { imagetools } from 'vite-imagetools';
import checker from 'vite-plugin-checker';

export default defineConfig(({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return {
    build: {
      outDir: '../public',
      emptyOutDir: true,
      sourcemap: true,
    },
    plugins: [
      tailwindcss(),
      react(),
      checker({
        // e.g. use TypeScript check
        typescript: true,
      }),
      imagetools(),
      mantineTheme({
        input: "./src/lib/theme/theme.ts",
      }),
    ],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      port: 3000,
      proxy: {
        '/api': {
          target: 'http://backend:3001',
          changeOrigin: true,
        }
      }
    }
  }
})
