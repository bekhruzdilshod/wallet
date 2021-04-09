package wallet

import (
	"testing"
	"reflect"
)

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