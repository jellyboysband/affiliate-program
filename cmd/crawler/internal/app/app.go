package app

type (
	App interface {
		Save(d Document) error
	}
	Store interface {
		Send(d Document) error
	}
	app struct {
		store Store
	}
)

func New(s Store) App {
	return &app{store: s}
}

// TODO generate affileta URL
func (app *app) Save(d Document) error {
	return app.store.Send(d)
}
