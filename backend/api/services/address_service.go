package services

import (
	"fmt"

	"backend/api/models"
	"backend/api/repository"
)

// GetAddresses fetches all addresses for a user
func GetAddresses(userID string) ([]models.Address, error) {
	addresses, err := repository.GetAddressesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}
	return addresses, nil
}

// CreateAddress validates and creates a new address for a user
func CreateAddress(userID string, req models.CreateAddressRequest) (*models.Address, error) {
	// Check user does not exceed 5 saved addresses
	existing, err := repository.GetAddressesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check addresses: %w", err)
	}
	if len(existing) >= 5 {
		return nil, fmt.Errorf("maximum of 5 saved addresses allowed")
	}

	address, err := repository.CreateAddress(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return address, nil
}

// UpdateAddress validates and updates an existing address
func UpdateAddress(addressID, userID string, req models.UpdateAddressRequest) error {
	existing, err := repository.GetAddressByID(addressID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch address: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("address not found")
	}

	return repository.UpdateAddress(addressID, userID, req)
}

// DeleteAddress removes an address belonging to a user
func DeleteAddress(addressID, userID string) error {
	existing, err := repository.GetAddressByID(addressID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch address: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("address not found")
	}

	// Prevent deleting the default address if other addresses exist
	if existing.IsDefault {
		all, _ := repository.GetAddressesByUserID(userID)
		if len(all) > 1 {
			return fmt.Errorf("cannot delete default address, set another address as default first")
		}
	}

	return repository.DeleteAddress(addressID, userID)
}

// SetDefaultAddress sets an address as the user's default delivery address
func SetDefaultAddress(addressID, userID string) error {
	existing, err := repository.GetAddressByID(addressID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch address: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("address not found")
	}

	req := models.UpdateAddressRequest{
		Label:     existing.Label,
		FullName:  existing.FullName,
		Phone:     existing.Phone,
		County:    existing.County,
		Town:      existing.Town,
		Street:    existing.Street,
		IsDefault: true,
	}

	return repository.UpdateAddress(addressID, userID, req)
}
