package driver

import "github.com/IkezawaYuki/c_root/internal/croot/repository/filter"

type DBDriver interface {
	First(model interface{}, filter filter.Filter) error
}
