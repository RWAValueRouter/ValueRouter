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
-- Table structure for table `smw_group_info`
--

DROP TABLE IF EXISTS `groups_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `groups_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `tx_type` varchar(24) COLLATE utf8mb4_bin NOT NULL COMMENT '交易类型',
  `account` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '交易发送方',
  `nonce` bigint(20) COLLATE utf8mb4_bin NOT NULL COMMENT '交易发送方当前的nonce值',
  `key_type`varchar(24) COLLATE utf8mb4_bin NOT NULL COMMENT '加密算法类型',
  `group_id` varchar(128) COLLATE utf8mb4_bin NOT NULL COMMENT '组ID',
  `thres_hold` varchar(24) COLLATE utf8mb4_bin NOT NULL COMMENT '门限值 ，2/3',
  `mode` tinyint(2) NOT NULL COMMENT '0：需要审批，gid用subid，1：不需要审批，自动审批，2：需要审批，但是现将签名发送到总组gid就是总组的id，然后再子组竞争，谁先签名，谁完成交易',
  `accept_timeout` int(11) NOT NULL COMMENT '接收组的超时时间',
  `sigs` mediumtext COLLATE utf8mb4_bin NOT NULL COMMENT 'enode list separated by |',
  `timestamp` varchar(64) NOT NULL COMMENT '交易发送时间戳',
  `key_id` varchar(128) COLLATE utf8mb4_bin COMMENT 'mpc address key id',
  `uuid` varchar(128) COLLATE utf8mb4_bin COMMENT 'uniq identifier',
  `local_system_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
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

-- Dump completed on 2023-02-16 14:06:07
