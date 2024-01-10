-- 账号表
CREATE TABLE `admin` (
     `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
     `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
     `password` varchar(96) NOT NULL DEFAULT '' COMMENT '密码',
     `role_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '角色ID',
     `name` varchar(256) NOT NULL DEFAULT '' COMMENT '姓名',
     `nickname` varchar(256) NOT NULL DEFAULT '' COMMENT '昵称',
     `mobile` varchar(256) NOT NULL DEFAULT '' COMMENT '手机号',
     `last_login_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最后一次登录时间',
     `last_login_ip` char(20) NOT NULL DEFAULT '' COMMENT '最后登录ip',
     `status` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1启用,2禁用',
     `is_delete` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除 0否1是',
     `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
     `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
     PRIMARY KEY (`id`) USING BTREE,
     KEY `idx_role_id` (`role_id`) USING BTREE,
     KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='账号表';

-- 公告
CREATE TABLE `notice` (
      `id` int(8) NOT NULL AUTO_INCREMENT,
      `title` varchar(256) NOT NULL DEFAULT '' COMMENT '标题',
      `image` text COMMENT '图片',
      `cate_id` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '种类 1功能更新 2活动通知',
      `is_delete` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除 0否1是',
      `sort` int(1) unsigned NOT NULL DEFAULT '1' COMMENT '排序',
      `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
      `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
      PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='公告';


CREATE TABLE `user` (
        `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
        `openid` varchar(64) NOT NULL DEFAULT '' COMMENT 'openid',
        `anonymous_openid` varchar(256) NOT NULL DEFAULT '' COMMENT '接口返回的匿名登录凭证',
        `unionid` varchar(64) NOT NULL DEFAULT '' COMMENT 'unionId',
        `session_key` varchar(64) NOT NULL DEFAULT '' COMMENT 'session_key',
        `nickname` varchar(256) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '昵称',
        `gender` tinyint(1) NOT NULL DEFAULT '0' COMMENT '用户的性别，1=男性，2=女性，0=未知',
        `city` varchar(256) NOT NULL DEFAULT '' COMMENT '用户所在城市',
        `country` varchar(256) NOT NULL DEFAULT '' COMMENT '用户所在国家',
        `province` varchar(256) NOT NULL DEFAULT '' COMMENT '用户所在省份',
        `language` varchar(256) NOT NULL DEFAULT '' COMMENT '用户的语言，简体中文为zh_CN',
        `avatar_url` varchar(256) NOT NULL DEFAULT '' COMMENT '头像',
        `mobile` varchar(16) NOT NULL COMMENT '手机号，带区号',
        `country_code` varchar(8) NOT NULL COMMENT '国家码',
        `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
        `vip_expires_time` int(10) NOT NULL DEFAULT '0' COMMENT '会员到期时间',
        `is_banned` tinyint(1) NOT NULL DEFAULT '0' COMMENT '用户是否被禁用，0否，1是',
        `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
        `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后更新时间',
        PRIMARY KEY (`id`) USING BTREE,
        KEY `idx_unionid` (`unionid`) USING BTREE,
        KEY `idx_openid` (`openid`) USING BTREE,
        KEY `idx_anonymous_openid` (`anonymous_openid`(255)) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户信息表';
