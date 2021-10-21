package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// get semua campaign tanpa parameter
func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	// memanggil relasi, dengan Preload("NamaField struct", "condition")
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// get campaign berdasarkan userID
func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign

	// memanggil relasi, dengan Preload("NamaField struct", "condition")
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
