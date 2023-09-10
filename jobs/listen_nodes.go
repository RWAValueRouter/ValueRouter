package jobs

import (
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/db"
)

type Node struct {
	Ip_addr     string
	Name        string
	Email       string
	Telegram_id string
	Enode       string
}

// getRegisteredNodeInfo listen registered node and stored it into db
func getRegisteredNodeInfo() {
	log.Info("listen register nodes")
	//TODO: tmp static data will be removed in the future
	var nodes []Node
	nodes = append(nodes, Node{
		Ip_addr:     "127.0.0.1:3794",
		Name:        "test node1",
		Email:       "test1@gmail.com",
		Telegram_id: "test1",
		Enode:       "enode://efcfd7d814734ccd1d83c799b5dd70ceb08ca6e490083aa84b6b72ec9c0e17bb21d6dbf410804086170c42d1e12386b29a9485a408292bd4f45e2d922c7bb0ef@127.0.0.1:30823",
	})

	nodes = append(nodes, Node{
		Ip_addr:     "127.0.0.1:3793",
		Name:        "test node2",
		Email:       "test2@gmail.com",
		Telegram_id: "test2",
		Enode:       "enode://95d5b413c834903f57c9aefad07b62c02f7f515bf37d0b30fde1e84c9cd7c1bc96253995001c39706367271584a8733d6577bbb4bcad399a7dd93cefea67b523@127.0.0.1:30824",
	})

	nodes = append(nodes, Node{
		Ip_addr:     "127.0.0.1:3792",
		Name:        "test node3",
		Email:       "test3@gmail.com",
		Telegram_id: "test3",
		Enode:       "enode://663d1f73ac3738f71ab22cb2a12d79c9b045ae423a0641bbfdd3af03354fe76ccab8002ce50ffd6ae3eeb5acd7fea40429b11dba6d59758ba2f899c091cae70c@127.0.0.1:30825",
	})

	for _, n := range nodes {
		c, err := db.Conn.GetIntValue("select count(ip_addr) from nodes_info where ip_addr = ?", n.Ip_addr)
		if err != nil {
			log.Error("DB Error", "GetIntValue", err.Error())
			return
		}
		if c > 0 {
			log.Info("DB Exist", "Ip addr", n.Ip_addr)
			continue
		}
		_, err = db.Conn.CommitOneRow("insert into nodes_info(ip_addr,name,email,telegram_id,enode) values(?,?,?,?,?)", n.Ip_addr, n.Name, n.Email, n.Telegram_id, n.Enode)
		if err != nil {
			log.Error("DB Error", "CommitOneRow", err.Error())
			return
		}
	}
}

func init() {
	jobs.AddFunc("@every 1m", getRegisteredNodeInfo)
}
