package tick

import "github.com/hootuu/eggcone/fdn/tick/def"

type Listener interface {
	GetName() string
	Match(job *def.Job) bool
	Deal(job *def.Job) (any, error)
}
