package respository

import (
	"github.com/stretchr/testify/mock"
	"grabvn-golang-bootcamp/week_4/practice/todo/protobuf"
)

type FakeTodoRepo struct {
	Repo mock.Mock
}

func (fakeRepo *FakeTodoRepo) Create(item *protobuf.CreateTodoRequest) (result string, err error) {
	ret := fakeRepo.Repo.Called(item)

	if rf, ok := ret.Get(0).(func(*protobuf.CreateTodoRequest) string); ok {
		result = rf(item)
	} else {
		if ret.Get(0) != nil {
			result = ret.Get(0).(string)
		} else if ret.Get(1) != nil {
			err = ret.Error(1)
		}
	}
	return
}

func (fakeRepo *FakeTodoRepo) Update(id string, item *protobuf.TodoRequestUpdateInfo) (err error) {
	ret := fakeRepo.Repo.Called(id, item)

	if rf, ok := ret.Get(0).(func(string, *protobuf.TodoRequestUpdateInfo) error); ok {
		return rf(id, item)
	} else {
		if ret.Get(0) == nil {
			return
		}
		return ret.Error(1)
	}
}

func (fakeRepo *FakeTodoRepo) Get(id string) (result *protobuf.Todo, err error) {
	ret := fakeRepo.Repo.Called(id)

	if rf, ok := ret.Get(0).(func(string) *protobuf.Todo); ok {
		result = rf(id)
	} else {
		if ret.Get(0) != nil {
			result = ret.Get(0).(*protobuf.Todo)
		} else {
			return nil, ret.Error(1)
		}
	}
	if rf, ok := ret.Get(1).(func(string) error); ok {
		err = rf(id)
	} else {
		err = ret.Error(1)
	}
	return
}

func (fakeRepo *FakeTodoRepo) GetList(limit int32, marker string, complete bool) (result []*protobuf.Todo, err error) {
	ret := fakeRepo.Repo.Called(limit, marker, complete)
	if rf, ok := ret.Get(0).(func(int32, string, bool) []*protobuf.Todo); ok {
		result = rf(limit, marker, complete)
	} else {
		if ret.Get(0) != nil {
			result = ret.Get(0).([]*protobuf.Todo)
		} else {
			return nil, ret.Error(1)
		}
	}

	if ret.Get(1) != nil {
		err = ret.Error(1)
	}

	return
}

func (fakeRepo FakeTodoRepo) Delete(id string) (err error) {
	ret := fakeRepo.Repo.Called(id)

	if ret.Get(0) != nil {
		return ret.Error(0)
	}
	return
}
