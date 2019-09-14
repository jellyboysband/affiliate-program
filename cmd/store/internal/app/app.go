package app

import "context"

type (
	Store interface {
		Save(ctx context.Context, productInfo ArgSaveProduct) (*Product, error)
		Products(ctx context.Context) ([]Product, error)
	}

	App interface {
		Save(ctx context.Context, productInfo ArgSaveProduct) (*Product, error)
		ListProduct(ctx context.Context) ([]Product, error)
	}

	app struct {
		s Store
	}
)

func New(s Store) App {
	return &app{s}
}

func (app *app) Save(ctx context.Context, productInfo ArgSaveProduct) (*Product, error) {
	return app.s.Save(ctx, productInfo)
}

func (app *app) ListProduct(ctx context.Context) ([]Product, error) {
	return app.s.Products(ctx)
}
