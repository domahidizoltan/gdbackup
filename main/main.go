package main

import (
	"flag"

	humanize "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
)

func main() {
	argLogLevel := flag.String("loglevel", DefaultLogLevel, "possible leg levels: trace, debug, info, warn, error, fatal, panic")
	gdIgnorePath := flag.String("gdignore-path", DefaultGDIgnorePath, "path to gdignore.yaml")
	delay := flag.String("delay", DefaultDownloadDelay, "min-max random delay range in seconds before downloading a file")
	maxParallelDownloads := flag.Int("max-parallel-downloads", DefaultMaxParallelDownloads, "max parallel file downloads")
	flag.Parse()

	configureLogger(argLogLevel)

	ignoreItems := NewIgnoreItems(*gdIgnorePath)
	log.Debug("ignoreItems: %v\r\n", ignoreItems)

	// client := NewOauth2Client()
	// driveClient := NewDriveClient(client)
	driveClient := &FakeDriveClient{}
	rootFolder := &drive.File{Id: DriveRoot, Name: DriveRoot, MimeType: string(Folder), Size: 0}
	downloader := NewDownloader(driveClient, ignoreItems)
	downloader.ConfigureDownloader(*delay, *maxParallelDownloads)

	downloader.GetDriveItem(rootFolder, DriveRoot)
	downloader.WaitUntilFinished()

	fileCount := downloader.DownloadStats.totalFileCount
	totalSize := humanize.Bytes(uint64(downloader.DownloadStats.totalSize))
	log.Infof("Finished backup of %d files [%s]", fileCount, totalSize)

}

func configureLogger(argLogLevel *string) {
	logLevel, err := log.ParseLevel(*argLogLevel)
	if err != nil {
		log.Warnf("Unable to parse loglevel. Falling back to info.")
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "060102 15:04:05.000",
	})
}
