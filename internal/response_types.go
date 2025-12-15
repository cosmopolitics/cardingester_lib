package cardingester
import (
	"time"
)

type Bulk_Response struct {
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	Data    []Data `json:"data"`
}

type Data struct {
	Object          string    `json:"object"`
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	UpdatedAt       time.Time `json:"updated_at"`
	URI             string    `json:"uri"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Size            int       `json:"size"`
	DownloadURI     string    `json:"download_uri"`
	ContentType     string    `json:"content_type"`
	ContentEncoding string    `json:"content_encoding"`
}
