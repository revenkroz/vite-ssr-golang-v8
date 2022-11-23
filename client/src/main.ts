import { createSSRApp } from 'vue'
import { createMemoryHistory, createRouter, createWebHistory } from 'vue-router'
import routes from 'virtual:generated-pages'
import App from './App.vue'

import '@unocss/reset/tailwind.css'
import './styles/main.css'
import 'uno.css'

const isServer = typeof window === 'undefined'
export function createApp(context: any = {}) {
  // eslint-disable-next-line no-console
  console.log(context)

  const app = createSSRApp(App)
  const router = createRouter({
    history: isServer ? createMemoryHistory() : createWebHistory(import.meta.env.BASE_URL),
    routes,
  })
  app.use(router)

  return { app, router }
}
