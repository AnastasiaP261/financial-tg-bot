package messages

import (
	"context"

	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type MessageSender interface {
	SendMessage(text string, userID int64, userName string) error
	SendImage(img []byte, chatID int64, userName string) error
	SendKeyboard(text string, chatID int64, buttonTexts []string, userName string) error
}

type PurchasesModel interface {
	AddPurchase(ctx context.Context, userID int64, rawSum, category, rawDate string) (purchases.ExpensesAndLimit, error)

	AddCategory(ctx context.Context, category string) error
	GetAllCategories(ctx context.Context) ([]purchases.CategoryRow, error)

	Report(ctx context.Context, period purchases.Period, userID int64) (txt string, img []byte, err error)

	ChangeUserCurrency(ctx context.Context, userID int64, currency purchases.Currency) error
	ChangeUserLimit(ctx context.Context, userID int64, rawLimit string) error
	AddCategoryToUser(ctx context.Context, userID int64, category string) error
	GetUserCategories(ctx context.Context, userID int64) ([]string, error)

	ToPeriod(str string) (purchases.Period, error)
	StrToCurrency(str string) (purchases.Currency, error)
	CurrencyToStr(cy purchases.Currency) (string, error)
}

type StatusStore interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type Model struct {
	tgClient       MessageSender
	purchasesModel PurchasesModel
	statusStore    StatusStore
}

func New(tgClient MessageSender, purchasesModel PurchasesModel, redis StatusStore) *Model {
	return &Model{
		tgClient:       tgClient,
		purchasesModel: purchasesModel,
		statusStore:    redis,
	}
}