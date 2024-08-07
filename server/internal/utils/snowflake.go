// File: internal/utils/snowflake.go

package utils

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	userNode  *snowflake.Node
	matchNode *snowflake.Node
	once      sync.Once
)

func InitSnowflake() error {
	var err error
	once.Do(func() {
		userNode, err = snowflake.NewNode(1) // Node ID 1 for user IDs
		if err != nil {
			return
		}
		matchNode, err = snowflake.NewNode(2) // Node ID 2 for match IDs
	})
	return err
}

func GenerateUID() int64 {
	return userNode.Generate().Int64()
}

func GenerateMID() int64 {
	return matchNode.Generate().Int64()
}
