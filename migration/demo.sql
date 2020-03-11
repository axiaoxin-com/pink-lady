-- phpMyAdmin SQL Dump
-- version 4.6.3
-- https://www.phpmyadmin.net/
--
-- Host: 100.115.82.43:3846
-- Generation Time: 2019-12-01 16:11:43
-- 服务器版本： 5.6.28-cdb2016-log
-- PHP Version: 7.1.1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `clm`
--

-- --------------------------------------------------------

--
-- 表的结构 `alert_filter_rule`
--

CREATE TABLE `alert_filter_rule` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键ID',
  `alert_policy_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '关联的告警策略ID',
  `relation` tinyint(4) NOT NULL DEFAULT '0' COMMENT '与或关系（1=与，2=或）',
  `field` varchar(128) NOT NULL DEFAULT '' COMMENT '日志原始字段名',
  `operating` varchar(8) NOT NULL DEFAULT '' COMMENT '操作符',
  `value` varchar(255) NOT NULL DEFAULT '' COMMENT '筛选值（多个值逗号分隔）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='告警筛选规则表';

-- --------------------------------------------------------

--
-- 表的结构 `alert_policy`
--

CREATE TABLE `alert_policy` (
  `id` bigint(20) NOT NULL COMMENT '主键id',
  `appid` int(11) NOT NULL DEFAULT '0' COMMENT 'AppId',
  `uin` varchar(64) NOT NULL DEFAULT '' COMMENT 'Uin',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '策略名称',
  `metric_set_id` int(11) NOT NULL DEFAULT '0' COMMENT '指标集ID',
  `notice_frequency_sec` int(11) NOT NULL DEFAULT '0' COMMENT '通知频率（通知间隔秒数）',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 1=已开启 2=未开启 3=已失效',
  `alert_group_id` varchar(1024) NOT NULL DEFAULT '' COMMENT '告警接收组 逗号分隔',
  `alert_channel` varchar(16) NOT NULL DEFAULT '0' COMMENT '告警接收渠道 1=邮件 2=短信 3=微信',
  `notice_period_begin` int(11) NOT NULL DEFAULT '0' COMMENT '通知时段开始时间（从00:00:00开始计算的秒数）',
  `notice_period_end` int(11) NOT NULL DEFAULT '0' COMMENT '通知时段结束时间（从00:00:00开始计算的秒数）',
  `url_scheme` varchar(5) NOT NULL DEFAULT '' COMMENT '回调url的scheme',
  `callback_url` varchar(255) NOT NULL DEFAULT '' COMMENT '回调url 不包含scheme部分',
  `latest_alert_time` varchar(32) NOT NULL DEFAULT '' COMMENT '最后告警时间（产生告警后更改该字段）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='告警策略表';

-- --------------------------------------------------------

--
-- 表的结构 `alert_trigger_rule`
--

CREATE TABLE `alert_trigger_rule` (
  `id` bigint(20) NOT NULL COMMENT '主键id',
  `alert_policy_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '关联的告警策略主键id',
  `metric_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '指标id',
  `metric_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '指标类型 1=普通指标 2=复合指标',
  `operating` varchar(12) NOT NULL DEFAULT '' COMMENT '操作符',
  `value` varchar(64) NOT NULL DEFAULT '' COMMENT '指标阈值',
  `continuous_cycle_count` int(11) NOT NULL DEFAULT '0' COMMENT '持续周期个数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `relation` tinyint(4) DEFAULT '0' COMMENT '与或关系（1=与，2=或）'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='告警策略触发条件表';

--
-- Indexes for dumped tables
--

--
-- Indexes for table `alert_filter_rule`
--
ALTER TABLE `alert_filter_rule`
  ADD PRIMARY KEY (`id`),
  ADD KEY `i_alert_policy_id` (`alert_policy_id`);

--
-- Indexes for table `alert_policy`
--
ALTER TABLE `alert_policy`
  ADD PRIMARY KEY (`id`),
  ADD KEY `i_name` (`name`);

--
-- Indexes for table `alert_trigger_rule`
--
ALTER TABLE `alert_trigger_rule`
  ADD PRIMARY KEY (`id`),
  ADD KEY `i_alert_policy_id` (`alert_policy_id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `alert_filter_rule`
--
ALTER TABLE `alert_filter_rule`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID';
--
-- 使用表AUTO_INCREMENT `alert_policy`
--
ALTER TABLE `alert_policy`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id';
--
-- 使用表AUTO_INCREMENT `alert_trigger_rule`
--
ALTER TABLE `alert_trigger_rule`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id';
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
