import { createApp } from '~/main'

const { app, router } = createApp({})
router.push(window.location.pathname)

router.isReady().then(() => {
  app.mount('#app', true)
})
