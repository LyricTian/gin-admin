package structure

import (
	"github.com/jinzhu/copier"
)

func Copy(s, ts interface{}) error {
	return copier.Copy(ts, s)
}
