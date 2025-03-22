// Code generated by mockery v2.52.1. DO NOT EDIT.

package classifier_mocks

import (
	context "context"

	model "github.com/bitmagnet-io/bitmagnet/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// LocalSearch is an autogenerated mock type for the LocalSearch type
type LocalSearch struct {
	mock.Mock
}

type LocalSearch_Expecter struct {
	mock *mock.Mock
}

func (_m *LocalSearch) EXPECT() *LocalSearch_Expecter {
	return &LocalSearch_Expecter{mock: &_m.Mock}
}

// ContentById provides a mock function with given fields: _a0, _a1
func (_m *LocalSearch) ContentByID(_a0 context.Context, _a1 model.ContentRef) (model.Content, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ContentByID")
	}

	var r0 model.Content
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ContentRef) (model.Content, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ContentRef) model.Content); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(model.Content)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ContentRef) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LocalSearch_ContentById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContentByID'
type LocalSearch_ContentById_Call struct {
	*mock.Call
}

// ContentById is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 model.ContentRef
func (_e *LocalSearch_Expecter) ContentById(_a0 interface{}, _a1 interface{}) *LocalSearch_ContentById_Call {
	return &LocalSearch_ContentById_Call{Call: _e.mock.On("ContentByID", _a0, _a1)}
}

func (_c *LocalSearch_ContentById_Call) Run(run func(_a0 context.Context, _a1 model.ContentRef)) *LocalSearch_ContentById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.ContentRef))
	})
	return _c
}

func (_c *LocalSearch_ContentById_Call) Return(_a0 model.Content, _a1 error) *LocalSearch_ContentById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LocalSearch_ContentById_Call) RunAndReturn(run func(context.Context, model.ContentRef) (model.Content, error)) *LocalSearch_ContentById_Call {
	_c.Call.Return(run)
	return _c
}

// ContentBySearch provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *LocalSearch) ContentBySearch(_a0 context.Context, _a1 model.ContentType, _a2 string, _a3 model.Year) (model.Content, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for ContentBySearch")
	}

	var r0 model.Content
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ContentType, string, model.Year) (model.Content, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ContentType, string, model.Year) model.Content); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(model.Content)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ContentType, string, model.Year) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LocalSearch_ContentBySearch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContentBySearch'
type LocalSearch_ContentBySearch_Call struct {
	*mock.Call
}

// ContentBySearch is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 model.ContentType
//   - _a2 string
//   - _a3 model.Year
func (_e *LocalSearch_Expecter) ContentBySearch(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *LocalSearch_ContentBySearch_Call {
	return &LocalSearch_ContentBySearch_Call{Call: _e.mock.On("ContentBySearch", _a0, _a1, _a2, _a3)}
}

func (_c *LocalSearch_ContentBySearch_Call) Run(run func(_a0 context.Context, _a1 model.ContentType, _a2 string, _a3 model.Year)) *LocalSearch_ContentBySearch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.ContentType), args[2].(string), args[3].(model.Year))
	})
	return _c
}

func (_c *LocalSearch_ContentBySearch_Call) Return(_a0 model.Content, _a1 error) *LocalSearch_ContentBySearch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LocalSearch_ContentBySearch_Call) RunAndReturn(run func(context.Context, model.ContentType, string, model.Year) (model.Content, error)) *LocalSearch_ContentBySearch_Call {
	_c.Call.Return(run)
	return _c
}

// NewLocalSearch creates a new instance of LocalSearch. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLocalSearch(t interface {
	mock.TestingT
	Cleanup(func())
}) *LocalSearch {
	mock := &LocalSearch{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
