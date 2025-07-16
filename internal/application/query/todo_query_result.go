package query

import (
	"github/imfropz/go-ddd/internal/application/common"
)

type TodoQueryResult struct {
	Result *common.TodoResult
}

type TodoQueryListResult struct {
	Result []*common.TodoResult
}
