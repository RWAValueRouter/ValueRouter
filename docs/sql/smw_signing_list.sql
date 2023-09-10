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

DROP TABLE IF EXISTS `signing_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `signing_list` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key_id` varchar(128) COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc sign key ID',
  `user_account` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '用户account',
  `group_id` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '组ID',
  `ip_port` varchar(512) COLLATE utf8mb4_bin NOT NULL COMMENT 'gid对应用户对应的ipport地址',
  `key_type` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'key类型',
  `mode` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'Mode模式',
  `msg_hash`  varchar(1024) COLLATE utf8mb4_bin NOT NULL COMMENT 'msg signed hash values , using | to separate multiple hashes',
  `msg_context` mediumtext COLLATE utf8mb4_bin NOT NULL COMMENT 'msg_context list separated by |',
  `nonce` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'user account nonce',
  `public_key` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc address public key',
  `mpc_address` varchar(128) COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc address',
  `threshold` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT '门限制',
  `timestamp` varchar(128) COLLATE utf8mb4_bin COMMENT 'timestamp',
  `status` tinyint(2) NOT NULL DEFAULT 0 COMMENT '0:pending , 1 Done',
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
alter table signing_list add column `chain_id` int COLLATE utf8mb4_bin DEFAULT 0 COMMENT 'chain id';
alter table signing_list add column `chain_type` int COLLATE utf8mb4_bin DEFAULT 0 COMMENT 'chain type';

create index idx_signing_list_kid_status on signing_list (key_id, status);
create index idx_signing_list_acct_status on signing_list (user_account, status);