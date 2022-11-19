package main

import "context"

func main() {
	ctx := context.Background()

	app, err := initialize()
	if err != nil {
		panic(err)
	}

	app.logger.Info(ctx, "Hello, World!", app.logger.Field("config", app.config.Config()))
}
