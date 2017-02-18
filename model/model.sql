/* NOTE: AUTO-GENERATED by midc, DON'T edit!! */

CREATE TABLE IF NOT EXISTS `user` (
	`id` BIGINT(20)   COMMENT '随机唯一Id',
	`account_type` BIGINT(20)   COMMENT '账号类型',
	`account` VARCHAR(128) UNIQUE  COMMENT '账号',
	`nickname` VARCHAR(32)   COMMENT '昵称',
	`avatar` VARCHAR(256)   COMMENT '头像',
	`qrcode` TEXT   COMMENT '二维码',
	`gender` BIGINT(20)   COMMENT '性别',
	`birthday` VARCHAR(32)   COMMENT '生日',
	`id_card_type` BIGINT(20)   COMMENT '身份证件类型',
	`id_card` VARCHAR(64)   COMMENT '证件唯一标识',
	`encrypted_password` VARCHAR(64)   COMMENT '加密后密码',
	`password_salt` VARCHAR(64)   COMMENT '加密密码的盐',
	`created_at` VARCHAR(32)   COMMENT '账号创建时间',
	`created_ip` VARCHAR(32)   COMMENT '账号创建时IP',
	`last_login_at` VARCHAR(32)   COMMENT '最后登陆时间',
	`last_login_ip` VARCHAR(32)   COMMENT '最后登陆时IP',
	PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `client` (
	`id` VARCHAR(64)   COMMENT 'oauth 客户端唯一Id',
	`secret` VARCHAR(64)   COMMENT '密码',
	`name` VARCHAR(64)   COMMENT '应用名称',
	`description` TEXT   COMMENT '应用描述',
	`scope` TEXT   COMMENT '授权范围',
	`callback_url` VARCHAR(256)   COMMENT '回调地址',
	PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `access_token` (
	`id` BIGINT(20) AUTO_INCREMENT  COMMENT '递增唯一Id',
	`uid` BIGINT(20)   COMMENT '用户Id',
	`created_at` VARCHAR(32)   COMMENT '创建时间',
	`modified_at` VARCHAR(32)   COMMENT '修改时间',
	`expire_at` VARCHAR(32)   COMMENT '到期时间',
	`token` VARCHAR(64) UNIQUE  COMMENT '令牌',
	`refresh_token` VARCHAR(64) UNIQUE  COMMENT '刷新用令牌',
	`resource_owner` VARCHAR(64)   COMMENT '资源所有者',
	`client_id` VARCHAR(64)   COMMENT '客户Id',
	`scope` TEXT   COMMENT '可访问权限范围',
	PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `authorization_request` (
	`id` BIGINT(20) AUTO_INCREMENT  COMMENT '递增唯一Id',
	`created_at` VARCHAR(32)   COMMENT '创建时间',
	`authorization_code` VARCHAR(64)   COMMENT '认证码',
	`uid` BIGINT(20)   COMMENT '用户Id',
	`redirect_uri` VARCHAR(256)   COMMENT '重定向URI',
	`response_type` VARCHAR(64)   COMMENT '返回类型',
	`state` VARCHAR(128)   COMMENT '自定义状态',
	`client_id` VARCHAR(64)   COMMENT '客户端Id',
	`granted_scopes` TEXT   COMMENT '授权范围',
	`requested_scopes` TEXT   COMMENT '请求范围',
	PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `telno_verify_code` (
	`id` BIGINT(20) AUTO_INCREMENT  COMMENT '递增唯一Id',
	`telno` VARCHAR(32)   COMMENT '手机号码',
	`code` VARCHAR(64) UNIQUE  COMMENT '验证码',
	`expired_at` VARCHAR(32)   COMMENT '到期时间',
	PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `email_verify_code` (
	`id` BIGINT(20) AUTO_INCREMENT  COMMENT '递增唯一Id',
	`email` VARCHAR(64)   COMMENT 'email 地址',
	`code` VARCHAR(64) UNIQUE  COMMENT '验证码',
	`expired_at` VARCHAR(32)   COMMENT '到期时间',
	PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `session` (
	`id` VARCHAR(64)   COMMENT '唯一Id,用作cookie',
	`uid` BIGINT(20)   COMMENT '关联的用户Id',
	`expire_at` VARCHAR(32)   COMMENT '到期时间',
	PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

