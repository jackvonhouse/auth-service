package validator

import (
	"github.com/beevik/guid"
)

func IsValidGUID(g string) bool { return guid.IsGuid(g) }
