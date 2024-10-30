package golang

type RouteItem struct {
	Method string
	Import string
	Path   string
}

type RouteGroup struct {
	Name  string
	Items []RouteItem
}
