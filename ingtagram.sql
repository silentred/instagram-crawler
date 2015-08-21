/*
Navicat MySQL Data Transfer

Source Server         : local
Source Server Version : 50617
Source Host           : localhost:3306
Source Database       : ingtagram

Target Server Type    : MYSQL
Target Server Version : 50617
File Encoding         : 65001

Date: 2015-08-21 14:06:36
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for picture
-- ----------------------------
DROP TABLE IF EXISTS `picture`;
CREATE TABLE `picture` (
  `orig_id` varchar(50) NOT NULL,
  `pic_url` varchar(300) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL COMMENT '0未下载； 1下载完成',
  `created_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`orig_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of picture
-- ----------------------------
INSERT INTO `picture` VALUES ('1', 'http://sdf', '1', '123');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `orig_id` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `access_token` varchar(255) DEFAULT NULL,
  `last_auth_time` int(11) DEFAULT NULL,
  `valid` tinyint(1) DEFAULT NULL COMMENT '0表示 token 失效，1表示可用',
  PRIMARY KEY (`orig_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of user
-- ----------------------------
