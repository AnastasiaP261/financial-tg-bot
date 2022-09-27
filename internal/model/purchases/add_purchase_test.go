package purchases_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	mocks "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases/_mocks"
)

func Test_AddPurchase_OnlySum(t *testing.T) {
	t.Run("целое число", func(t *testing.T) {
		repo := mocks.NewMockRepo(gomock.NewController(t))
		model := purchases.New(repo)

		repo.EXPECT().AddPurchase(gomock.Any()).Return(nil)

		err := model.AddPurchase(123, "123", "", "")
		assert.NoError(t, err)
	})

	t.Run("дробное число", func(t *testing.T) {
		repo := mocks.NewMockRepo(gomock.NewController(t))
		model := purchases.New(repo)

		repo.EXPECT().AddPurchase(gomock.Any()).Return(nil)

		err := model.AddPurchase(123, "234.5", "", "")
		assert.NoError(t, err)
	})

	t.Run("невалидное число", func(t *testing.T) {
		repo := mocks.NewMockRepo(gomock.NewController(t))
		model := purchases.New(repo)

		err := model.AddPurchase(123, "12o.o5", "", "")
		assert.Error(t, err, purchases.ErrSummaParsing)
	})
}

func Test_AddPurchase_SumAndCategory(t *testing.T) {
	t.Run("добавление траты по уже существующей категории", func(t *testing.T) {
		repo := mocks.NewMockRepo(gomock.NewController(t))
		model := purchases.New(repo)

		repo.EXPECT().CategoryExist(gomock.Any()).Return(true)
		repo.EXPECT().AddPurchase(gomock.Any()).Return(nil)

		err := model.AddPurchase(123, "234.5", "some category", "")
		assert.NoError(t, err)
	})

	t.Run("добавление траты по не существующей категории", func(t *testing.T) {
		repo := mocks.NewMockRepo(gomock.NewController(t))
		model := purchases.New(repo)

		repo.EXPECT().CategoryExist(gomock.Any()).Return(false)

		err := model.AddPurchase(123, "234.5", "some category", "")
		assert.Error(t, err, purchases.ErrCategoryNotExist)
	})
}

func Test_AddPurchase_SumAndCategoryAndDate(t *testing.T) {
	t.Run("добавление с валидной датой", func(t *testing.T) {
		repo := mocks.NewMockRepo(gomock.NewController(t))
		model := purchases.New(repo)

		repo.EXPECT().CategoryExist(gomock.Any()).Return(true)
		repo.EXPECT().AddPurchase(gomock.Any()).Return(nil)

		err := model.AddPurchase(123, "234.5", "some category", "01.01.2022")
		assert.NoError(t, err)
	})

	t.Run("добавление с не валидной датой", func(t *testing.T) {
		repo := mocks.NewMockRepo(gomock.NewController(t))
		model := purchases.New(repo)

		repo.EXPECT().CategoryExist(gomock.Any()).Return(true)

		err := model.AddPurchase(123, "234.5", "some category", "01-01-2022")
		assert.Error(t, err, purchases.ErrDateParsing)
	})
}
