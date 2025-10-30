package models

import (
	"time"

	"github.com/Stenoliv/didlydoodash_api/internal/db/datatypes"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Organisation struct {
	Base

	// Identifiers and ownership
	Name        string  `gorm:"size:100;not null;uniqueIndex;" json:"name"`
	Slug        string  `gorm:"size:120;not null;uniqueIndex;" json:"slug"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	OwnerID     string  `gorm:"not null" json:"-"`
	Owner       User    `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"owner"`

	// Contact / metadata
	Website    *string    `gorm:"" json:"website,omitempty"`
	LogoURL    *string    `gorm:"" json:"logoUrl,omitempty"`
	Location   *string    `gorm:"" json:"location,omitempty"`
	Timezone   string     `gorm:"default:'UTC';" json:"timezone"`
	IsActive   bool       `gorm:"default:true" json:"isActive"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty"`

	// Business / usage context
	Settings datatypes.JSONB `gorm:"type:jsonb;default:'{}'" json:"settings"`

	// Derived / relations
	Members  []OrganisationMember `gorm:"foreignKey:OrganisationID;" json:"members,omitempty"`
	Projects []Project            `gorm:"foreignKey:OrganisationID;" json:"projects"`
	Chats    []ChatRoom           `gorm:"foreignKey:OrganisationID;" json:"chatRooms,omitempty"`
}

func (o *Organisation) SaveOrganisation(tx *gorm.DB) (err error) {
	// Set slug automatically
	o.Slug = slug.Make(o.Name)

	if err = tx.Create(&o).Error; err != nil {
		return err
	}
	return nil
}

func (o *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	err = o.GenerateID()
	if err != nil {
		return err
	}

	return nil
}

func (o *Organisation) AfterFind(tx *gorm.DB) (err error) {
	if o.Owner.ID == "" {
		if err := tx.Preload("Owner").Find(&o).Error; err != nil {
			return err
		}
	}
	return nil
}
