package testkit

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	// modules package is shared, so we can register test packages as well.
	modules.Register(make(Module))
	modules.Register(make(GenerativeModule))
	modules.Register(make(RerankerModule))
}

const ModuleName = "testkit-module"

// Module implements [internal.Module] for a map of arbitrary values.
// "testkit-module" is registered in the modules.registry.
type Module map[string]any

func (Module) Name() string { return ModuleName }

const GenerativeModuleName = "generative-testkit"

// GenerativeModule is a "generative" variant of [Module].
type GenerativeModule Module

func (GenerativeModule) Name() string { return GenerativeModuleName }

const RerankerModuleName = "reranker-testkit"

// RerankerModule is a "reranker" variant of [Module].
type RerankerModule Module

func (RerankerModule) Name() string { return RerankerModuleName }
