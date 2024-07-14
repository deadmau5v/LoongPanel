package v1

type route struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Comment string `json:"comment"`
}

var RouteList []route

func AddRouteComment(method, path, comment string) {
	RouteList = append(RouteList, route{Path: path, Method: method, Comment: comment})
}

func GetRoutes() []route {
	return RouteList
}

func GetRouteComment(method, path string) string {
	for _, route := range RouteList {
		if route.Method == method && route.Path == path {
			return route.Comment
		}
	}
	return ""
}
