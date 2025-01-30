package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type UrlParam struct {
	Url      string
	Position int
}

func GetUrlParamId(path UrlParam) (int, error) {
	pathParts := strings.Split(path.Url, "/")
	if len(pathParts) != path.Position+1 {
		return 0, fmt.Errorf("wrong path")
	}

	id := pathParts[path.Position]

	return strconv.Atoi(id)
}
