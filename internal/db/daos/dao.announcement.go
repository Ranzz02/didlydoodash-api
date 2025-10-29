package daos

import (
	"github.com/Stenoliv/didlydoodash_api/internal/db"
	"github.com/Stenoliv/didlydoodash_api/internal/db/models"
)

func GetAnnouncements(OrganisationID string) ([]models.Announcement, error) {
	var announcement []models.Announcement
	err := db.DB.Model(&models.Announcement{}).Where("organisation_id = ? ", OrganisationID).Find(&announcement).Error

	if err != nil {
		return nil, err
	}
	return announcement, nil
}
func GetAnnouncement(aID string) (*models.Announcement, error) {
	var a *models.Announcement

	err := db.DB.Model(&models.Announcement{}).
		Where("id = ?", aID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return a, nil

}
