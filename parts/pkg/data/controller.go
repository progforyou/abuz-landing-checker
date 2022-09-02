package data

import (
	"AbuzLandingChecker/parts/pkg/fp"
	"AbuzLandingChecker/parts/pkg/tools"
	"errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UsersController struct {
	CreateHash func() (string, error)
	UpdateIP   func(IP, hashString, UA, city string) error
	GetAll     func() ([]Users, error)
	GetById    func(id int) (Users, error)
	Save       func(obj Users) error
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewUsersController(db *gorm.DB, baseLog zerolog.Logger) UsersController {
	fpController := fp.CreateController()
	log := baseLog.With().Str("model", "users").Logger()
	if err := db.AutoMigrate(&Users{}); err != nil {
		log.Fatal().Err(err).Msg("auto-migrate")
	}

	return UsersController{
		CreateHash: func() (string, error) {
			var objs []Users
			var obj Users
			var generatedHash string

			tx := db.Model(&Users{}).Find(&objs)
			generatedHash = tools.RandStringBytes(24)
			generatedHash += letterBytes[len(objs)/52 : len(objs)/52+1]
			obj.GeneratedHash = generatedHash

			tx = db.Model(&Users{}).Create(&obj)
			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				return "", tx.Error
			}
			return generatedHash, nil
		},
		UpdateIP: func(IP, hashString, UA, city string) error {
			var obj Users
			var objCheck Users
			tx := db.Model(&Users{}).Find(&objCheck, "generated_hash", hashString)
			if tx.RowsAffected == 0 {
				err := errors.New("hash not found")
				log.Error().Err(err).Msg("hash not found")
				return err
			}

			tx = db.Model(&Users{}).Find(&obj, "uniq_hash", hashString+IP)

			obj.IPLocation = city
			if objCheck.ID > 0 {
				obj.TelegramName = objCheck.TelegramName
			}

			if obj.IP == "" {
				obj.FP = UA
				obj.IP = IP
				obj.GeneratedHash = hashString
				obj.UniqHash = hashString + IP
				obj.Count = 1
				obj.IsAntik = fpController.Check(UA)
				tx = db.Save(&obj)
				if tx.Error != nil {
					log.Error().Err(tx.Error).Msg("db update ip new error")
					return tx.Error
				}
				return nil
			}

			if obj.IP == IP {
				obj.GeneratedHash = hashString
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
				obj.IsAntik = fpController.Check(UA)
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
			tx := db.Model(&Users{}).Order("generated_hash").Find(&obj)

			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				return nil, tx.Error
			}
			return obj, nil
		},
		GetById: func(id int) (Users, error) {
			var obj Users
			tx := db.Model(&Users{}).Where("id = ?", id).Find(&obj)

			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				return obj, tx.Error
			}
			return obj, nil
		},
		Save: func(obj Users) error {
			tx := db.Save(&obj)
			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db update ip error")
				return tx.Error
			}
			return nil
		},
	}
}
