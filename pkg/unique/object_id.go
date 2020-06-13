package unique

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectID Define alias
type ObjectID = primitive.ObjectID

// NewObjectID Create object id
func NewObjectID() ObjectID {
	return primitive.NewObjectID()
}
