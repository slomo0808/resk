-- 红包商品表的简化版本，为了方便测试评估，至保留了部分字段
-- 红包商品表
drop table if exists `goods`;

create table `goods`
(
  `id`              bigint(20)     not null auto_increment comment '商品id，自动生成',
  `envelope_no`     varchar(32)    not null comment '红包编号',
  `remain_amount`   decimal(30, 6) not null default '0.000000',
  `remain_quantity` int(10)        not null comment '红包剩余数量',
  `created_at`      datetime(3)    not null default current_timestamp(3) comment '创建时间',
  `updated_at`      datetime(3)    not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
  primary key (`id`) using btree,
  unique key `envelope_no_idx` (`envelope_no`) using btree
) engine = InnoDB
  AUTO_INCREMENT = 156
  DEFAULT CHARSET = utf8
  ROW_FORMAT = DYNAMIC;

-- 红包商品表的无符号类型字段版本
drop table if exists `goods_unsigned`;
create table `goods_unsigned`
(
  `id`              bigint(20)              not null auto_increment comment '商品id，自动生成',
  `envelope_no`     varchar(32)             not null comment '红包编号',
  `remain_amount`   decimal(30, 6) unsigned not null default '0.000000',
  `remain_quantity` int(10) unsigned        not null comment '红包剩余数量',
  `created_at`      datetime(3)             not null default current_timestamp(3) comment '创建时间',
  `updated_at`      datetime(3)             not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
  primary key (`id`) using btree,
  unique key `envelope_no_idx` (`envelope_no`) using btree
) engine = InnoDB
  AUTO_INCREMENT = 156
  DEFAULT CHARSET = utf8
  ROW_FORMAT = DYNAMIC;