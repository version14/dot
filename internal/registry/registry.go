package registry

type Registry struct {
	Name        string
	Description string
	Plugins     []Plugin
}

func (registry *Registry) RegisterPlugin(plugin Plugin) {
	registry.Plugins = append(registry.Plugins, plugin)
}
