package receive

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

type Batch struct {
	Metadata json.RawMessage `json:"metadata"`
	Image    string          `json:"image,omitempty"`
	Mask     *Mask           `json:"mask,omitempty"`
	Raw      string          `json:"raw,omitempty"`
}

type Mask struct {
	Image  string          `json:"image,omitempty"`
	Fruits json.RawMessage `json:"fruits,omitempty"`
}

func readBase64File(f multipart.File) (string, error) {
	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func GetMultipartElement(elementName string,
	form *multipart.Form, batch *Batch) error {
	if files, ok := form.File[elementName]; ok && len(files) > 0 {
		file, err := files[0].Open()
		if err != nil {
			return fmt.Errorf("Error processing %s in multipart: %w",
				elementName, err)
		}
		defer file.Close()

		encoded, err := readBase64File(file)
		if err != nil {
			return fmt.Errorf("Error processing %s in multipart: %w",
				elementName, err)
		}

		switch elementName {
		case "image":
			batch.Image = encoded
		case "mask":
			if batch.Mask == nil {
				batch.Mask = &Mask{}
			}
			batch.Mask.Image = encoded
		case "raw":
			batch.Raw = encoded
		}
	}

	return nil
}

func ParseMultipart(form *multipart.Form) ([]byte, error) {
	batch := Batch{}

	metadata := form.Value["metadata"]
	if len(metadata) == 0 {
		return nil, fmt.Errorf("metadata field required")
	}
	batch.Metadata = json.RawMessage(metadata[0])

	GetMultipartElement("image", form, &batch)
	GetMultipartElement("mask", form, &batch)
	GetMultipartElement("raw", form, &batch)

	return json.Marshal(batch)
}
