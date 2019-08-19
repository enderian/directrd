package terminals

import (
	"log"

	"github.com/enderian/directrd/pkg/types"
)

func loadTerminals() {
	var terminals []*types.Terminal
	if err := ctx.DB().Find(&terminals).Error; err != nil {
		log.Fatalf("failed to load terminals: %v", err)
		return
	}
	log.Printf("loaded %d terminals from the database.", len(terminals))
}