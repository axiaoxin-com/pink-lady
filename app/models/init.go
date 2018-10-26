// package models save database models
// define your database models in the package
package models

import (
	"github.com/jinzhu/gorm"
)

type BaseModel gorm.Model

// init function here for package models
func init() {

}
