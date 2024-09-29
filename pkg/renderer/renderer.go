package renderer

import (
	"fmt"
	"rogchap.com/v8go"
)

// Renderer renders a React application to HTML.
type Renderer struct {
	pool          *IsolatePool
	ssrScriptName string
}

// NewRenderer creates a new server side renderer for a given script.
func NewRenderer(scriptContents string) *Renderer {
	ssrScriptName := "server.js"

	return &Renderer{
		pool:          NewIsolatePool(scriptContents, ssrScriptName),
		ssrScriptName: ssrScriptName,
	}
}

// Render renders the provided path to HTML.
func (r *Renderer) Render(urlPath string) (string, error) {
	iso := r.pool.Get()
	defer r.pool.Put(iso)

	ctx := v8go.NewContext(iso.Isolate)
	defer ctx.Close()

	iso.RenderScript.Run(ctx)

	renderCmd := fmt.Sprintf(`ssrRender("%s")`, urlPath)
	val, err := ctx.RunScript(renderCmd, r.ssrScriptName)
	if err != nil {
		return "", formatError(err)
	}

	renderedHtml := ""

	if val.IsPromise() {
		result, err := resolvePromise(ctx, val, err)
		if err != nil {
			return "", formatError(err)
		}

		renderedHtml = result.String()
	} else {
		renderedHtml = val.String()
	}

	return renderedHtml, nil
}
