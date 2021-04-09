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


type Service struct {
	nextAccountID	int64
	accounts	[]*types.Account
	payments []*types.Payment
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