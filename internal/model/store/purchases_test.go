package store

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func Test_AddPurchase(t *testing.T) {
	s := New()
	s.Users = []userRow{
		{sync.RWMutex{}, user{UserID: 123, Currency: RUB}},
	}

	date, _ := time.Parse("02.01.2006", "01.01.2000")
	err := s.AddPurchase(model.AddPurchaseReq{
		UserID:   123,
		Sum:      100.50,
		Category: "some category",
		Date:     date,
	})

	assert.NoError(t, err)
	assert.Equal(t,
		purchase{
			UserID:   123,
			Sum:      100.50,
			Category: "some category",
			Date:     date,
		},
		s.purchaseAccessRead()[0],
	)
}

func Test_GetUserPurchasesFromDate(t *testing.T) {
	s := New()

	t1, _ := time.Parse("02.01.2006", "01.10.2022")
	t2, _ := time.Parse("02.01.2006", "02.10.2022")
	t3, _ := time.Parse("02.01.2006", "03.10.2022")
	t4, _ := time.Parse("02.01.2006", "04.10.2022")

	s.Users = []userRow{
		{sync.RWMutex{}, user{UserID: 123, Currency: RUB}},
	}
	s.Purchases = []purchaseRow{
		{sync.RWMutex{}, purchase{UserID: 123, Sum: 100, Category: "some category 1", Date: t1, USDRatio: 2, CNYRatio: 2, EURRatio: 2}},
		{sync.RWMutex{}, purchase{UserID: 123, Sum: 100, Category: "some category 2", Date: t2, USDRatio: 2, CNYRatio: 2, EURRatio: 2}},
		{sync.RWMutex{}, purchase{UserID: 123, Sum: 200, Category: "some category 1", Date: t3, USDRatio: 2, CNYRatio: 2, EURRatio: 2}},
		{sync.RWMutex{}, purchase{UserID: 123, Sum: 100, Category: "some category 1", Date: t4, USDRatio: 2, CNYRatio: 2, EURRatio: 2}},
	}

	res, err := s.GetUserPurchasesFromDate(t2, 123)

	assert.NoError(t, err)
	assert.EqualValues(t, []model.Purchase{
		{PurchaseCategory: "some category 2", Summa: 100, RateToRUB: model.RateToRUB{USD: 2, EUR: 2, CNY: 2}},
		{PurchaseCategory: "some category 1", Summa: 200, RateToRUB: model.RateToRUB{USD: 2, EUR: 2, CNY: 2}},
		{PurchaseCategory: "some category 1", Summa: 100, RateToRUB: model.RateToRUB{USD: 2, EUR: 2, CNY: 2}},
	}, res)
}
