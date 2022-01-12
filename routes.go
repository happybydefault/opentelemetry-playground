package playground

func (p *Playground) registerRoutes() {
	p.mux.Get("/", p.withInstrumentation("hello world handler", p.helloWorldHandler))
}
