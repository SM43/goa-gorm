package userapi
import (
	"github.com/jinzhu/gorm"
)

type (

	//User model represents categories which associated with group of tags
	User struct {
		gorm.Model
		Name string `gorm:"size:100;not null;unique"`
	}

)