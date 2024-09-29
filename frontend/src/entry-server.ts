import { renderToString } from '@vue/server-renderer'
import { makeApp } from '~/main'

export async function render(url: string) {
  const { app, router } = makeApp()
  await router.push(url)

  const ctx: any = {}

  return await renderToString(app, ctx)
}

async function ssrRender(url: string) {
  return await render(url)
}

(globalThis as any).ssrRender = ssrRender
