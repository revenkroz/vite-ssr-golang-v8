import { defineConfig } from 'vite'
import config from './vite.config'

export default defineConfig({
  ...config,
  build: {
    rollupOptions: {
      input: {
        server: 'src/entry-server.ts',
      },
      output: {
        format: 'cjs',
        entryFileNames: '[name].js',
        inlineDynamicImports: true,
      },
    },
  },
})
