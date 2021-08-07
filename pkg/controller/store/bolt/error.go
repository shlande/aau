package bolt

import (
	"fmt"
	store2 "github.com/shlande/dmhy-rss/pkg/controller/store"
)

var (
	ErrOpenTx = fmt.Errorf("%w:无法打开事务", store2.ErrOperation)
)
