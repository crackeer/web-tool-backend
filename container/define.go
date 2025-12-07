package container

type Description struct {
	Title     string      `json:"title"`
	InputForm interface{} `json:"input_form"`
}

type Tool interface {
	GetDescription() *Description
	RecvInput(bytes []byte) error
	Run(func(string)) error
}

var (
	ToolMap map[string]Tool
)

func init() {
	ToolMap = make(map[string]Tool)
}

func RegisterTool(name string, tool Tool) {
	ToolMap[name] = tool
}

func GetTool(name string) Tool {
	return ToolMap[name]
}

func GetToolConfig() []map[string]interface{} {
	names := make([]map[string]interface{}, 0, len(ToolMap))
	for name, tool := range ToolMap {
		desc := tool.GetDescription()
		names = append(names, map[string]interface{}{
			"name":       name,
			"title":      desc.Title,
			"input_form": desc.InputForm,
		})
	}
	return names
}
