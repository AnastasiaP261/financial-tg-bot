package purchases

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
)

var (
	ErrSummaParsing        = errors.New("summa parsing error")
	ErrLimitParsing        = errors.New("limit parsing error")
	ErrDateParsing         = errors.New("date parsing error")
	ErrCategoryNotExist    = errors.New("such category doesn't exist")
	ErrUnknownPeriod       = errors.New("unknown period")
	ErrInvalidDate         = errors.New("invalid date")
	ErrUserHasntCategory   = errors.New("this user hasn't such category")
	ErrCreateReportRequest = errors.New("create report request failed")
)

// Repo репозиторий
type Repo interface {
	GetRate(ctx context.Context, y int, m int, d int) (bool, currency.RateToRUB, error)
	AddRate(ctx context.Context, y int, m int, d int, rates currency.RateToRUB) error

	UserCreateIfNotExist(ctx context.Context, userID int64) error
	ChangeCurrency(ctx context.Context, userID int64, currency currency.Currency) error
	GetUserInfo(ctx context.Context, userID int64) (User, error)
	ChangeUserLimit(ctx context.Context, userID int64, newLimit float64) error
	AddCategoryToUser(ctx context.Context, userID int64, catName string) error
	UserHasCategory(ctx context.Context, userID int64, categoryID uint64) (bool, error)
	GetUserCategories(ctx context.Context, userID int64) ([]string, error)

	AddPurchase(ctx context.Context, req AddPurchaseReq) error
	GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]Purchase, error)
	GetUserPurchasesSumFromMonth(ctx context.Context, userID int64, fromDate time.Time) (float64, error)

	GetCategoryID(ctx context.Context, categoryName string) (uint64, error)
	AddCategory(ctx context.Context, categoryName string) error
	GetAllCategories(ctx context.Context) ([]CategoryRow, error)
}

type ExchangeRateGetter interface {
	GetExchangeRateToRUB() currency.RateToRUB
	GetExchangeRateToRUBFromDate(ctx context.Context, y, m, d int) (currency.RateToRUB, error)
}

type ReportsStore interface {
	Delete(ctx context.Context, key string) error
}

type BrokerMsgCreator interface {
	SendNewMsg(key string, value string) error
}

type Model struct {
	Repo               Repo
	ExchangeRatesModel ExchangeRateGetter
	ReportsStore       ReportsStore
	BrokerMsgCreator   BrokerMsgCreator
}

func New(repo Repo, exchangeRatesModel ExchangeRateGetter, reportsStore ReportsStore, producer BrokerMsgCreator) *Model {
	return &Model{
		Repo:               repo,
		ExchangeRatesModel: exchangeRatesModel,
		ReportsStore:       reportsStore,
		BrokerMsgCreator:   producer,
	}
}
