package main

import (
	"flag"
	"os"

	humanize "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
)

func main() {
	argLogLevel := flag.String("loglevel", DefaultLogLevel, "possible leg levels: trace, debug, info, warn, error, fatal, panic")
	gdIgnorePath := flag.String("gdignore-path", DefaultGDIgnorePath, "path to gdignore.yaml")
	delay := flag.String("delay", DefaultDownloadDelay, "min-max random delay range in seconds before downloading a file")
	maxParallelDownloads := flag.Int("max-parallel-downloads", DefaultMaxParallelDownloads, "max parallel file downloads")
	backupPath := flag.String("backup-path", DefaultBackupPath, "path to save backup files")
	flag.Parse()

	configureLogger(argLogLevel, *backupPath)

	ignoreItems := NewIgnoreItems(*gdIgnorePath)
	log.Debugf("ignoreItems: %v\r\n", ignoreItems)

	client := NewOauth2Client()
	driveClient := NewDriveClient(client)
	// driveClient := &FakeDriveClient{}
	rootFolder := &drive.File{Id: DriveRoot, Name: DriveRoot, MimeType: string(Folder), Size: 0}
	downloader := NewDownloader(driveClient, ignoreItems)
	downloader.ConfigureDownloader(*delay, *maxParallelDownloads, *backupPath)

	downloader.GetDriveItem(rootFolder, DriveRoot)
	downloader.WaitUntilFinished()

	fileCount := downloader.DownloadStats.totalFileCount
	totalSize := humanize.Bytes(uint64(downloader.DownloadStats.totalSize))
	log.Infof("Finished backup of %d files [%s]", fileCount, totalSize)

}

func configureLogger(argLogLevel *string, path string) {
	path = getValidDir(path)
	logFilePath := joinToPath(path, DefaultLogFile)
	log.Infof("Logging to file " + logFilePath)
	logFile, err := os.OpenFile(logFilePath, os.O_TRUNC | os.O_CREATE | os.O_RDWR, 0666)
    if err != nil {
        log.Errorf("Error opening file: %v", err)
    }
	log.SetOutput(logFile)


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
