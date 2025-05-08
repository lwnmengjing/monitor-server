package system

import (
	"runtime"

	"github.com/mss-boot-io/mss-boot/pkg/migration"
	migrationModel "github.com/mss-boot-io/mss-boot/pkg/migration/models"
	"gorm.io/gorm"

	"service-http/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1691804837583Tables)
}

func _1691804837583Tables(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		err := tx.Migrator().AutoMigrate(
			new(models.Model),
			new(models.Field),
		)
		if err != nil {
			return err
		}

		return tx.Create(&migrationModel.Migration{
			Version: version,
		}).Error
	})
}
