import { makeApp } from '~/main'

const { app, router } = makeApp()
router.push(window.location.pathname)

app.mount('#app', true)
