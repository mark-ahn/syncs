package syncs

//go:generate genny -in rc__template.go -out rc__template__gen.go gen "_Prefix_=Of Some=interface{}"

import (
	"github.com/cheekybits/genny/generic"
)

type Some generic.Type
