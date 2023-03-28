package nextdrawservice

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type NextDrawService struct {
	db *bun.DB
}

func NewNextDrawService(db *bun.DB) NextDrawService {
	return NextDrawService{
		db: db,
	}
}

type NextDraw struct {
	bun.BaseModel `bun:"table:nextdraw_nextdraw,alias:nd"`
	ID            int64  `json:"id" bun:",pk,autoincrement"`
	DateString    string `bun:",unique"`
	Prize         string
	CreatedAt     time.Time `json:"createdAt" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updatedAt" bun:",nullzero,notnull,default:current_timestamp"`
}

type INextDraw interface {
	GetDate() string
	GetPrize() string
}

func (s NextDrawService) Save(n INextDraw) (INextDraw, error) {
	d := NextDraw{
		DateString: n.GetDate(),
		Prize:      n.GetPrize(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	ctx := context.TODO()

	_, err := s.db.NewInsert().Model(&d).Exec(ctx)
	return n, err

}
