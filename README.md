# dump
备份MySQL数据库,速度快

配置文件目录:./config

备份时,运行dump.go,根据配置文件中的databases选取不同的数据库,对应生成:dumps/生成日期/数据库/*.sql

运行SQL时,运行run.go,根据配置文件中的databases,分别选取 dumps/生成日期/数据库名称/*.sql 运行
