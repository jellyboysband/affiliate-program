package store

import (
	"context"
	"github.com/jellyboysband/eye/cmd/store/internal/app"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	dbName         = "eye"
	collectionName = "products"
)

type (
	store struct {
		db *mongo.Client
	}

	product struct {
		Id            primitive.ObjectID `bson:"_id"`
		AliID         int                `bson:"ali_id"`
		Rating        float64            `bson:"rating"`
		URL           string             `bson:"url"`
		Title         string             `bson:"title"`
		TotalSales    int                `bson:"total_sales"`
		RatingProduct string             `bson:"rating_product"`
		TotalComment  int                `bson:"total_comment"`
		Discount      float64            `bson:"discount"`
		Max           price              `bson:"max"`
		Min           price              `bson:"min"`
		Shop          shop               `bson:"shop"`
	}

	price struct {
		Currency string  `bson:"currency"`
		Cost     float64 `bson:"cost"`
	}

	shop struct {
		ID           int    `bson:"id"`
		Name         string `bson:"name"`
		Followers    int    `bson:"followers"`
		PositiveRate string `bson:"positive_rate"`
	}
)

func New(db *mongo.Client) app.Store {
	return &store{db: db}
}

func (s *store) Save(ctx context.Context, product app.ArgSaveProduct) (*app.Product, error) {
	collection := s.db.Database(dbName).Collection(collectionName)

	document := convertToMONGO(product)
	document.Id = primitive.NewObjectIDFromTimestamp(time.Now())
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert document")
	}

	return s.Product(ctx, document.Id.String())
}

func (s *store) Product(ctx context.Context, id string) (*app.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}

	filter := bson.M{"_id": objectID}
	document := &product{}

	collection := s.db.Database(dbName).Collection(collectionName)
	err = collection.FindOne(ctx, filter).Decode(document)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document")
	}

	return convertToAppFormat(document), nil
}

func (s *store) Products(ctx context.Context) ([]app.Product, error) {
	collection := s.db.Database(dbName).Collection(collectionName)

	filter := bson.M{}
	c, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get products")
	}

	var products []product
	err = c.All(ctx, &products)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode documents")
	}

	return convertList(products), nil
}

func convertList(products []product) []app.Product {
	convertsProducts := make([]app.Product, len(products))
	for i := range products {
		convertsProducts[i] = *convertToAppFormat(&products[i])
	}
	return convertsProducts
}

func convertToAppFormat(document *product) *app.Product {
	return &app.Product{
		Id:            document.Id.String(),
		AliID:         document.AliID,
		URL:           document.URL,
		Title:         document.Title,
		TotalSales:    document.TotalSales,
		RatingProduct: document.RatingProduct,
		TotalComment:  document.TotalComment,
		Discount:      document.Discount,
		Max: app.Price{
			Currency: document.Max.Currency,
			Cost:     document.Max.Cost,
		},
		Min: app.Price{
			Currency: document.Min.Currency,
			Cost:     document.Min.Cost,
		},
		Shop: app.Shop{
			ID:           document.Shop.ID,
			Name:         document.Shop.Name,
			Followers:    document.Shop.Followers,
			PositiveRate: document.Shop.PositiveRate,
		},
	}
}

func convertToMONGO(appFormat app.ArgSaveProduct) product {
	return product{
		AliID:         appFormat.AliID,
		Rating:        appFormat.Rating,
		URL:           appFormat.URL,
		Title:         appFormat.Title,
		TotalSales:    appFormat.TotalSales,
		RatingProduct: appFormat.RatingProduct,
		TotalComment:  appFormat.TotalComment,
		Discount:      appFormat.Discount,
		Max: price{
			Currency: appFormat.Max.Currency,
			Cost:     appFormat.Max.Cost,
		},
		Min: price{
			Currency: appFormat.Min.Currency,
			Cost:     appFormat.Min.Cost,
		},
		Shop: shop{
			ID:           appFormat.Shop.ID,
			Name:         appFormat.Shop.Name,
			Followers:    appFormat.Shop.Followers,
			PositiveRate: appFormat.Shop.PositiveRate,
		},
	}
}
