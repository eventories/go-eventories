package core

import (
	"github.com/eventories/go-eventories/core/interaction"
	"github.com/eventories/go-eventories/database"
)

type Worker struct {
	eth *interaction.Interactor

	db         database.Database
	checkpoint uint64
}
