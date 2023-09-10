-- MySQL dump 10.13  Distrib 8.0.13, for macos10.14 (x86_64)
--
-- Host: localhost    Database: smw
-- ------------------------------------------------------
-- Server version	8.0.23

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `smw_enodes_info`
--

DROP TABLE IF EXISTS `signs_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `signs_detail` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key_id` varchar(128) COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc sign key ID',
  `user_account` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '用户account',
  `group_id` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '组ID',
  `enode` varchar(512) COLLATE utf8mb4_bin NOT NULL COMMENT '节点对应的enode',
  `ip_port` varchar(512) COLLATE utf8mb4_bin NOT NULL COMMENT 'gid对应用户对应的ipport地址',
  `threshold` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT '门限制',
  `msg_hash`  varchar(1024) COLLATE utf8mb4_bin NOT NULL COMMENT 'msg signed hash values , using | to separate multiple hashes',
  `msg_context` mediumtext COLLATE utf8mb4_bin NOT NULL COMMENT 'msg_context list separated by |',
  `rsv` mediumtext COLLATE utf8mb4_bin COMMENT 'rsv list separated by |',
  `public_key` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc address public key',
  `mpc_address` varchar(128) COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc address',
  `reply_initializer` tinyint(2) COMMENT '0:not initializer ,1: initializer',
  `reply_status` varchar(128) COLLATE utf8mb4_bin COMMENT 'reply status of creating mpc wallet',
  `reply_timestamp` varchar(128) COLLATE utf8mb4_bin COMMENT 'reply timestamp',
  `reply_enode` varchar(512) COLLATE utf8mb4_bin COMMENT 'reply enode',
  `initiator_public_key` varchar(128) COLLATE utf8mb4_bin COMMENT 'initial public key',
  `key_type` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'key类型',
  `mode` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'Mode模式',
  `status` tinyint(2) NOT NULL DEFAULT 0 COMMENT '0:pending , 1 SUCCESS , 2 FAIL, 3 Timeout, 4 Tx-Pending, 5 Tx-Confirmed, 6 Tx-Failed, 7 Tx-NotValid, 8 Tx-Replaced',
  `error` varchar(512) COLLATE utf8mb4_bin COMMENT 'error message',
  `tip`  varchar(512) COLLATE utf8mb4_bin COMMENT 'tip message',
  `txid` varchar(128) COLLATE utf8mb4_bin COMMENT 'signing tx id',
  `chain_type` int COLLATE utf8mb4_bin DEFAULT 0 COMMENT 'chain type',
  `local_system_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `chain_id` int COLLATE utf8mb4_bin DEFAULT 0 COMMENT 'chain id',
  `sign_type` int COLLATE utf8mb4_bin DEFAULT 0 COMMENT 'sign type',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-02-16 13:35:55

alter table signs_detail add column `chain_id` int COLLATE utf8mb4_bin DEFAULT 0 COMMENT 'chain id';
alter table signs_detail add column `sign_type` int COLLATE utf8mb4_bin DEFAULT 0 COMMENT 'sign_type';

create index idx_signs_detail_kid_acct on signs_detail (key_id, user_account);
create index idx_signs_detail_status on signs_detail (status);
