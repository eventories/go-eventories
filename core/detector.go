package core

import (
	"github.com/eventories/go-eventories/core/interaction"
	"github.com/eventories/go-eventories/core/table"
	"github.com/eventories/go-eventories/database"
)

type Detector struct {
	eth *interaction.Interactor
	tab *table.Table

	db         database.Database
	checkpoint uint64
}
