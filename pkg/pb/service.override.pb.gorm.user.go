package pb

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
	
	"github.com/infobloxopen/atlas-app-toolkit/gorm/resource"
)

// AfterToORM will marshall Group from request to ORM object
func (m *User) AfterToORM(ctx context.Context, a *UserORM) error {
	for _, g := range m.GroupList {
		containerTagId, err := resource.DecodeInt64(&Group{}, g)
		if err != nil {
			return err
		}
		a.GroupList = append(a.GroupList, &GroupORM{Id: containerTagId})
	}
	
	return nil
}

// AfterToPB copies the Group resource information from ORM object to PB object
func (m *UserORM) AfterToPB(ctx context.Context, a *User) error {
	
	for _, ct := range m.GroupList {
		containerTagId, err := resource.Encode(&Group{}, ct.Id)
		if err != nil {
			return err
		}
		a.GroupList = append(a.GroupList, containerTagId)
	}
	return nil
}

func (m *UserORM) BeforeStrictUpdateSave(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	if err := db.Model(m).Association("Groups").Replace(m.GroupList).Error; err != nil {
		return nil, err
	}
	return db, nil
}
