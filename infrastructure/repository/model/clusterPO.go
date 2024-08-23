package model

// ClusterPO 集群
type ClusterPO struct {
	Id             int64  `gorm:"primaryKey;comment:主键"`
	Name           string `gorm:"size:32;not null;comment:集群名称"`
	FopsAddr       string `gorm:"size:64;not null;comment:集群地址"`
	FScheduleAddr  string `gorm:"size:64;not null;comment:调度中心地址"`
	DockerHub      string `gorm:"size:256;not null;comment:托管地址"`
	DockerUserName string `gorm:"size:32;not null;comment:账户名称"`
	DockerUserPwd  string `gorm:"size:64;not null;comment:账户密码"`
	DockerNetwork  string `gorm:"size:64;not null;comment:Docker网络"`
	IsLocal        bool   `gorm:"size:1;not null;default:0;comment:本地集群"`
}
