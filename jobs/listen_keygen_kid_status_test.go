package jobs

import (
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ListenKeyGenSuite struct {
	suite.Suite
}

func (suite *ListenKeyGenSuite) SetupTest() {
	common.Conf.DbConfig.DbDriverSource = "gotest:12345678@tcp(127.0.0.1:3306)/smw"
	common.Conf.DbConfig.DbDriverName = "mysql"
	db.Init()
	jobs.Stop()
}

func (suite *ListenKeyGenSuite) TestListenKeyGen() {
	if db.Conn.IsConnected() != nil {
		return
	}
	exist, err := db.Conn.GetIntValue("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'smw'  AND table_name = 'accounts_info'")
	var affected int64
	if assert.NoError(suite.T(), err) {
		if exist == 0 {
			affected, err = db.Conn.CommitOneRow("CREATE TABLE `accounts_info` (\n  `id` int NOT NULL AUTO_INCREMENT,\n  `gid` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '组ID',\n  `threshold` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '门限制',\n  `user_account` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户account',\n  `ip_port` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'gid对应用户对应的ipport地址',\n  `enode` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '节点对应的enode',\n  `key_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'mpc address key ID',\n  `public_key` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'public key',\n  `mpc_address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'mpc address',\n  `initializer` tinyint DEFAULT NULL COMMENT '0:not initializer ,1: initializer',\n  `reply_status` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'reply status of creating mpc wallet',\n  `reply_timestamp` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'timestamp',\n  `reply_enode` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'reply enode',\n  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0:pending , 1 SUCCESS , 2 FAIL, 3 Timeout',\n  `error` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'error message',\n  `tip` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'tip message',\n  `uuid` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'uniq identifier',\n  `local_system_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin")
			if assert.NoError(suite.T(), err) {
				assert.True(suite.T(), affected == 0)
			}
		}
	}
	_, err = db.Conn.CommitOneRow("INSERT INTO `accounts_info`(`id`,`gid`,`threshold`,`user_account`,`ip_port`, `enode`, `key_id`, `public_key`, `mpc_address`, `initializer`, `reply_status`, `reply_timestamp`, `reply_enode`, `status`, `error`, `tip`, `uuid`)  VALUES(70000,'2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73','2/2','0x24affae9c683b7615d4130300288e348e4b5d091','127.0.0.1:13793','enode://2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96@127.0.0.1:30824','0x846f37fc3c0817f2bf489ad69ef3a704f3f37b610605c97346848fb995bade08',NULL,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'0f465e35-9470-406a-b3d7-d92f603a5911')")
	assert.NoError(suite.T(), err)
	listenKeygenKidStatus()
	time.Sleep(5 * time.Second)
	_, err = db.Conn.CommitOneRow("delete from `accounts_info` where id = ?", 70000)
	assert.NoError(suite.T(), err)
}

func TestKeyGen(t *testing.T) {
	suite.Run(t, new(ListenKeyGenSuite))
}
