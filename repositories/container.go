package repositories

import (
	"errors"
	"fmt"
	"sort"
	"yard-calculation/models"

	"gorm.io/gorm"
)

type ContainerRepository struct {
	DB *gorm.DB
}

func NewContainerRepository(db *gorm.DB) *ContainerRepository {
	return &ContainerRepository{DB: db}
}

func (r *ContainerRepository) GetYardByName(name string) (*models.Yard, error) {
	var yard models.Yard
	// Preload Plans juga
	if err := r.DB.Preload("Blocks.Plans").First(&yard, "id = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("yard with name %s not found", name)
		}
		return nil, err
	}
	return &yard, nil
}

func (r *ContainerRepository) GetBlockByName(blockName string, yardID string) (*models.Block, error) {
	var block models.Block
	if err := r.DB.First(&block, "id = ? AND yard_id = ?", blockName, yardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("block with name %s in yard %s not found", blockName, yardID)
		}
		return nil, err
	}
	return &block, nil
}

// Simulasi pengisian Occupancy dari database
func (r *ContainerRepository) LoadBlockOccupancy(block *models.Block) error {
	if block.Occupancy == nil {
		block.Occupancy = make(map[string]bool)
	}

	var containers []models.Container
	// Hanya muat kontainer yang ditempatkan di block ini
	if err := r.DB.Where("block_id = ? AND is_placed = ?", block.ID, true).Find(&containers).Error; err != nil {
		return err
	}

	for _, c := range containers {
		// Untuk kontainer 40ft, tandai dua slot sebagai terisi
		block.Occupancy[fmt.Sprintf("%d-%d-%d", c.Slot, c.Row, c.Tier)] = true
		if c.Size == 40 {
			block.Occupancy[fmt.Sprintf("%d-%d-%d", c.Slot+1, c.Row, c.Tier)] = true
		}
	}
	return nil
}

// Simulasi pengecekan apakah posisi kosong
func (r *ContainerRepository) IsPositionAvailable(block *models.Block, slot, row, tier int) bool {
	key := fmt.Sprintf("%d-%d-%d", slot, row, tier)
	occupied, exists := block.Occupancy[key]
	return !exists || !occupied
}

// Simulasi pengecekan apakah posisi untuk container 40ft kosong
func (r *ContainerRepository) IsPositionAvailable40ft(block *models.Block, slot, row, tier int) bool {
	if slot+1 > block.TotalSlot { // Gunakan TotalSlot
		return false
	}
	if !r.IsPositionAvailable(block, slot, row, tier) || !r.IsPositionAvailable(block, slot+1, row, tier) {
		return false
	}
	return true
}

func (r *ContainerRepository) GetPlansForSpec(yardID, blockID string, size int, height float64, ctype string) ([]models.YardPlan, error) {
	var plans []models.YardPlan
	if err := r.DB.Where("yard_id = ? AND block_id = ? AND planned_size = ? AND planned_height = ? AND planned_type = ?", yardID, blockID, size, height, ctype).Find(&plans).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Tidak ada rencana spesifik, kembalikan slice kosong
			return []models.YardPlan{}, nil
		}
		return nil, err
	}
	// Urutkan rencana berdasarkan prioritas (misalnya, tier -> row -> slot terendah dulu)
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].MinTier != plans[j].MinTier {
			return plans[i].MinTier < plans[j].MinTier
		}
		if plans[i].MinRow != plans[j].MinRow {
			return plans[i].MinRow < plans[j].MinRow
		}
		return plans[i].MinSlot < plans[j].MinSlot
	})
	return plans, nil
}

func (r *ContainerRepository) FindSuggestedPosition(yard *models.Yard, size int, height float64, ctype string) (*models.Container, error) {
	// Iterasi semua blocks di yard
	for _, block := range yard.Blocks {
		plans, err := r.GetPlansForSpec(yard.ID, block.ID, size, height, ctype)
		if err != nil {
			// Log error dan lanjutkan ke block berikutnya
			fmt.Printf("Error getting plans for block %s: %v\n", block.ID, err)
			continue
		}
		if len(plans) == 0 {
			// Tidak ada rencana untuk spesifikasi ini di block ini, lanjutkan
			continue
		}

		if err := r.LoadBlockOccupancy(&block); err != nil {
			// Log error dan lanjutkan ke block berikutnya
			fmt.Printf("Error loading occupancy for block %s: %v\n", block.ID, err)
			continue
		}

		// Iterasi dalam rencana-rencana block ini
		for _, plan := range plans {
			// Iterasi dalam area rencana (Tier -> Row -> Slot)
			for t := plan.MinTier; t <= plan.MaxTier; t++ {
				for r_idx := plan.MinRow; r_idx <= plan.MaxRow; r_idx++ {
					for s := plan.MinSlot; s <= plan.MaxSlot; s++ {
						// Pastikan tetap dalam batas total block
						if s > block.TotalSlot || r_idx > block.TotalRow || t > block.TotalTier {
							continue
						}
						if size == 20 {
							if r.IsPositionAvailable(&block, s, r_idx, t) {
								suggested := &models.Container{
									YardID:  yard.ID,
									BlockID: block.ID, // Gunakan block.ID dari iterasi
									// YardPlanID: &plan.ID, // Jika ingin mengikat ke plan
									Slot: s, Row: r_idx, Tier: t,
									Size: size, Height: height, Type: ctype,
								}
								return suggested, nil
							}
						} else if size == 40 {
							// Pastikan slot berikutnya juga dalam area rencana dan total block
							if s+1 > plan.MaxSlot || s+1 > block.TotalSlot {
								continue
							}
							if r.IsPositionAvailable40ft(&block, s, r_idx, t) {
								suggested := &models.Container{
									YardID:  yard.ID,
									BlockID: block.ID, // Gunakan block.ID dari iterasi
									// YardPlanID: &plan.ID, // Jika ingin mengikat ke plan
									Slot: s, Row: r_idx, Tier: t,
									Size: size, Height: height, Type: ctype,
								}
								return suggested, nil
							}
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("no suitable position found within planned areas for container spec (size: %d, height: %.1f, type: %s) in any block of yard %s", size, height, ctype, yard.ID)
}

func (r *ContainerRepository) CreateContainer(container *models.Container) error {
	// Pastikan IsPlaced di set true saat ditempatkan
	container.IsPlaced = true
	return r.DB.Create(container).Error
}

func (r *ContainerRepository) GetContainerByNumber(containerNumber string) (*models.Container, error) {
	var container models.Container
	if err := r.DB.Where("container_number = ? AND is_placed = ?", containerNumber, true).First(&container).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("container with number %s not found or not placed", containerNumber)
		}
		return nil, err
	}
	return &container, nil
}

func (r *ContainerRepository) UpdateContainer(container *models.Container) error {
	// Misalnya, saat pickup, set IsPlaced ke false
	return r.DB.Save(container).Error
}

// Fungsi untuk melepaskan posisi di block occupancy (opsional, bisa juga di update saja)
func (r *ContainerRepository) ReleasePosition(block *models.Block, slot, row, tier int) {
	key := fmt.Sprintf("%d-%d-%d", slot, row, tier)
	delete(block.Occupancy, key)
}
