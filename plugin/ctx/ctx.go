package ctx

import "context"

func genGlobalName(pluginName string, varName string) string {
	return "plugin:" + pluginName + ":" + varName
}

func WithPluginVar(ctx context.Context, pluginName string, varName string, value any) context.Context {
	globalVarName := genGlobalName(pluginName, varName)

	return context.WithValue(ctx, globalVarName, value)
}

func GetPluginVar(ctx context.Context, pluginName string, varName string) any {
	globalVarName := genGlobalName(pluginName, varName)

	return ctx.Value(globalVarName)
}

func GetPluginVarString(ctx context.Context, pluginName string, varName string) string {
	val := GetPluginVar(ctx, pluginName, varName)
	if val == nil {
		return ""
	}
	return val.(string)
}

func GetPluginVarInt(ctx context.Context, pluginName string, varName string) int {
	val := GetPluginVar(ctx, pluginName, varName)
	if val == nil {
		return 0
	}
	return val.(int)
}

var globalVars map[string]string

func init() {
	globalVars = make(map[string]string)
}

func AddGlobalVar(key, value string) {
	if _, ok := globalVars[key]; !ok {
		globalVars[key] = value
	}
}

func GetGlobalVar(key string) string {
	return globalVars[key]
}
