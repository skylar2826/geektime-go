package __filter_builder

type FilterBuilder func(next Filter) Filter

type Filter func(c *Context)

var filterBuilders = make(map[string]FilterBuilder, 4)

func RegisterFilterBuilder(name string, builder FilterBuilder) {
	filterBuilders[name] = builder
}

func GetFilterBuilder(name string) FilterBuilder {
	return filterBuilders[name]
}
