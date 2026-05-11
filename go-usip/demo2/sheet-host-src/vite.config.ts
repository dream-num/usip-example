import process from 'node:process'
import { defineConfig, loadEnv } from 'vite'

import packageJson from './package.json'

export default ({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const univerEndpoint = env.UNIVER_ENDPOINT || 'http://localhost:8000'
  return defineConfig({
    server: {
      cors: true,
      proxy: {
        '/universer-api': {
          target: univerEndpoint,
          changeOrigin: true,
          ws: true,
        },
      },
      allowedHosts: ['local.univer.plus'],
    },
    define: {
      'process.env.UNIVER_CLIENT_LICENSE': `"${env.UNIVER_CLIENT_LICENSE}"` || '"%%UNIVER_CLIENT_LICENSE_PLACEHOLDER%%"',
      'process.env.UNIVER_VERSION': `"${packageJson.dependencies['@univerjs/presets']}"`,
    },
    base: '/sheet/',
    worker: {
      format: 'es',
      rollupOptions: {
        output: {
          entryFileNames: 'worker.js',
        },
      },
    },
    build: {
      outDir: '../web/public/sheet-host',
      emptyOutDir: true,
      rollupOptions: {
        output: {
          entryFileNames: 'main.js',
        },
      },
    },
  })
}
