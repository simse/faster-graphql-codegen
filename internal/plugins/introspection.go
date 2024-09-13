package plugins

func (p* PluginTask) Introspect() {
    p.Output.WriteString("{\"message\":\"introspection plugin is not implemented\"}\n")
}