package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"sync"
	"sync/atomic"
	"time"
)

const (
	PublishUserLikedPlaces = "liked_place_vec"
	PublishMetaData        = "metadata_vec"
	PublishLikedPlaceLists = "liked_place_lists_vec"
)

type QueueHandler interface {
	PublishUserVec(ctx context.Context, userID int, messageType string) error
}

type PublishID struct {
	ID int `json:"id"`
}

type PlacePublish struct {
	PlaceHeaderID int    `json:"place_header_id"`
	Description   string `json:"description"`
	TagIDs        []int  `json:"tag_ids"`
	FeatureIDs    []int  `json:"feature_ids"`
}

type rabbitQueueHandler struct {
	conn             *amqp.Connection
	reconnectRetries atomic.Int32
	reconnectFunc    func()
}

func NewQueueHandler() QueueHandler {
	handler := &rabbitQueueHandler{reconnectFunc: func() {}}
	handler.connect()
	return handler
}

func (r *rabbitQueueHandler) connect() {
	connAddr := fmt.Sprintf("amqp://%v:%v@%v:%v/",
		viper.GetString("RABBIT_USER"), viper.GetString("RABBIT_PASS"),
		viper.GetString("RABBIT_HOST"), viper.GetString("RABBIT_PORT"))

	var once sync.Once
	reconnectFunc := func() {
		time.Sleep(time.Second)
		if r.conn == nil || r.conn.IsClosed() {
			if r.reconnectRetries.Load() == viper.GetInt32("RABBIT_MAX_RETRIES") {
				panic("rabbit max retries exceeded")
			}
			r.reconnectRetries.Add(1)
			once.Do(r.connect)
		}
	}

	var err error
	r.conn, err = amqp.Dial(connAddr)
	if err != nil {
		fmt.Println(fmt.Sprintf("error on connecting to rabbit: %v, retrying...", err))
		reconnectFunc()
	}

	r.reconnectFunc = reconnectFunc

	r.reconnectRetries.Store(0)
}

func (r *rabbitQueueHandler) PublishUserVec(ctx context.Context, userID int, messageType string) error {
	body, err := json.Marshal(PublishID{userID})
	if err != nil {
		panic(fmt.Sprintf("wtf just happened in rabbit.go func publish user vec?"))
	}

	if messageType != PublishMetaData && messageType != PublishUserLikedPlaces && messageType != PublishLikedPlaceLists {
		return errors.New("wrong message type dude")
	}

	ch, err := r.conn.Channel()
	if err != nil {
		r.reconnectFunc()
		for {
			if !r.conn.IsClosed() {
				break
			}
		}
	}

	err = ch.PublishWithContext(ctx,
		"",
		viper.GetString("USER_VEC_QUEUE"),
		false,
		false,
		amqp.Publishing{
			Headers:   map[string]interface{}{"message_type": messageType},
			Body:      body,
			Timestamp: time.Now(),
		})

	return err
}
