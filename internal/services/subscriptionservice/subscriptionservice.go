package subscriptionservice

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Subscription struct {
	bun.BaseModel `bun:"table:subscription_subscription,alias:sub"`
	ID            int64     `json:"id" bun:",pk,autoincrement"`
	ChatID        int64     `json:"chatId" bun:",unique"`
	Threshold     uint32    `json:"threshold"`
	CreatedAt     time.Time `json:"createdAt" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updatedAt" bun:",nullzero,notnull,default:current_timestamp"`
}

func newSubscription(chatID int64, v uint32) Subscription {
	return Subscription{
		ChatID:    chatID,
		Threshold: v,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type SubscriptionService struct {
	db *bun.DB
}

func NewSubscriptionService(db *bun.DB) SubscriptionService {
	return SubscriptionService{
		db: db,
	}
}

func (ss *SubscriptionService) Subscribe(chatID int64, limit uint32) (*Subscription, error) {
	ns := newSubscription(chatID, limit)
	s := &ns
	ctx := context.TODO()

	_, err := ss.db.NewInsert().Model(s).On("CONFLICT (chat_id) DO UPDATE").Set("threshold = EXCLUDED.threshold").Set("updated_at = current_timestamp").Exec(ctx)
	return s, err
}
