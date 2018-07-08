package main

import (
	"log"

	"mj/utils"

	"github.com/88bytes/nano"
	"github.com/88bytes/nano/serialize/json"
)

// 不错，测试通过，已经可以完美的和麻将的客户端通信了
// 服务器的优化，支持分布式，更优的网络数据大小，更小的GC 这些以后再慢慢来吧
func main() {
	utils.UserInfoUtilInst = utils.NewUserInfoUtil()
	utils.FreeMatchMgrInst = utils.NewFreeMatchMgr()
	utils.PVPMgrInst = utils.NewPVPMgr()

	nano.SetSerializer(json.NewSerializer())
	nano.Register(&Login{})
	nano.Register(&Match{})
	nano.Register(&PVP{})
	// nano.EnableDebug()
	log.SetFlags(log.LstdFlags | log.Llongfile)
	nano.Listen(":8010")
}
