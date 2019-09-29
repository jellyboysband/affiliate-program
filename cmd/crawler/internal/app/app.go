package app

type (
	App interface {
		Save(d Document) error
	}
	Store interface {
		Send(d Document, rate float64) error
	}
	app struct {
		store      Store
		ourRateMin float64
	}
)

func New(s Store, ourRateMin float64) App {
	return &app{store: s, ourRateMin: ourRateMin}
}

// TODO generate affileta URL
func (app *app) Save(d Document) error {
	ourRate := d.Rating()

	switch ourRate >= app.ourRateMin {
	case true:
		return app.store.Send(d, ourRate)
	default:
		return nil // Dont send.
	}
}
