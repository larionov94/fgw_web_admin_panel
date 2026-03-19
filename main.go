package main

import (
	"fgw_web_admin_panel/pkg/logg"
)

const skipNofS = 4 // skipNofS кол-во пропускаемых кадров стека.

func main() {
	logger, _ := logg.NewLogger()
	defer logger.Close()

	logger.LogI("I3000 Это только начало.", skipNofS)
}
