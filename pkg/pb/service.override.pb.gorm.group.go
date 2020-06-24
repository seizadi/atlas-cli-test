package pb

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"

	"github.com/infobloxopen/atlas-app-toolkit/gorm/resource"
)

// AfterToORM will marshall User from request to ORM object
func (m *Group) AfterToORM(ctx context.Context, a *GroupORM) error {
	for _, item := range m.UserList {
		id, err := resource.DecodeInt64(&User{}, item)
		if err != nil {
			return err
		}
		a.UserList = append(a.UserList, &UserORM{Id: id})
	}

	return nil
}

// AfterToPB copies the User resource information from ORM object to PB object
// Construct the user_list field in the PB message as a list of ID's of
// all the users that belong in this group. The list of all the users is
// available in the GroupORM object.
func (m *GroupORM) AfterToPB(ctx context.Context, a *Group) error {

	for _, item := range m.UserList {
		id, err := resource.Encode(&User{}, item.Id)
		if err != nil {
			return err
		}
		a.UserList = append(a.UserList, id)
	}
	return nil
}

func (m *GroupORM) BeforeStrictUpdateSave(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	if err := db.Model(m).Association("Users").Replace(m.UserList).Error; err != nil {
		return nil, err
	}
	return db, nil
}
