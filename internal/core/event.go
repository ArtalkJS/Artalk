package core

// -------------------------------------------------------------------
// Event data
// -------------------------------------------------------------------

type BootstrapEvent struct {
	App *App
}

type TerminateEvent struct {
	App *App
}
