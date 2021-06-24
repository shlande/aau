package blot

import (
	"fmt"
)

var (
	ErrOpenTx = fmt.Errorf("%w:无法打开事务")
)
