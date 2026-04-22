// A simple BubbleTea program that runs both natively and in the browser via WASM.
package main

import (
	"log"

	booba "github.com/NimbleMarkets/go-booba"
	"github.com/NimbleMarkets/go-booba-example/internal/model"
)

func main() {
	if err := booba.Run(model.InitialModel()); err != nil {
		log.Fatal(err)
	}
}
