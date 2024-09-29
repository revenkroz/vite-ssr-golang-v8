# Golang SSR for Vite Apps Example

Minimal example of using V8 to render a Frontend App built with Vite.

## Used Packages

- [v8go](https://github.com/rogchap/v8go) - Go bindings for V8 JavaScript engine
- [Vitesse Lite](https://github.com/antfu/vitesse-lite) - Frontend App Template written in Vue and built with Vite

## Covered features

- [x] Fully functional frontend application (not just a simple "Hello World")
- [x] SSR for the frontend application
- [x] V8 Isolate pool to avoid creating a new Isolate for each request
- [x] Embedded frontend application in the binary to reduce the file system calls

## How to run the example

```bash
docker compose up -d
curl http://localhost:8080/hi/test
```

Note: if there is some issues with building the image, remove `--platform=linux/amd64` options from Dockerfile. This was used to avoid issues when running on Apple M1 architecture.

## Pros and Cons

Pros:
- No need to use Node.js
- All in one binary

Cons:
- The more big js-bundle size the slower rendering
- You can't use Vite features like hot module replacement
- To build the server you need a lot of c-libs installed


## Keynotes

### Separate config for Vite Build

The default config is configured in `vite.config.ts` and the additional configuration for the SSR build is located in [`vite.config.prod.ts`](client/vite.config.prod.ts).
SSR config builds the frontend application with target `cjs` as it is required for the V8 engine.

### Entry point for SSR and Client are different

See:
- Client-side entry: [`client/src/entry-client.ts`](client/src/entry-client.ts)
- Server-side entry: [`client/src/entry-server.ts`](client/src/entry-server.ts)

### How to speed up SSR more

Split frontend build for server for multiple files (currently it's a single file).
A `require` function must be implemented in the V8 context to load the files.
It's better to cache the required files in memory.

### Hot Reloading

It's not possible to use hot reloading with V8. For frontend development it's better to use Vite directly and store code it in another repo.
For backend development use any watcher to rebuild all (e.g. [air](https://github.com/cosmtrek/air)).

### Build without Vite and with Go only

Some cool features of Vite will be missing (e.g. Glob Import, Dynamic Import, hot module replacement, etc.),
but it's possible to build the frontend application with [ESBuild](https://github.com/evanw/esbuild) - a Go-based bundler. It has a Go API and it's very fast.
Actually, it's used by Vite under the hood.

## Approaches to render the frontend application
### First Approach

Fastest one. The idea is to run SSR script, get function with sensitive args and run it.

#### Go Code

`iso` object has `Isolate *v8go.Isolate` and `RenderScript *v8go.UnboundScript` (it was compiled before).

```go
// renderer.go
// ...

func (r *Renderer) Render(urlPath string) (string, error) {
	iso := r.pool.Get()
	defer r.pool.Put(iso)

	ctx := v8go.NewContext(iso.Isolate)
	defer ctx.Close()

	iso.RenderScript.Run(ctx)

	renderCmd := fmt.Sprintf(`ssrRender("%s")`, urlPath)
	val, err := ctx.RunScript(renderCmd, r.ssrScriptName)
	if err != nil {
		if jsErr, ok := err.(*v8go.JSError); ok {
			err = fmt.Errorf("%v", jsErr.StackTrace)
		}
		return "", nil
	}

	return val.String(), nil
}
```

#### JS Code

```typescript
// entry-server.ts
// ...

function ssrRender(url: string) {
  return render(url).then((html) => {
    return html
  })
}

(globalThis as any).ssrRender = ssrRender
```


### Second Approach

Create global object in Go with `render` function that recieves rendered html, concat string in Go and return.

#### Go Code

`iso` object has `Isolate *v8go.Isolate` and `RenderScript *v8go.UnboundScript` (it was compiled before).

```go
// renderer.go
// ...

func (r *Renderer) Render(urlPath string) (string, error) {
	iso := r.pool.Get()
	defer r.pool.Put(iso)

	outputHTML := ""
	
	ssrObject := v8go.NewObjectTemplate(iso.Isolate)
	ssrObject.Set("href", urlPath)
	ssrObject.Set("render", v8go.NewFunctionTemplate(iso.Isolate, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			outputHTML = args[0].String()
		}
		return nil
	}))
	
	globalObject := v8go.NewObjectTemplate(iso.Isolate)
	globalObject.Set("ssr", ssrObject)
	
	ctx := v8go.NewContext(iso.Isolate, globalObject)
	defer ctx.Close()
	
	start := time.Now()
	iso.RenderScript.Run(ctx)
	//if _, err := ctx.RunScript(r.scriptSource, r.Path); err != nil {
	//	return "", err
	//}
	fmt.Println("Script run:", time.Since(start))
	
	return outputHTML, nil
}
```

#### JS Code

```typescript
// entry-server.ts
// ...

if (typeof ssr !== 'undefined') {
  render(ssr.href).then((html) => {
    ssr.render(html)
  })
}
```
