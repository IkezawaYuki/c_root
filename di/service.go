package di

import (
	"github.com/IkezawaYuki/c_root/internal/croot/infrastructre/driver"
	"github.com/IkezawaYuki/c_root/internal/croot/interfaces"
	"github.com/IkezawaYuki/c_root/internal/croot/usecase"
)

func NewCustomerService() usecase.CustomerService {
	return usecase.NewCustomerService(
		interfaces.NewCustomerAdapter(),
		interfaces.NewInstagramAdapter(
			driver.NewHttpClient()),
	)
}
