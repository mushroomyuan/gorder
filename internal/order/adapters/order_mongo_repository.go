package adapters

import (
	"context"
	"time"

	domain "github.com/mushroomyuan/gorder/order/domain/order"
	"github.com/mushroomyuan/gorder/order/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dbName   = viper.GetString("monogo.db-name")
	collName = viper.GetString("monogo.coll-name")
)

type OrderRepositoryMongo struct {
	db *mongo.Client
}

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{db: db}
}

func (m *OrderRepositoryMongo) collection() *mongo.Collection {
	return m.db.Database(dbName).Collection(collName)
}

type orderModel struct {
	MonogoID    primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `bson:"customer_id"`
	Status      string             `bson:"status"`
	PaymentLink string             `bson:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (m *OrderRepositoryMongo) marshalToModel(order *domain.Order) *orderModel {
	return &orderModel{
		MonogoID:    primitive.NewObjectID(),
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}

func (*OrderRepositoryMongo) logWithTag(tag string, err error, result any) {
	l := logrus.WithFields(logrus.Fields{
		"tag":            "order_repository_mongo",
		"performed_time": time.Now().Unix(),
		"result":         result,
		"err":            err,
	})
	if err != nil {
		l.Infof("%s_fail", tag)
	} else {
		l.Infof("%s_success", tag)
	}
}

func (m *OrderRepositoryMongo) unmarshal(o *orderModel) *domain.Order {
	return &domain.Order{
		ID:          o.MonogoID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       o.Items,
	}
}

func (m *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	defer m.logWithTag("create", err, created)
	write := m.marshalToModel(order)
	res, err := m.collection().InsertOne(ctx, write)
	if err != nil {
		return nil, err
	}
	created = order
	order.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return order, nil
}

func (m *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	defer m.logWithTag("get", err, got)
	read := &orderModel{}
	mongoID, _ := primitive.ObjectIDFromHex(id)
	cond := bson.M{"_id": mongoID}
	if err = m.collection().FindOne(ctx, cond).Decode(read); err != nil {
		return
	}
	if len(read.MonogoID) == 0 {
		return nil, domain.NotFoundError{OrderID: id}
	}
	return m.unmarshal(read), nil
}

func (m *OrderRepositoryMongo) Update(
	ctx context.Context,
	o *domain.Order,
	updateFn func(context.Context, *domain.Order,
	) (*domain.Order, error)) (err error) {
	defer m.logWithTag("update", err, nil)
	if o == nil {
		panic("got nil order")
	}
	session, err := m.db.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(ctx)
	if err = session.StartTransaction(); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = session.AbortTransaction(ctx)
		} else {
			_ = session.CommitTransaction(ctx)
		}
	}()

	// inside transaction
	oldOrder, err := m.Get(ctx, o.ID, o.CustomerID)
	if err != nil {
		return
	}
	updated, err := updateFn(ctx, o)
	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
	res, err := m.collection().UpdateOne(
		ctx,
		bson.M{"_id": mongoID, "customer_id": o.CustomerID},
		bson.M{"$set": bson.M{
			"status":       updated.Status,
			"payment_link": updated.PaymentLink,
		}},
	)
	if err != nil {
		return
	}
	m.logWithTag("finish_update", err, res)
	return

}
