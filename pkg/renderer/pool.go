package renderer

import (
	"rogchap.com/v8go"
	"runtime"
	"sync"
)

type IsolatePool struct {
	pool *sync.Pool
}

type IsolateContainer struct {
	Isolate      *v8go.Isolate
	RenderScript *v8go.UnboundScript
}

func NewIsolatePool(ssrScriptContents string, ssrScriptName string) *IsolatePool {
	return &IsolatePool{
		pool: newIsolatePool(ssrScriptContents, ssrScriptName),
	}
}

func newIsolatePool(ssrScriptContents string, ssrScriptName string) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			isolate := v8go.NewIsolate()
			script, _ := isolate.CompileUnboundScript(ssrScriptContents, ssrScriptName, v8go.CompileOptions{})

			runtime.SetFinalizer(isolate, func(iso *v8go.Isolate) {
				if iso != nil {
					iso.Dispose()
				}
			})

			return &IsolateContainer{
				Isolate:      isolate,
				RenderScript: script,
			}
		},
	}
}

func (p *IsolatePool) Get() *IsolateContainer {
	item := p.pool.Get()
	isolate, ok := item.(*IsolateContainer)
	if !ok {
		panic(item.(error))
	}

	return isolate
}

func (p *IsolatePool) Put(isolateContainer *IsolateContainer) {
	p.pool.Put(isolateContainer)
}
