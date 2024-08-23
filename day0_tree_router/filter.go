package __tree_router

type FilterBuilder func(next Filter) Filter

type Filter func(c *Context)
