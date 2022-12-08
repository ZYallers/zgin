package route

import (
	v000 "github.com/ZYallers/zgin/example/controller/v000"
	"github.com/ZYallers/zgin/helper/restful"
	"github.com/ZYallers/zgin/types"
)

var Restful types.Restful

func init() {
	Restful = restful.Register(Restful, &v000.Person{})
}

