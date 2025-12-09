package container

type ToolFunc func(string, func(string)) (string, error)

var (
	ToolMap map[string]ToolFunc = make(map[string]ToolFunc)
)

func RegisterTool(name string, tool ToolFunc) {
	ToolMap[name] = tool
}

func GetTool(name string) ToolFunc {
	return ToolMap[name]
}
