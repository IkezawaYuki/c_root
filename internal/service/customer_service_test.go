package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerService_CreateCustomer(t *testing.T) {
	fmt.Println("TestCustomerService_CreateCustomer")
	customer, err := customerSrv.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(customer))
}
