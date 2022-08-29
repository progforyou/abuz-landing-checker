package axdata

import (
	"axp/target/generated-sources/proto/dominous"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Token    string `gorm:"type:VARBINARY(32);index:idx_uniq_user_session,unique;column:token"`
	Name     string
	Level    int
	Rank     int
	Exp      int64
	Ico      int
	Bot      bool    `gorm:"index"`
	Room     *uint64 `gorm:"index"`
	Version  string
	Platform dominous.PUserPlatform `gorm:""`
	Social   UserSocial             `gorm:"embedded"`
}

type UserSocial struct {
	SocialId   string
	SocialType dominous.PSocialType `gorm:""`
}

//func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
//	user.token = user.Token.String()
//	user.session = user.Session.String()
//	return
//}
