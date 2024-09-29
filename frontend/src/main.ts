import { createApp, createSSRApp } from 'vue'
import { createMemoryHistory } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router/auto'
import App from './App.vue'

import '@unocss/reset/tailwind.css'
import './styles/main.css'
import 'uno.css'

const isServer = typeof window === 'undefined'

export function makeApp() {
  const app = isServer ? createSSRApp(App) : createApp(App)
  const router = createRouter({
    history: isServer ? createMemoryHistory() : createWebHistory(import.meta.env.BASE_URL),
  })

  app.use(router)

  return {
    app,
    router,
  }
}
