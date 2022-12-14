package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type user struct {
	UserID      int64         `db:"id"`
	Currency    Currency      `db:"curr"` // выбранная пользователем валюта
	CategoryIDs pq.Int64Array `db:"category_ids"`
	Limit       float64       `db:"month_limit"`
}

// Currency тип валюты
type Currency string

const (
	// RUB валюта - рубль
	RUB Currency = "RUB"

	// USD валюта - доллар
	USD Currency = "USD"

	// EUR валюта - евро
	EUR Currency = "EUR"

	// CNY валюта - китайский юань
	CNY Currency = "CNY"
)

func currencyToModelTypeConv(c Currency) (currency.Currency, error) {
	switch c {
	case RUB:
		return currency.RUB, nil
	case USD:
		return currency.USD, nil
	case EUR:
		return currency.EUR, nil
	case CNY:
		return currency.CNY, nil
	default:
		return 0, errors.New("invalid currency")
	}
}
func currencyFromModelTypeConv(c currency.Currency) (Currency, error) {
	switch c {
	case currency.RUB:
		return RUB, nil
	case currency.USD:
		return USD, nil
	case currency.EUR:
		return EUR, nil
	case currency.CNY:
		return CNY, nil
	default:
		return "", errors.New("invalid currency")
	}
}

// UserCreateIfNotExist проверяет, что такой юзер есть в базе, и, если его нет, создает такого юзера.
// Нужно вызывать эту функцию в начале каждой другой команды. Это позволит лениво создать запись о пользователе и
// снимет с модели ответственность за нормализацию данных
func (s *Service) UserCreateIfNotExist(ctx context.Context, userID int64) error {
	ok, err := s.userExist(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "userExist")
	}

	if !ok {
		if err = s.addUser(ctx, userID); err != nil {
			return errors.Wrap(err, "addUser")
		}
	}

	return nil
}

// userExist проверка, что такой юзер уже создан в базе
func (s *Service) userExist(ctx context.Context, userID int64) (bool, error) {
	res, err := s.getUserInfo(ctx, userID)
	if errors.Is(err, ErrUserDoesntExists) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "getUserInfo")
	}

	if res.UserID == 0 {
		return false, nil
	}

	return true, nil
}

// addUser добавляет юзера с такой айдишкой в базу
func (s *Service) addUser(ctx context.Context, userID int64) error {
	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tblUsers).
		Columns(tblUsersColID).
		Values(userID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// ChangeCurrency смена валюты пользователя
func (s *Service) ChangeCurrency(ctx context.Context, userID int64, currency currency.Currency) error {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}

	curr, err := currencyFromModelTypeConv(currency)
	if err != nil {
		return errors.Wrap(err, "currencyFromModelTypeConv")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tblUsers).
		Set(tblUsersColCurrency, curr).
		Where(sq.Eq{tblUsersColID: userID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// GetUserInfo возвращает информацию о пользователе в формате модели
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (model.User, error) {
	res, err := s.getUserInfo(ctx, userID)
	if err != nil && errors.Is(err, ErrUserDoesntExists) {

		if err := s.addUser(ctx, userID); err != nil {
			return model.User{}, errors.Wrap(err, "addUser")
		}

		res, err = s.getUserInfo(ctx, userID)
		if err != nil {
			return model.User{}, errors.Wrap(err, "getUserInfo")
		}

	} else if err != nil {
		return model.User{}, errors.Wrap(err, "getUserInfo")
	}

	curr, err := currencyToModelTypeConv(res.Currency)
	if err != nil {
		return model.User{}, errors.Wrap(err, "currencyToModelTypeConv")
	}

	return model.User{
		UserID:     res.UserID,
		Currency:   curr,
		Categories: res.CategoryIDs,
		Limit:      res.Limit,
	}, nil
}

// getUserInfo возвращает информацию о пользователе (для использования внутри пакета)
func (s *Service) getUserInfo(ctx context.Context, userID int64) (user, error) {
	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblUsersColID, tblUsersColCurrency, tblUsersColLimit, tblUsersColCategoriesIDs).
		From(tblUsers).
		Where(sq.Eq{
			tblUsersColID: userID,
		}).
		ToSql()
	if err != nil {
		return user{}, errors.Wrap(err, "query creating error")
	}

	rows, err := s.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return user{}, errors.Wrap(err, "db.QueryxContext")
	}

	data := user{}
	if err = readX(rows, &data); err != nil {
		return user{}, errors.Wrap(err, "readX")
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) || data.UserID == 0 {
		return user{}, ErrUserDoesntExists
	} else if err != nil {
		return user{}, errors.Wrap(err, "db.QueryxContext")
	}

	return data, nil
}

func (s *Service) ChangeUserLimit(ctx context.Context, userID int64, newLimit float64) error {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tblUsers).
		Set(tblUsersColLimit, newLimit).
		Where(sq.Eq{tblUsersColID: userID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

func (s *Service) AddCategoryToUser(ctx context.Context, userID int64, catName string) error {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}
	catID, err := s.GetCategoryID(ctx, catName)
	if err != nil {
		return errors.Wrap(err, "GetCategoryID")
	}

	q, args, err := sq.Expr(`UPDATE users 
							SET category_ids=array_append(category_ids, $1)
							WHERE id=$2`,
		catID, userID).ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

func (s *Service) UserHasCategory(ctx context.Context, userID int64, categoryID uint64) (bool, error) {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return false, errors.Wrap(err, "UserCreateIfNotExist")
	}

	q, args, err := sq.Expr(`SELECT array_position(category_ids, $1) as pos
							FROM users
							WHERE id=$2;`, categoryID, userID).ToSql()
	if err != nil {
		return false, errors.Wrap(err, "query creating error")
	}

	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, "db.QueryRowContext")
	}
	var pos sql.NullInt16
	if err = read(rows, &pos); err != nil {
		return false, errors.Wrap(err, "read")
	}

	return pos.Int16 != 0, nil
}

func (s *Service) GetUserCategories(ctx context.Context, userID int64) ([]string, error) {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return nil, errors.Wrap(err, "UserCreateIfNotExist")
	}

	userInfo, err := s.getUserInfo(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "getUserInfo")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblCategoriesColID, tblCategoriesColCategoryName).
		From(tblCategories).
		Where(sq.Eq{
			tblCategoriesColID: []int64(userInfo.CategoryIDs),
		}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query creating error")
	}

	var categories []category
	if err = s.db.SelectContext(ctx, &categories, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db.QueryRowContext")
	}

	data := make([]string, len(categories))
	for i := range categories {
		data[i] = categories[i].Category
	}

	return data, nil
}
