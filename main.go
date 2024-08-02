package main

import (
	"github.com/IkezawaYuki/c_root/interface/router"
)

func main() {
	_ = router.NewRouter().Run(":8082")
}
