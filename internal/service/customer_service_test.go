package service

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCustomerService_CreateCustomer(t *testing.T) {
	fmt.Println("TestCustomerService_CreateCustomer")
	customers, err := customerSrv.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(customers))

	err = customerSrv.CreateCustomer(context.Background(), &entity.Customer{
		ID:             1,
		Name:           "test",
		Password:       "test",
		Email:          "email",
		WordpressURL:   "wordpressURL",
		FacebookToken:  Ptr("facebookToken"),
		StartDate:      Ptr(time.Now()),
		InstagramID:    Ptr("instagramID"),
		InstagramName:  Ptr("instagramName"),
		DeleteHashFlag: 0,
	})
	assert.NoError(t, err)

	customer, err := customerSrv.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "test", customer.Name)
	fmt.Println(customer)
}

func Ptr[T any](v T) *T {
	return &v
}
