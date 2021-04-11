package wallet

import (
	"errors"
	"github.com/bekhruzdilshod/wallet/pkg/types"
	"github.com/google/uuid"
)



var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greated than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrPaymentNotCreated = errors.New("can't create payment")
var ErrFavoriteNotFound = errors.New("favorite not found")


type Service struct {
	nextAccountID	int64
	accounts	[]*types.Account
	payments []*types.Payment
	favorites []*types.Favorite
}


// s.RegisterAccount создает новый аккаунт и помещает его в хранилище Service
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	
	// Проверка номера телефона на уникальность
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	// Регистрация нового аккаунта
	s.nextAccountID++
	account := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}


// s.Deposit пополняет счет существующего аккаунта на указанную сумму
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}


// s.Pay осуществляет оплату по указанной категории и списывает средства с аккаунта пользователя
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}


// FindAccountByID осуществляет поиск аккаунта в хранилище Service по уникальному ID
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {

	var account *types.Account
	for _, acc := range s.accounts {
		if accountID == acc.ID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil

}

// FindPaymentByID осуществляет поиск платежа в хранилище Service по уникальному ID
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	
	var payment *types.Payment
	for _, pm := range s.payments {
		if paymentID == pm.ID {
			payment = pm
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}


// Reject Возвращает сумму совершенного платежа на счет плательщика и отменяет сам платеж
func (s *Service) Reject(paymentID string) error {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil || payment == nil {
		return ErrPaymentNotFound
	}
	
	accountToRefund, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return ErrAccountNotFound
	}

	payment.Status = types.PaymentStatusFail
	accountToRefund.Balance += payment.Amount
	return nil
}


// s.Repeat позволяет повторить платеж (меняется только ID)
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {

	// Ищем платеж по его ID
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	new_payment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, ErrPaymentNotCreated
	}

	return new_payment, nil

}


func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	
	var favorite *types.Favorite
	for _, fv := range s.favorites {
		if favoriteID == fv.ID {
			favorite = fv
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}



// s.FavoritePayment создает избранное из конкретного платежа
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	// Создание избранного платежа
	favorite := &types.Favorite{
		ID: uuid.New().String(),
		Name: name,
		AccountID: payment.AccountID,
		Amount: payment.Amount,
		Category: payment.Category,
	}

	// Помещение платежа в хранилище сервиса
	s.favorites = append(s.favorites, favorite)
	
	return favorite, nil

}


// s.PayFromFavorite совершает платеж из конкретного избранного (Favorite) 
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	// Поиск избранного по ID
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	// Совершение платежа по параметрам найденного избранного
	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}