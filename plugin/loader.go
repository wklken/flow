package plugin

import (
	"sort"

	"github.com/justinas/alice"
	"github.com/wklken/flow/plugin/basic_auth"
	"github.com/wklken/flow/plugin/file_logger"
	"github.com/wklken/flow/plugin/request_id"
)

func New(name string) Plugin {
	switch name {
	case "request_id":
		return &request_id.Plugin{}
	case "basic_auth":
		return &basic_auth.Plugin{}
	case "file_logger":
		return &file_logger.Plugin{}
	}
	return nil
}

func BuildPluginChain(plugins ...Plugin) alice.Chain {
	// sort the plugin by priority
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Priority() < plugins[j].Priority()
	})

	// build the alice chain
	chain := alice.New()
	for _, plugin := range plugins {
		chain = chain.Append(plugin.Handler)
	}

	return chain
}
