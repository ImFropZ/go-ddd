package query

import "github/imfropz/go-ddd/internal/application/common"

type UserQueryResult struct {
	Result *common.UserResult
}

type UserQueryListResult struct {
	Result []*common.UserResult
}
