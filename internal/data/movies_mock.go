// Hey, listen! Don't edit these files!
// They're auto-generated by mockery. Go look in the
// README under the Mocks section for more info.

// Code generated by mockery v2.36.0. DO NOT EDIT.

package data

import mock "github.com/stretchr/testify/mock"

// MockMovies is an autogenerated mock type for the Movies type
type MockMovies struct {
	mock.Mock
}

type MockMovies_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMovies) EXPECT() *MockMovies_Expecter {
	return &MockMovies_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: id
func (_m *MockMovies) Delete(id int64) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMovies_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockMovies_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - id int64
func (_e *MockMovies_Expecter) Delete(id interface{}) *MockMovies_Delete_Call {
	return &MockMovies_Delete_Call{Call: _e.mock.On("Delete", id)}
}

func (_c *MockMovies_Delete_Call) Run(run func(id int64)) *MockMovies_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockMovies_Delete_Call) Return(_a0 error) *MockMovies_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMovies_Delete_Call) RunAndReturn(run func(int64) error) *MockMovies_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: id
func (_m *MockMovies) Get(id int64) (*Movie, error) {
	ret := _m.Called(id)

	var r0 *Movie
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*Movie, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *Movie); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Movie)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMovies_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockMovies_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - id int64
func (_e *MockMovies_Expecter) Get(id interface{}) *MockMovies_Get_Call {
	return &MockMovies_Get_Call{Call: _e.mock.On("Get", id)}
}

func (_c *MockMovies_Get_Call) Run(run func(id int64)) *MockMovies_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockMovies_Get_Call) Return(_a0 *Movie, _a1 error) *MockMovies_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMovies_Get_Call) RunAndReturn(run func(int64) (*Movie, error)) *MockMovies_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: movie
func (_m *MockMovies) Insert(movie *Movie) error {
	ret := _m.Called(movie)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Movie) error); ok {
		r0 = rf(movie)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMovies_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type MockMovies_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - movie *Movie
func (_e *MockMovies_Expecter) Insert(movie interface{}) *MockMovies_Insert_Call {
	return &MockMovies_Insert_Call{Call: _e.mock.On("Insert", movie)}
}

func (_c *MockMovies_Insert_Call) Run(run func(movie *Movie)) *MockMovies_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*Movie))
	})
	return _c
}

func (_c *MockMovies_Insert_Call) Return(_a0 error) *MockMovies_Insert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMovies_Insert_Call) RunAndReturn(run func(*Movie) error) *MockMovies_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: movie
func (_m *MockMovies) Update(movie *Movie) error {
	ret := _m.Called(movie)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Movie) error); ok {
		r0 = rf(movie)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMovies_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockMovies_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - movie *Movie
func (_e *MockMovies_Expecter) Update(movie interface{}) *MockMovies_Update_Call {
	return &MockMovies_Update_Call{Call: _e.mock.On("Update", movie)}
}

func (_c *MockMovies_Update_Call) Run(run func(movie *Movie)) *MockMovies_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*Movie))
	})
	return _c
}

func (_c *MockMovies_Update_Call) Return(_a0 error) *MockMovies_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMovies_Update_Call) RunAndReturn(run func(*Movie) error) *MockMovies_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMovies creates a new instance of MockMovies. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMovies(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMovies {
	mock := &MockMovies{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}