import { renderToString } from '@vue/server-renderer'
import { createApp } from '~/main'

export async function render(url: string) {
  const { app, router } = createApp({})
  await router.push(url)

  const ctx: any = {}

  return await renderToString(app, ctx)
}

function ssrRender(url: string) {
  return render(url).then((html) => {
    return html
  })
}

(globalThis as any).ssrRender = ssrRender
