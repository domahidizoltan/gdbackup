package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
)

type FakeDriveClient struct{}

func (f FakeDriveClient) DownloadFile(file *drive.File, path string) {
	log.Debugf(">> download %s\r\n", path)
}

func (f FakeDriveClient) ExportFile(file *drive.File, path string, exportMimeType string) {
	log.Debugf(">> export %s as %s", path, exportMimeType)
}

func (f FakeDriveClient) GetItems(parent *drive.File) []*drive.File {
	var files = make([]*drive.File, 0)
	var parentID string
	if parent == nil {
		parentID = DriveRoot
	} else {
		parentID = parent.Id
	}

	switch parentID {
	case "DocumentsId":
		files = []*drive.File{
			&drive.File{Id: "publicId", Name: "public", Size: 0, MimeType: string(Folder)},
			&drive.File{Id: "documentsFile1", Name: "homework.doc", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "privateId", Name: "private", Size: 0, MimeType: string(Folder)},
		}
	case "privateId":
		files = []*drive.File{
			&drive.File{Id: "privateFile1", Name: "ignore-file1", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "privateFile2", Name: "ignore-file2", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "privateFile3", Name: "office1.zip", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "privateFile4", Name: "office2.zip", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "privateFile5", Name: "office1.doc", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "privateFile6", Name: "longvideo.mp4", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "privateFile7", Name: "longvideo.mpg", Size: 1, MimeType: string(Unknown)},
		}
	case "ImagesId":
		files = []*drive.File{
			&drive.File{Id: "rawfilesId", Name: "rawfiles", Size: 0, MimeType: string(Folder)},
			&drive.File{Id: "publicId", Name: "public", Size: 0, MimeType: string(Folder)},
		}
	case "publicId":
		files = []*drive.File{
			&drive.File{Id: "publicFile1", Name: "office-pics.zip", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "publicFile2", Name: "outofoffice.zip", Size: 1, MimeType: string(Unknown)},
		}
	case "VideosId":
		files = []*drive.File{
			&drive.File{Id: "videoFile1", Name: "team-building.mpg", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "videoFile2", Name: "party.mp4", Size: 1, MimeType: string(Unknown)},
		}

	case DriveRoot:
		files = []*drive.File{
			&drive.File{Id: "DocumentsId", Name: "Documents", Size: 0, MimeType: string(Folder)},
			&drive.File{Id: "ImagesId", Name: "Images", Size: 0, MimeType: string(Folder)},
			&drive.File{Id: "VideosId", Name: "Videos", Size: 0, MimeType: string(Folder)},
			&drive.File{Id: "rootFile1", Name: "shared-for-gmail.txt", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "rootFile2", Name: "shared-for-facebook.txt", Size: 1, MimeType: string(Unknown)},
			
		}

	default:
		files = []*drive.File{}
	}
	return files
}
