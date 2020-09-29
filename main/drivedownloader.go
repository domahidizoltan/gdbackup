package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
)

type ExportItemType string

const (
	Docx ExportItemType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	Xlsx ExportItemType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	Pptx ExportItemType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
)

type DriveDownloader struct {
	ignoreItems              map[string][]string
	ignoreItemsRootLevelKeys []string
	backupDir                string
	backupDirName            string
	client                   DriveClient
	downloadQueue            chan downloadItem
	semaphore                chan struct{}
	inProgress               sync.WaitGroup
	DownloadStats            DownloadStats
	downloadConfig           downloadConfig
}

type downloadItem struct {
	item     *drive.File
	savePath string
}

type DownloadStats struct {
	totalFileCount int
	totalSize      int64
}

type downloadConfig struct {
	minDelay             int
	maxDelay             int
	maxParallelDownloads int
	backupPath           string
}

func NewDownloader(client DriveClient, ignoreItems map[string][]string) *DriveDownloader {
	var rootKeys []string
	for key, _ := range ignoreItems {
		rootKeys = append(rootKeys, key)
	}

	folder := "gdbackup" + time.Now().Format("20060102")

	var waitgroup sync.WaitGroup

	downloader := &DriveDownloader{
		ignoreItems:              ignoreItems,
		ignoreItemsRootLevelKeys: rootKeys,
		backupDir:                joinToPath(DefaultBackupPath, folder),
		backupDirName:            folder,
		client:                   client,
		downloadQueue:            make(chan downloadItem),
		semaphore:                make(chan struct{}, DefaultMaxParallelDownloads),
		inProgress:               waitgroup,
		DownloadStats:            DownloadStats{0, 0},
		downloadConfig:           downloadConfig{},
	}

	downloader.ConfigureDownloader(DefaultDownloadDelay, DefaultMaxParallelDownloads, DefaultBackupPath)
	rand.Seed(time.Now().UnixNano())
	go downloader.init()

	return downloader
}

func (this *DriveDownloader) init() {
	for job := range this.downloadQueue {
		log.Tracef("received download job %v", job)
		go func(job downloadItem) {
			this.semaphore <- struct{}{}
			this.fetchSingleItem(job.item, job.savePath)
			<-this.semaphore
		}(job)
	}
}

func (this *DriveDownloader) ConfigureDownloader(delay string, maxParallelDownloads int, backupPath string) {
	if strings.Count(delay, "-") == 1 {
		delays := strings.Split(delay, "-")

		minDelay, err := strconv.Atoi(delays[0])
		if err != nil {
			log.Warnf("Unable to parse min delay %s. Falling back to default %s", minDelay, this.downloadConfig.minDelay)
		} else {
			this.downloadConfig.minDelay = minDelay
		}

		maxDelay, err := strconv.Atoi(delays[1])
		if err != nil {
			log.Warnf("Unable to parse max delay %s. Falling back to default %s", maxDelay, this.downloadConfig.maxDelay)
		} else {
			this.downloadConfig.maxDelay = maxDelay
		}
	}

	backupPath = getValidDir(backupPath)

	this.downloadConfig.maxParallelDownloads = maxParallelDownloads
	this.downloadConfig.backupPath = backupPath
	this.backupDir = joinToPath(backupPath, this.backupDirName)
	this.semaphore = make(chan struct{}, maxParallelDownloads)
}

func (this *DriveDownloader) WaitUntilFinished() {
	this.inProgress.Wait()
}

func (this *DriveDownloader) GetDriveItem(item *drive.File, currentPath string) {
	name := item.Name
	itemPath := this.getPath(currentPath, name)
	if item.MimeType == string(Folder) {
		nextItems := this.client.GetItems(item)
		if len(nextItems) > 0 {
			this.makeOsDir(currentPath)
		}

		hasRootParent := name == DriveRoot
		this.walkFolder(nextItems, itemPath, hasRootParent)
	} else {
		log.Trace("downloading item: %s\r\n", joinToPath(currentPath, name))
		savePath := joinToPath(this.backupDir, itemPath)

		this.inProgress.Add(1)
		this.downloadQueue <- downloadItem{item, savePath}
	}

}

