package dao

import (
	"github.com/stretchr/testify/mock"
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

func Test_db_FetchMessage(t *testing.T) {
	dbMock := dbMock{}
	type returnData struct {
		val string
		err error
	}
	tests := []struct {
		name     string
		db       DBService
		language string
		want     returnData
		wantErr  returnData
	}{
		// TODO: Add test cases.
		{"test 1 english ", dbMock, "en", returnData{"aaa", nil}, returnData{"aaa", nil}},
		{"test 1 vietnames ", dbMock, "vi", returnData{"xin chao", nil}, returnData{"na ni", nil}},
		{"test 1 default ", dbMock, "", returnData{"xxxx", nil}, returnData{"na", nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			language := ""
			if len(tt.language) != 0 {
				language = tt.language
			}

			dbMock.On("FetchMessage", language).Return(tt.want.val, tt.want.err)

			got, err := dbMock.FetchMessage(language)
			if (err != nil) || (got != tt.wantErr.val) {
				t.Errorf("db.FetchDefaultMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_db_FetchDefaultMessage(t *testing.T) {
	dbMock := dbMock{}
	type returnData struct {
		val string
		err error
	}
	tests := []struct {
		name    string
		db      DBService
		want    returnData
		wantErr returnData
	}{
		// TODO: Add test cases.
		{"test 1 english ", dbMock, returnData{"aaa", nil}, returnData{"aaa", nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock.On("FetchDefaultMessage").Return("aaa", nil)
			got, err := dbMock.FetchDefaultMessage()
			if (err != nil) && (got != tt.wantErr.val) {
				t.Errorf("db.FetchDefaultMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
