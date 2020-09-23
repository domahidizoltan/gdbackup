package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
)

type DriveClient interface {
	GetItems(parent *drive.File) []*drive.File
	DownloadFile(file *drive.File, path string)
	ExportFile(file *drive.File, path string, exportMimeType string)
}

type DefaultDriveClient struct {
	service *drive.Service
}

type DriveItemType string

const (
	Folder       DriveItemType = "application/vnd.google-apps.folder"
	Document     DriveItemType = "application/vnd.google-apps.document"
	Spreadsheet  DriveItemType = "application/vnd.google-apps.spreadsheet"
	Presentation DriveItemType = "application/vnd.google-apps.presentation"
	Unknown      DriveItemType = "application/vnd.google-apps.unknown"
)

func NewDriveClient(client *http.Client) DefaultDriveClient {
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return DefaultDriveClient{
		service: srv,
	}
}

func (client DefaultDriveClient) GetItems(parent *drive.File) []*drive.File {
	query := fmt.Sprintf("'%s' in parents and trashed=false", parent.Id)
	res, err := client.service.Files.List().
		Spaces("drive").
		Q(query).
		IncludeItemsFromAllDrives(false).
		Fields("files(id, name, size, mimeType)").
		Do()

	if err != nil {
		log.Errorf("Unable to retrieve files from %s: %v", err, parent.Name)
	}

	return res.Files
}

func (client DefaultDriveClient) DownloadFile(file *drive.File, path string) {
	saveResponse(path, createFile, func() (*http.Response, error) {
		return client.service.Files.Get(file.Id).Download()
	})
}

func (client DefaultDriveClient) ExportFile(file *drive.File, path string, exportMimeType string) {
	saveResponse(path, createFile, func() (*http.Response, error) {
		return client.service.Files.Export(file.Id, exportMimeType).Download()
	})
}

func createFile(path string) *os.File {
	out, err := os.Create(path)
	if err != nil {
		log.Panicf("Unable to create file: %v", err)
	}
	return out
}

func saveResponse(path string, fileAction func(string) *os.File, driveClientAction func() (*http.Response, error)) {
	outFile := fileAction(path)
	defer outFile.Close()

	response, err := driveClientAction()
	if err != nil {
		log.Errorf("Unable to download file: %v", err)
		if os.Remove(path) != nil {
			log.Warn("Unable to delete create file at " + path)
		}
	} else {
		log.Tracef("Headers %v:", response.Header)
		io.Copy(outFile, response.Body)
	}
	defer response.Body.Close()

	checkDownloadSize(response.Header, outFile)
}

func checkDownloadSize(responseHeaders http.Header, outFile *os.File) {
	contentLength := responseHeaders.Get("Content-Length")
	responseSize, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Warnf("Unable to parse Content-Length: %v", err)
		responseSize = 0
	}

	fileInfo, err := outFile.Stat()
	if err != nil {
		log.Warnf("Unable to get file info for %s: %v", outFile.Name(), err)
	}

	if int64(responseSize) != fileInfo.Size() {
		log.Warnf("Downloaded file size %d does not match %d", responseSize, fileInfo.Size())
	}
}
