package main

import (
	"os"
	"strings"
)

const (
	PathSep                     = string(os.PathSeparator)
	DriveRoot                   = "root"
	DefaultLogLevel             = "info"
	DefaultGDIgnorePath         = "gdignore.yaml"
	DefaultDownloadDelay        = "5-20"
	DefaultMaxParallelDownloads = 5
)

func joinToPath(tokens ...string) string {
	var items []string
	for _, item := range tokens {
		items = append(items, item)
	}
	return strings.Join(items, PathSep)
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if value == item {
			return true
		}
	}
	return false
}

func containsWithWildcard(items []string, value string) bool {
	for _, item := range items {
		idx := strings.Index(item, "*")
		if idx > -1 && strings.HasPrefix(value, item[:idx]) && strings.HasSuffix(value, item[idx+1:]) {
			return true
		}
	}
	return false
}

func appendToMap(itemMap map[string][]string, key string, item string) {
	val, ok := itemMap[key]
	if ok {
		itemMap[key] = append(val, item)
	} else {
		itemMap[key] = []string{item}
	}
}
