package wallet

import (
	"fmt"
	"reflect"
	"testing"
	"github.com/bekhruzdilshod/wallet/pkg/types"
	"github.com/google/uuid"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}


// Тестовый аккаунт
type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone: "+992900801441",
	balance: 1000,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	}{
		{amount: 100, category: "fun"},
	},
}


// Создание тестового аккаунта
func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	// Регистрация пользователя и пополняем его счет
	account, err := s.addAccountWithBalance(data.phone, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't register account, error = %v", err)
	}

	// выполняем платежи
	// Создаем слайс из платежей нужной длины
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("Can't make payments, error = %v", err)
		}
	}

	return account, payments, nil
}


// tS.addAccountWithBalance создает аккаунт по указанному номеру с указанным балансом по-умолчанию
func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	// Создаем аккаунт
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("Can't register account: %v", err)
	}

	// Пополняем счёт
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("Can't deposit account: %v", err)
	}

	return account, nil

}

// Успешная регистрация аккаунта
func TestService_RegisterAccount_success(t *testing.T) {
	s := &Service{}
	_, err := s.RegisterAccount("+992900801441")

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
}

// Аккаунт с таким номером телефона уже существует
func TestService_RegisterAccount_alreadyRegistered(t *testing.T) {
	s := &Service{}
	_, err := s.RegisterAccount("+992900801441")

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	_, err2 := s.RegisterAccount("+992900801441")
	if err2 != ErrPhoneRegistered {
		t.Errorf("Phone registered twice!")
	}
}

func TestService_FindAccountByID_found(t *testing.T) {
	s := &Service{}

	acc, err := s.RegisterAccount("+992900801441")
	if err != nil {
		t.Errorf("%v", err)
	}

	found, err := s.FindAccountByID(acc.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(found, acc) {
		t.Error("Doesn't work!")
	}

}

func TestService_FindAccountByID_notFound(t *testing.T) {
	s := &Service{}

	found, err := s.FindAccountByID(12345)
	if err == nil {
		t.Errorf("%v", found.ID)
	}

}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// Попробуем найти платеж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't find payment, error = %v", err)
	}

	// Сравниваем платежи
	if !reflect.DeepEqual(got, payment) {
		t.Errorf("FindPaymentByID(): wrong payment returned: %v", err)
		return
	}

}


func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// Попробуем найти несуществующий платеж
	got, err := s.FindPaymentByID(string(uuid.New().String()))  // payment.Status с целью избежать ошибки "declared but not used"
	if err == nil {
		t.Errorf("FindPaymentByID(): must return error, returned nil. Payment got = %v", got)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}

}


func TestService_Reject_notFound(t *testing.T) {
	s := &Service{}
	_, err := s.FindPaymentByID("1240123")
	if err == nil {
		t.Error("Oops! Something went wrong...")
	}
}

func TestService_Reject_found(t *testing.T) {

	expectedBalanceAfterReject := types.Money(100)

	s := &Service{}

	account, err := s.RegisterAccount("+992900801441")
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	s.Deposit(account.ID, 100)
	payment, err := s.Pay(account.ID, 10, "fun")
	if err != nil {
		t.Errorf("%v", err)
	}

	s.Reject(payment.ID)

	if payment.Status != types.PaymentStatusFail || account.Balance != expectedBalanceAfterReject {
		t.Error("Something went wrong... Oops!")
	}
}


func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	
	// Создаем тестовый аккаунт
	account, err := s.addAccountWithBalance("+992900801441", 100)
	if err != nil {
		t.Errorf("%v", err)
	}

	payment, err := s.Pay(account.ID, 25, "fun")
	if err != nil {
		t.Errorf("%v", err)
	}

	repeated_payment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if repeated_payment.Amount != payment.Amount || repeated_payment.Category != payment.Category {
		t.Errorf("Payments repeated wrong: categories (%v, %v), amounts (%v, %v)", payment.Category, repeated_payment.Category, payment.Amount, repeated_payment.Amount)
	} 
}


func TestService_Repeat_paymentNotFound(t *testing.T) {
	s := newTestService()
	
	// Создаем тестовый аккаунт
	account, err := s.addAccountWithBalance("+992900801441", 100)
	if err != nil {
		t.Errorf("%v", err)
	}

	_, er := s.Pay(account.ID, 25, "fun")
	if er != nil {
		t.Errorf("%v", er)
	}

	repeated_payment, err := s.Repeat(uuid.New().String())
	if err == nil {
		t.Errorf("Repeated not-exist payment: %v", repeated_payment)
	}


}
