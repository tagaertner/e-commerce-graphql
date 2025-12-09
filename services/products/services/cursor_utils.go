package services

import (
	"encoding/base64"
	"fmt"
	
	
	"github.com/tagaertner/e-commerce-graphql/services/products/models"
)

// EncodeCursor creates a cursor from a story
func EncodeCursor(p *models.Product) string{
	cursor := fmt.Sprintf("%s", p.ID)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}

// DecodeCursor extracts timestamp and ID from cursor
func DecodeCursor(cursor string) (string, error) {
    decoded, err := base64.StdEncoding.DecodeString(cursor)
    if err != nil {
        return "", err
    }
    return string(decoded), nil
}