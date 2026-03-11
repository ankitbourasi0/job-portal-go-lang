package repository

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ankitbourasi0/job-portal/internal/database"
	"github.com/ledongthuc/pdf"
)

type GuestRepository struct {
	Queries *database.Queries
}

func NewGuestRepository(dbQueries *database.Queries) *GuestRepository {
	return &GuestRepository{Queries: dbQueries}
}

// Convert PDF to Text,
// accept: ReaderAt , Size of the file
func (r *GuestRepository) PdfParser(content io.ReaderAt, size int64) (string, error) {
	//1. Create Reader
	reader, err := pdf.NewReader(content, size)
	if err != nil {
		return "", fmt.Errorf("failed to create pdf reader: %v", err)
	}

	var buffer bytes.Buffer

	//2. Extract text page by page
	numberOfPages := reader.NumPage()
	for i := 0; i < numberOfPages; i++ {
		page := reader.Page(i)
		//Check if the page's value is null or if its content key is null then
		if page.V.IsNull() || page.V.Key("Contents").Kind() == pdf.Null {
			continue //skip the page
		}

		// this  method extract only text content
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to read pdf text: %v", err)
		}
		//add the text in buffer
		buffer.WriteString(text)
	}

	return buffer.String(), nil
}
