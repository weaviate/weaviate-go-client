package graphql

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

type NearImageArgumentBuilder struct {
	image         string
	fileReader    io.Reader
	withCertainty bool
	certainty     float32
	withDistance  bool
	distance      float32
}

// WithImage base64 encoded image
func (b *NearImageArgumentBuilder) WithImage(image string) *NearImageArgumentBuilder {
	b.image = image
	return b
}

// WithReader the image file
func (b *NearImageArgumentBuilder) WithReader(fileReader io.Reader) *NearImageArgumentBuilder {
	b.fileReader = fileReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearImageArgumentBuilder) WithCertainty(certainty float32) *NearImageArgumentBuilder {
	b.withCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearImageArgumentBuilder) WithDistance(distance float32) *NearImageArgumentBuilder {
	b.withDistance = true
	b.distance = distance
	return b
}

func (b *NearImageArgumentBuilder) getImage(image string, fileReader io.Reader) string {
	if fileReader != nil {
		content, err := io.ReadAll(fileReader)
		if err != nil {
			return err.Error()
		}
		return base64.StdEncoding.EncodeToString(content)
	}
	if strings.HasPrefix(image, "data:") {
		base64 := ";base64,"
		indx := strings.LastIndex(image, base64)
		return image[indx+len(base64):]
	}
	return image
}

// Build build the given clause
func (b *NearImageArgumentBuilder) build() string {
	clause := []string{}
	if len(b.image) > 0 || b.fileReader != nil {
		clause = append(clause, fmt.Sprintf("image: \"%s\"", b.getImage(b.image, b.fileReader)))
	}
	if b.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", b.certainty))
	}
	if b.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", b.distance))
	}
	return fmt.Sprintf("nearImage:{%s}", strings.Join(clause, " "))
}
