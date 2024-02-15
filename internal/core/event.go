package core

import "github.com/ArtalkJS/Artalk/internal/config"

// -------------------------------------------------------------------
// Event data
// -------------------------------------------------------------------

type BootstrapEvent struct {
	App *App
}

type TerminateEvent struct {
	App *App
}

type ConfUpdatedEvent struct {
	App  *App
	Conf *config.Config
}
