package services

import (
	"fmt"
	"yard-calculation/models"
	"yard-calculation/repositories"
	"yard-calculation/schemas"
)

type ContainerService struct {
	Repo *repositories.ContainerRepository
}

func NewContainerService(repo *repositories.ContainerRepository) *ContainerService {
	return &ContainerService{Repo: repo}
}

func (s *ContainerService) GetSuggestedPosition(yardName, containerNumber string, size int, height float64, ctype string) (*schemas.SuggestContainerResponse, error) {
	// Ambil data yard dan blocks beserta plans
	yard, err := s.Repo.GetYardByName(yardName)
	if err != nil {
		return nil, err
	}

	// Cari posisi yang sesuai dengan rencana di *semua* block dalam yard
	suggestedContainer, err := s.Repo.FindSuggestedPosition(yard, size, height, ctype)
	if err != nil {
		return nil, err
	}

	// Isi nomor kontainer ke hasil saran
	suggestedContainer.ContainerNumber = containerNumber

	response := schemas.SuggestContainerResponse{
		Yard:  suggestedContainer.YardID,
		Block: suggestedContainer.BlockID,
		Slot:  suggestedContainer.Slot,
		Row:   suggestedContainer.Row,
		Tier:  suggestedContainer.Tier,
	}

	return &response, nil
}

// Ubah fungsi PlaceContainer untuk menerima informasi kontainer
func (s *ContainerService) PlaceContainerDetailed(yardName, containerNumber, blockName string, slot, row, tier int, size int, height float64, ctype string) error {
	block, err := s.Repo.GetBlockByName(blockName, yardName)
	if err != nil {
		return err
	}

	if err := s.Repo.LoadBlockOccupancy(block); err != nil {
		return fmt.Errorf("error loading block occupancy: %v", err)
	}

	// Validasi batas
	if slot < 1 || slot > block.TotalSlot || row < 1 || row > block.TotalRow || tier < 1 || tier > block.TotalTier {
		return fmt.Errorf("position out of bounds for block %s", blockName)
	}

	// Cek apakah kontainer sudah ditempatkan
	existingPlacedContainer, err := s.Repo.GetContainerByNumber(containerNumber)
	if err != nil && err.Error() != fmt.Sprintf("container with number %s not found or not placed", containerNumber) {
		return err // Error lain
	}
	if existingPlacedContainer != nil {
		return fmt.Errorf("container with number %s is already placed at %s-%d-%d-%d", containerNumber, existingPlacedContainer.BlockID, existingPlacedContainer.Slot, existingPlacedContainer.Row, existingPlacedContainer.Tier)
	}

	// Validasi ketersediaan posisi berdasarkan ukuran
	if size == 20 {
		if !s.Repo.IsPositionAvailable(block, slot, row, tier) {
			return fmt.Errorf("position %d-%d-%d in block %s is occupied", slot, row, tier, blockName)
		}
	} else if size == 40 {
		if !s.Repo.IsPositionAvailable40ft(block, slot, row, tier) {
			return fmt.Errorf("positions %d-%d-%d and %d-%d-%d in block %s are not available for 40ft container", slot, row, tier, slot+1, row, tier, blockName)
		}
	} else {
		return fmt.Errorf("unsupported container size: %d", size)
	}

	// --- Validasi Penempatan Sesuai Rencana ---
	// Cek apakah ada rencana yang sesuai untuk spesifikasi kontainer di posisi yang dituju
	plans, err := s.Repo.GetPlansForSpec(yardName, blockName, size, height, ctype)
	if err != nil {
		return fmt.Errorf("error checking placement plan: %v", err)
	}
	if len(plans) == 0 {
		return fmt.Errorf("no plan found for container spec (size: %d, height: %.1f, type: %s) at placement location in block %s", size, height, ctype, blockName)
	}

	// Cek apakah posisi (slot, row, tier) atau (slot, slot+1, row, tier) masuk ke salah satu plan
	validLocation := false
	for _, p := range plans {
		if tier >= p.MinTier && tier <= p.MaxTier &&
			row >= p.MinRow && row <= p.MaxRow &&
			slot >= p.MinSlot && slot <= p.MaxSlot {
			if size == 40 { // Jika 40ft, cek slot+1 juga
				if slot+1 >= p.MinSlot && slot+1 <= p.MaxSlot {
					validLocation = true
					break
				}
			} else { // Jika 20ft
				validLocation = true
				break
			}
		}
	}
	if !validLocation {
		return fmt.Errorf("placement location (%d-%d-%d) does not match planned area for container spec (size: %d, height: %.1f, type: %s) in block %s", slot, row, tier, size, height, ctype, blockName)
	}
	// --- Akhir Validasi Penempatan Sesuai Rencana ---

	// Buat entitas Container untuk disimpan
	containerToPlace := &models.Container{
		ContainerNumber: containerNumber,
		Size:            size,
		Height:          height,
		Type:            ctype,
		YardID:          yardName,
		BlockID:         blockName,
		Slot:            slot,
		Row:             row,
		Tier:            tier,
		IsPlaced:        true,
		// YardPlanID:      nil, // Atau set ke ID plan jika diperlukan
	}

	return s.Repo.CreateContainer(containerToPlace)
}

func (s *ContainerService) PickupContainer(yardName, containerNumber string) error {
	container, err := s.Repo.GetContainerByNumber(containerNumber)
	if err != nil {
		return err
	}

	if container.YardID != yardName {
		return fmt.Errorf("container %s is not located in yard %s", containerNumber, yardName)
	}

	// Update status menjadi tidak ditempatkan
	container.IsPlaced = false
	return s.Repo.UpdateContainer(container)
	// Secara logika, posisi sekarang "kosong", GORM akan menyimpan perubahan IsPlaced
	// Jika menggunakan cache Occupancy di Block, perlu diupdate juga disana.
	// Kita abaikan cache Occupancy untuk sementara atau update saat load ulang.
}
