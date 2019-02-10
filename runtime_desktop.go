// +build darwin freebsd linux windows
// +build !android
// +build !ios

package tge

import log "log"

func backend_Instanciate(app App) error {
	log.Println("Creation desktop runtime")
	return nil
}
