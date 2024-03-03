package ctx

import "context"

func genGlobalName(pluginName string, varName string) string {
	return "plugin:" + pluginName + ":" + varName
}

func WithPluginVar(ctx context.Context, pluginName string, varName string, value string) context.Context {
	globalVarName := genGlobalName(pluginName, varName)

	return context.WithValue(ctx, globalVarName, value)
}

func GetPluginVar(ctx context.Context, pluginName string, varName string) any {
	globalVarName := genGlobalName(pluginName, varName)

	return ctx.Value(globalVarName)
}

func GetPluginVarString(ctx context.Context, pluginName string, varName string) string {
	return GetPluginVar(ctx, pluginName, varName).(string)
}

func GetPluginVarInt(ctx context.Context, pluginName string, varName string) int {
	return GetPluginVar(ctx, pluginName, varName).(int)
}
