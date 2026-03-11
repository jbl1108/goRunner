package main

import "github.com/jbl1108/goRunner/config"

func main() {
	app, err := config.NewApplication()
	if err != nil {
		panic(err)
	}
	// Start the REST service in a separate goroutine
	go func() {
		err := app.RestService.Start()
		if err != nil {
			panic(err)
		}
	}()
	// Start the MQTT client (this will block)
	app.OutputPublisher.Connect()
	defer app.OutputPublisher.Disconnect()
	// Block forever
	select {}
}
