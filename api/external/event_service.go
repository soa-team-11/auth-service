package external

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

var ctx = context.Background()

type EventService struct {
	rdb *redis.Client
}

func NewEventService() *EventService {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return &EventService{rdb: rdb}
}

// Zahtev za shopping cart
func (es *EventService) PublishUserRegistered(ctx context.Context, userId string) {
	tracer := otel.Tracer("auth-service")
	_, span := tracer.Start(ctx, "EventService.PublishUserRegistered")
	defer span.End()

	event := map[string]string{"userId": userId}
	data, _ := json.Marshal(event)

	if err := es.rdb.Publish(ctx, "user-registered", data).Err(); err != nil {
		span.RecordError(err)
		log.Println("Failed to publish user-registered event:", err)
	} else {
		log.Println("Published user-registered event for user:", userId)
	}

	span.End()
}

// Kompenzacija
func (es *EventService) SubscribeCartCreationFailures(deleteUserFunc func(userID string) error) {
	pubsub := es.rdb.Subscribe(ctx, "cart-creation-failed")
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var data map[string]string
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				log.Println("Failed to parse failure event:", err)
				continue
			}

			userID := data["userId"]
			log.Println("Rolling back user:", userID)

			if err := deleteUserFunc(userID); err != nil {
				log.Println("Failed to delete user:", err)
			} else {
				log.Println("User deleted due to cart creation failure:", userID)
			}
		}
	}()
}
