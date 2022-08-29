package data

import (
	"AbuzLandingChecker/parts/pkg/tools"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UsersController struct {
	CreateHash func() (string, error)
	UpdateIP   func(IP, hashString, UA string) error
	GetAll     func() ([]Users, error)
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
			var objCheck Users
			tx := db.Model(&Users{}).Find(&objCheck, "generated_hash", hashString)
			if tx.RowsAffected == 0 {
				err := errors.New("hash not found")
				log.Error().Err(err).Msg("hash not found")
				return err
			}

			tx = db.Model(&Users{}).Find(&obj, "uniq_hash", hashString+IP)

			fmt.Println("hashString+IP", hashString+IP)
			fmt.Println("obj", obj.IP)
			fmt.Println("obj", obj.ID)
			fmt.Println("IP", IP)
			fmt.Println("QWE", obj.IP == IP)

			if obj.IP == "" {
				obj.FP = UA
				obj.IP = IP
				obj.UniqHash = hashString + IP
				obj.Count = 1
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
				newObj.UniqHash = hashString + IP
				newObj.GeneratedHash = hashString
				newObj.Count = 1
				tx = db.Create(&newObj)
				if tx.Error != nil {
					log.Error().Err(tx.Error).Msg("db create rat error")
					return tx.Error
				}
			}

			return nil
		},
		GetAll: func() ([]Users, error) {
			var obj []Users
			tx := db.Model(&Users{}).Find(&obj)

			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				return nil, tx.Error
			}
			return obj, nil
		},
	}
}
