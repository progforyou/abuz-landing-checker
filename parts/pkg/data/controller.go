package data

import (
	"AbuzLandingChecker/parts/pkg/tools"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UsersController struct {
	CreateHash func() (string, error)
	UpdateIP   func(IP, hashString, UA string) error
}

func NewUsersController(db *gorm.DB, baseLog zerolog.Logger) UsersController {
	log := baseLog.With().Str("model", "users").Logger()
	if err := db.AutoMigrate(&Users{}); err != nil {
		log.Fatal().Err(err).Msg("auto-migrate")
	}

	return UsersController{
		CreateHash: func() (string, error) {
			var obj Users
			var generatedHash string

			for {
				generatedHash = tools.RandStringBytes(25)
				obj.GeneratedHash = generatedHash
				tx := db.Model(&Users{}).Find(&obj, "generated_hash", generatedHash)
				if tx.RowsAffected == 0 {
					break
				}
			}

			tx := db.Model(&Users{}).Create(&obj)
			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				return "", tx.Error
			}
			return generatedHash, nil
		},
		UpdateIP: func(IP, hashString, UA string) error {
			var obj Users
			tx := db.Model(&Users{}).Find(&obj, "generated_hash", hashString)

			if obj.IP == "" {
				obj.FP = UA
				obj.IP = IP
				tx = db.Save(&obj)
				if tx.Error != nil {
					log.Error().Err(tx.Error).Msg("db update ip new error")
					return tx.Error
				}
				return nil
			}

			if obj.IP == IP {
				obj.Count += 1
				tx = db.Save(&obj)
				if tx.Error != nil {
					log.Error().Err(tx.Error).Msg("db update ip error")
					return tx.Error
				}
			} else {
				var newObj Users
				newObj.FP = UA
				newObj.IP = IP
				newObj.GeneratedHash = hashString
				tx = db.Create(&newObj)
				if tx.Error != nil {
					log.Error().Err(tx.Error).Msg("db create rat error")
					return tx.Error
				}
			}

			return nil
		},
	}
}
