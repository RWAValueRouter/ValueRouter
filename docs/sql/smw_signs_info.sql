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

DROP TABLE IF EXISTS `signs_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `signs_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `account` varchar(128) COLLATE utf8mb4_bin NOT NULL COMMENT '节点账户',
  `nonce` int COLLATE utf8mb4_bin NOT NULL COMMENT 'user account nonce',
  `pubkey` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '签名公钥',
  `msg_hash` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '代签名消息hash值',
  `msg_context` mediumtext COLLATE utf8mb4_bin NOT NULL COMMENT 'msg 原文',
  `key_type` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'sign type',
  `group_id` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT 'group id of sign , can be sub_gid or gid',
  `threshold` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'sign threshold',
  `mod` varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT 'sign mod',
  `accept_timeout` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT 'accept timeout',
  `timestamp` varchar(128) COLLATE utf8mb4_bin NOT NULL COMMENT 'sign timestamp',
  `key_id` varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT 'sign keyId',
  `status` tinyint(2) NOT NULL DEFAULT 0 COMMENT '0:normal ,-1: invalid',
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

-- Dump completed on 2023-02-16 13:35:55

create index idx_signs_info_kid on signs_info (key_id);
