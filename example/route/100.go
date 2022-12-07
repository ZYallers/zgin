package route

import (
	v100 "github.com/ZYallers/zgin/example/controller/v100"
	"github.com/ZYallers/zgin/helper/restful"
)

func init() {
	Restful = restful.Register(Restful, &v100.Person{})
}
