package services

import (
	"fmt"

	"backend/api/repository"
)

func UpdateVariant(variantID string, sku, color, storage string, price float64, stock int, images []string) error {
	err := repository.UpdateVariant(variantID, sku, color, storage, price, stock, images)
	if err != nil {
		return fmt.Errorf("failed to update variant: %w", err)
	}
	return nil
}
