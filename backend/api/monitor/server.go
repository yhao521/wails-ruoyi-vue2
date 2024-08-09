package monitor

import (
	"context"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/monitor"
	"mySparkler/pkg/utils/R"
	"sync"
)

type ServerAPI struct {
	ctx context.Context
	baseAPI.Base
}

var serverAPI *ServerAPI
var onceServerAPI sync.Once

// NewApp creates a new App application struct
func NewServerAPI() *ServerAPI {
	if serverAPI == nil {
		onceServerAPI.Do(func() {
			serverAPI = &ServerAPI{}
		})
	}
	return serverAPI
}

func (g *ServerAPI) ServerData() R.Result {
	var server = monitor.GetServerInfo()
	return R.ReturnSuccess(server)
}
