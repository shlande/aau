package bolt

import (
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/store"
)

var (
	ErrOpenTx = fmt.Errorf("%w:无法打开事务", store.ErrOperation)
)
