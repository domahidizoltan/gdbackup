package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	PathSep                     = string(os.PathSeparator)
	DriveRoot                   = "root"
	DefaultLogLevel             = "info"
	DefaultGDIgnorePath         = "gdignore.yaml"
	DefaultDownloadDelay        = "3-8"
	DefaultMaxParallelDownloads = 3
	DefaultLogFile              = "gdbackup.log"
	DefaultBackupPath           = ""
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

func makeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Errorf("Could not create folder %s: %v", dir, err)
		}
	}
}

func getValidDir(path string) string {
	if path == "" {
		workdir, err := os.Getwd()
		if err != nil {
			log.Warn("Could not get workdir: %v", err)
			workdir, _ = os.UserHomeDir()
		}
		path = workdir
	}
	return path
}