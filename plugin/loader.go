package plugin

import (
	"github.com/wklken/flow/plugin/basic_auth"
	"github.com/wklken/flow/plugin/request_id"
)

var Plugins map[string]Plugin

func init() {
	Plugins = make(map[string]Plugin)
	// FIXME: how to register the plugin
	Plugins["request_id"] = &request_id.Plugin{}
	Plugins["basic_auth"] = &basic_auth.Plugin{}
}
