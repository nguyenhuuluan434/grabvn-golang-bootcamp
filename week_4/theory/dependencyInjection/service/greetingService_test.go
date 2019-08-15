package service

import (
	"github.com/stretchr/testify/mock"
	"grabvn-golang-bootcamp/week_4/theory/dependencyInjection/dao"
	"testing"
)

type dbMock struct {
	mock.Mock
}

func (db dbMock) FetchMessage(lang string) (string, error) {
	args := db.Called(lang)
	return args.String(0), args.Error(1)
}

func (db dbMock) FetchDefaultMessage() (string, error) {
	args := db.Called()
	return args.String(0), args.Error(1)
	//return args.Get(0).(*MyObject), args.Get(1).(*AnotherObjectOfMine)
}

func Test_greeter_Greet(t *testing.T) {
	dbMock := &dbMock{}
	type fields struct {
		database dao.DBService
		lang     string
	}
	tests := []struct {
		name   string
		fields fields
		returnVal string
		want   string
	}{
		// TODO: Add test cases.
		{name: "test en", fields: fields{lang: "vi", database: dbMock}, want: "aaa",returnVal:"aaa"},
		{name: "test vn", fields: fields{lang: "en", database: dbMock}, want: "123",returnVal:"123"},
		{name: "test df", fields: fields{lang: "", database: dbMock}, want: "",returnVal:""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock.On("FetchMessage", tt.fields.lang).Return(tt.returnVal,nil)
			g := &greeter{
				database: tt.fields.database,
				lang:     tt.fields.lang,
			}

			if got := g.Greet(); got != tt.want {
				t.Errorf("greeter.Greet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_greeter_GreetInDefaultMsg(t *testing.T) {
	type fields struct {
		database dao.DBService
		lang     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &greeter{
				database: tt.fields.database,
				lang:     tt.fields.lang,
			}
			if got := g.GreetInDefaultMsg(); got != tt.want {
				t.Errorf("greeter.GreetInDefaultMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
