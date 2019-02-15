// +build android ios

package tge

import (
	log "log"
)

// Run main entry point of runtime
func Run(app App) error {
	log.Println("Run()")

	// -------------------------------------------------------------------- //
	// Create
	// -------------------------------------------------------------------- //
	settings := &defaultSettings
	err := app.OnCreate(settings)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
