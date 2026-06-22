package repository

import (
	"encoding/json"
	"fmt"

	"backend/pkg/database"
)

func UpdateVariant(variantID string, sku, color, storage string, price float64, stock int, images []string) error {
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return fmt.Errorf("failed to marshal images: %w", err)
	}

	_, err = database.DB.Exec(`
		UPDATE product_variants 
		SET sku = $1, color = $2, storage = $3, price = $4, stock = $5, images = $6, updated_at = NOW()
		WHERE id = $7`,
		sku, color, storage, price, stock, imagesJSON, variantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update variant: %w", err)
	}
	return nil
}
