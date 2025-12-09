package container

import (
	"os"
	"path/filepath"
	"strings"
	"web-tool-backend/util"
)

type Crd struct {
	Title   string      `json:"title"`
	Name    string      `json:"name"`
	Form    interface{} `json:"form"`
	Package string      `json:"package"`
}

// GetCrdList
//
//	@return []Crd
func GetCrdList() []Crd {
	retData := []Crd{}
	entries, err := os.ReadDir(cfg.CrdDir)
	if err != nil {
		return retData
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		entryPath := filepath.Join(cfg.CrdDir, entry.Name())
		crd := Crd{}
		err = util.ReadJsonFile(entryPath, &crd)
		if err != nil {
			continue
		}
		crd.Name = strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		retData = append(retData, crd)
	}
	return retData
}
