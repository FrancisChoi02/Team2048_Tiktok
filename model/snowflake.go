package model

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

// ID_Init(string, int64) 传入偏移参照的初始时间，机器ID
func ID_Init(startTime string, machineID int64) (err error) {
	var st time.Time
	//2006-01-02这里表示的是  startTime传入的模板格式
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000

	node, err = snowflake.NewNode(machineID)
	return
}

// GenID() 生成分布式ID
func GenID() int64 {
	return node.Generate().Int64()
}
