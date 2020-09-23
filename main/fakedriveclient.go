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
	case "OsztasId":
		files = []*drive.File{
			&drive.File{Id: "BeazasId", Name: "Beazas", Size: 0, MimeType: string(Folder)},
			&drive.File{Id: "OsztasId1", Name: "OsztasFile", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "ValamiId", Name: "valami", Size: 0, MimeType: string(Folder)},
		}
	case "ValamiId":
		files = []*drive.File{
			&drive.File{Id: "ValamiId1", Name: "v1", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "ValamiId2", Name: "v2", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "ValamiId3", Name: "v3", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "ValamiId4", Name: "b1.zip", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "ValamiId5", Name: "b2.zip", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "ValamiId6", Name: "c.zip", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "ValamiId7", Name: "d1.mp4", Size: 1, MimeType: string(Unknown)},
			&drive.File{Id: "ValamiId8", Name: "e2.mp4", Size: 1, MimeType: string(Unknown)},
		}
	case "RootFolderId":
		files = []*drive.File{
			&drive.File{Id: "RootFolderId1", Name: "rf1", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "RootFolderId2", Name: "rf2", Size: 1, MimeType: string(Document)},
		}

	case DriveRoot:
		files = []*drive.File{
			&drive.File{Id: "OsztasId", Name: "Osztas", Size: 0, MimeType: string(Folder)},
			&drive.File{Id: "root2", Name: "test", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "root3", Name: "rootfile", Size: 1, MimeType: string(Document)},
			&drive.File{Id: "RootFolderId", Name: "RootFolder", Size: 0, MimeType: string(Folder)},
		}

	default:
		files = []*drive.File{}
	}
	return files
}
