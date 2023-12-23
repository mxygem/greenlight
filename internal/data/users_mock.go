// Hey, listen! Don't edit these files!
// They're auto-generated by mockery. Go look in the
// README under the Mocks section for more info.

// Code generated by mockery v2.36.0. DO NOT EDIT.

package data

import mock "github.com/stretchr/testify/mock"

// MockUsers is an autogenerated mock type for the Users type
type MockUsers struct {
	mock.Mock
}

type MockUsers_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUsers) EXPECT() *MockUsers_Expecter {
	return &MockUsers_Expecter{mock: &_m.Mock}
}

// GetByEmail provides a mock function with given fields: email
func (_m *MockUsers) GetByEmail(email string) (*User, error) {
	ret := _m.Called(email)

	var r0 *User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUsers_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type MockUsers_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - email string
func (_e *MockUsers_Expecter) GetByEmail(email interface{}) *MockUsers_GetByEmail_Call {
	return &MockUsers_GetByEmail_Call{Call: _e.mock.On("GetByEmail", email)}
}

func (_c *MockUsers_GetByEmail_Call) Run(run func(email string)) *MockUsers_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockUsers_GetByEmail_Call) Return(_a0 *User, _a1 error) *MockUsers_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUsers_GetByEmail_Call) RunAndReturn(run func(string) (*User, error)) *MockUsers_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: user
func (_m *MockUsers) Insert(user *User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUsers_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type MockUsers_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - user *User
func (_e *MockUsers_Expecter) Insert(user interface{}) *MockUsers_Insert_Call {
	return &MockUsers_Insert_Call{Call: _e.mock.On("Insert", user)}
}

func (_c *MockUsers_Insert_Call) Run(run func(user *User)) *MockUsers_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*User))
	})
	return _c
}

func (_c *MockUsers_Insert_Call) Return(_a0 error) *MockUsers_Insert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUsers_Insert_Call) RunAndReturn(run func(*User) error) *MockUsers_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: user
func (_m *MockUsers) Update(user *User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUsers_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockUsers_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - user *User
func (_e *MockUsers_Expecter) Update(user interface{}) *MockUsers_Update_Call {
	return &MockUsers_Update_Call{Call: _e.mock.On("Update", user)}
}

func (_c *MockUsers_Update_Call) Run(run func(user *User)) *MockUsers_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*User))
	})
	return _c
}

func (_c *MockUsers_Update_Call) Return(_a0 error) *MockUsers_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUsers_Update_Call) RunAndReturn(run func(*User) error) *MockUsers_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUsers creates a new instance of MockUsers. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUsers(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUsers {
	mock := &MockUsers{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}