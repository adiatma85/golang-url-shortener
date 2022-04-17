// Code generated by mockery v2.10.0. DO NOT EDIT.

package repository

import (
	helpers "github.com/adiatma85/golang-rest-template-api/pkg/helpers"
	mock "github.com/stretchr/testify/mock"

	models "github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
)

// ProductRepositoryInterface is an autogenerated mock type for the ProductRepositoryInterface type
type ProductRepositoryInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: product
func (_m *ProductRepositoryInterface) Create(product models.Product) (models.Product, error) {
	ret := _m.Called(product)

	var r0 models.Product
	if rf, ok := ret.Get(0).(func(models.Product) models.Product); ok {
		r0 = rf(product)
	} else {
		r0 = ret.Get(0).(models.Product)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Product) error); ok {
		r1 = rf(product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: product
func (_m *ProductRepositoryInterface) Delete(product *models.Product) error {
	ret := _m.Called(product)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Product) error); ok {
		r0 = rf(product)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWithIds provides a mock function with given fields: ids
func (_m *ProductRepositoryInterface) DeleteWithIds(ids []uint64) error {
	ret := _m.Called(ids)

	var r0 error
	if rf, ok := ret.Get(0).(func([]uint64) error); ok {
		r0 = rf(ids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields:
func (_m *ProductRepositoryInterface) GetAll() (*[]models.Product, error) {
	ret := _m.Called()

	var r0 *[]models.Product
	if rf, ok := ret.Get(0).(func() *[]models.Product); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: email
func (_m *ProductRepositoryInterface) GetByEmail(email string) (*models.Product, error) {
	ret := _m.Called(email)

	var r0 *models.Product
	if rf, ok := ret.Get(0).(func(string) *models.Product); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: productId
func (_m *ProductRepositoryInterface) GetById(productId string) (*models.Product, error) {
	ret := _m.Called(productId)

	var r0 *models.Product
	if rf, ok := ret.Get(0).(func(string) *models.Product); ok {
		r0 = rf(productId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(productId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: pagination
func (_m *ProductRepositoryInterface) Query(pagination helpers.Pagination) (*helpers.Pagination, error) {
	ret := _m.Called(pagination)

	var r0 *helpers.Pagination
	if rf, ok := ret.Get(0).(func(helpers.Pagination) *helpers.Pagination); ok {
		r0 = rf(pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*helpers.Pagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(helpers.Pagination) error); ok {
		r1 = rf(pagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryWithCondition provides a mock function with given fields: q, pagination
func (_m *ProductRepositoryInterface) QueryWithCondition(q *models.Product, pagination helpers.Pagination) (*helpers.Pagination, error) {
	ret := _m.Called(q, pagination)

	var r0 *helpers.Pagination
	if rf, ok := ret.Get(0).(func(*models.Product, helpers.Pagination) *helpers.Pagination); ok {
		r0 = rf(q, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*helpers.Pagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Product, helpers.Pagination) error); ok {
		r1 = rf(q, pagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: product
func (_m *ProductRepositoryInterface) Update(product *models.Product) error {
	ret := _m.Called(product)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Product) error); ok {
		r0 = rf(product)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