func (this *DriveDownloader) getPath(currentPath string, name string) string {
	var path string
	if len(currentPath) == 0 || currentPath == DriveRoot {
		path = name
	} else {
		path = joinToPath(currentPath, name)
	}
	return path
}

func (this *DriveDownloader) fetchSingleItem(item *drive.File, savePath string) {
	defer this.inProgress.Done()

	random := rand.Intn(this.downloadConfig.maxDelay - this.downloadConfig.minDelay)
	wait, _ := time.ParseDuration(strconv.Itoa(this.downloadConfig.minDelay+random) + "s")
	log.Debugf("waiting %s before downloading: %s ...", wait, item.Name)
	time.Sleep(wait)

	readableSize := humanize.Bytes(uint64(item.Size))
	printPath := savePath[len(this.backupDir):]
	cut := len(printPath) - int(math.Min(float64(len(printPath)), 80))
	msg := fmt.Sprintf(">> downloading: %80s\t%6s", printPath[cut:], readableSize)
	log.Info(msg)

	switch DriveItemType(item.MimeType) {
	case Document:
		savePath = this.withExtension(savePath, "docx")
		this.client.ExportFile(item, savePath, string(Docx))
	case Spreadsheet:
		savePath = this.withExtension(savePath, "xlsx")
		this.client.ExportFile(item, savePath, string(Xlsx))
	case Presentation:
		savePath = this.withExtension(savePath, "pptx")
		this.client.ExportFile(item, savePath, string(Pptx))
	default:
		this.client.DownloadFile(item, savePath)
	}

	fileSize := this.getFileSize(item, savePath)
	this.DownloadStats.totalFileCount++
	this.DownloadStats.totalSize += fileSize
}

func (this *DriveDownloader) makeOsDir(currentPath string) {
	if currentPath == "root" {
		return
	}

	downloadDir := joinToPath(this.backupDir, currentPath)
	makeDir(downloadDir)
}

func (this *DriveDownloader) walkFolder(nextItems []*drive.File, itemPath string, hasRootParent bool) {
	ignoreItems, _ := this.ignoreItems[itemPath]
	isRootAnIgnoreKey := contains(this.ignoreItemsRootLevelKeys, DriveRoot)

	for _, next := range nextItems {
		isFolder := next.MimeType == string(Folder)

		isFileAtRootNode := hasRootParent && !isFolder && !isRootAnIgnoreKey
		if isFileAtRootNode {
			continue
		}

		name := next.Name
		isFolderIgnoreKeyAtRootNode := isFolder && contains(this.ignoreItemsRootLevelKeys, name)
		isNotIgnoredFileAtRootNode := !isFolder && !contains(this.ignoreItems[DriveRoot], name)
		isRootLevelInclude := hasRootParent && (isFolderIgnoreKeyAtRootNode || isNotIgnoredFileAtRootNode)
		isNonRootLevelInclude := !hasRootParent && !this.isIgnoredItem(ignoreItems, name)
		if isRootLevelInclude || isNonRootLevelInclude {
			this.GetDriveItem(next, itemPath)
		} else {
			log.Info("ignoring ", joinToPath(itemPath, name))
		}
	}
}

func (this *DriveDownloader) isIgnoredItem(ignoreItems []string, name string) bool {
	return contains(ignoreItems, name) || containsWithWildcard(ignoreItems, name)
}

func (this *DriveDownloader) withExtension(path string, extension string) string {
	if !this.hasExtension(path) {
		path += "." + extension
	}
	return path
}

func (this *DriveDownloader) hasExtension(path string) bool {
	tokens := strings.Split(path, PathSep)
	lastToken := tokens[len(tokens)-1]
	return strings.Contains(lastToken, ".")
}

func (this *DriveDownloader) getFileSize(item *drive.File, savePath string) int64 {
	fileSize := int64(0)
	fileInfo, err := os.Stat(savePath)
	if err != nil {
		log.Warnf("Could not get file size of %s: %v", savePath, err)
		fileSize = item.Size
	} else {
		fileSize = fileInfo.Size()
	}
	return fileSize
}
