package api

import (
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api/bright_horizons"
	"tadpoles-backup/internal/api/tadpoles"
	"tadpoles-backup/internal/schemas"
)

func GetProvider() schemas.Provider {
	switch config.Provider.String() {
	case config.BrightHorizons:
		return bright_horizons.NewProvider()
	default:
		return tadpoles.NewProvider()
	}
}
