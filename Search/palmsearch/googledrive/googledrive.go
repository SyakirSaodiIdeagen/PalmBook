package googledrive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"palmsearch/elasticsearch"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveFile struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"Name,omitempty"`
	WebUrl         string `json:"webUrl,omitempty"`
	DownloadLink   string `json:"downloadLink,omitempty"`
	Source         string `json:"source,omitempty"`
	CreatedAt      string `json:"createdDateTime,omitempty"`
	LastModifiedAt string `json:"lastModifiedDateTime,omitempty"`
}

func IndexGoogleDrive() {
	var filesData []DriveFile
	ctx := context.Background()

	jsonKey, err := os.ReadFile("googledrive/gd2.json")
	if err != nil {
		log.Fatalf("Failed to read JSON key file: %v", err)
	}

	creds, err := google.CredentialsFromJSON(ctx, jsonKey, drive.DriveScope)
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	service, err := drive.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create Drive service: %v", err)
	}

	filesList, err := service.Files.List().
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Q("mimeType != 'application/vnd.google-apps.folder'").
		Fields("files(id, name, webViewLink, webContentLink, createdTime, modifiedTime)").
		Do()

	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	if len(filesList.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, file := range filesList.Files {
			driveFile := DriveFile{
				ID:             file.Id,
				Name:           file.Name,
				WebUrl:         file.WebViewLink,
				DownloadLink:   file.WebContentLink,
				CreatedAt:      file.CreatedTime,
				LastModifiedAt: file.ModifiedTime,
				Source:         "GoogleDrive",
			}
			filesData = append(filesData, driveFile)
		}
	}

	for _, file := range filesData {
		fmt.Printf("File: %s (ID: %s)\n", file.Name, file.ID)
		fmt.Printf("View Link: %s\n", file.WebUrl)
		fmt.Printf("Download Link: %s\n\n", file.DownloadLink)
	}

	buf := getElasticBuf(filesData)
	esclient := elasticsearch.GetEsClient()
	elasticsearch.BulkInsert(esclient, buf)

}

func getElasticBuf(data []DriveFile) bytes.Buffer {

	var buf bytes.Buffer
	for _, doc := range data {
		id := "gd" + "_" + doc.ID
		//c.Set(id, "golang-bulk-index3", cache.DefaultExpiration)

		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "golang-bulk-index3", "_id" : "%v" } }%s`, id, "\n"))
		d, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("Cannot encode document %d: %s", doc.ID, err)
		}

		buf.Grow(len(meta) + len(d) + 1)
		buf.Write(meta)
		buf.Write(d)
		buf.WriteByte('\n')
	}

	return buf
}
