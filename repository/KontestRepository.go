// repository.go
package repository

import (
	"kontest-api/model"
)

// KontestRepository defines methods for contest data operations.
type KontestRepository interface {
	Save(kontest model.KontestModel)
	DeleteAll()
}
