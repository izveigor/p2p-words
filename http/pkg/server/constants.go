package server

import (
	"os"
	"strings"
)

func GetPathToStatic() string {
	if !strings.HasSuffix(os.Args[0], ".test") {
		return "pkg/server/static/"
	} else {
		return "./static/"
	}
}

var PATH_TO_STATIC string = GetPathToStatic()
