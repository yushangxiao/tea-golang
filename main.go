package main

import (
	"tea/router"
)

func main() {
	r := router.SetupRouter()
	r.Run("0.0.0.0:7860") // listen and serve on 0.0.0.0:8080
}
