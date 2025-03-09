package main

import (
	"context"
	"fmt"

	"loyalit/internal/adapters/app"

	_ "loyalit/docs"
	_ "time/tzdata"
)

//	@title			LoyalT API
//	@version		1.0
//	@description	API of LoyalT - platform for creating and managing loyalty programs
//	@termsOfService	http://swagger.io/terms/

// @host		prod-team-22-t62v97db.final.prodcontest.ru
// @BasePath	/api/v1
// @schemes	https
func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to init app: %v", err))
	}

	a.Run()
}
