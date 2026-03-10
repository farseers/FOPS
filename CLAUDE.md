# CLAUDE.md
## 1、项目概览
用于日常运维的平台应用

包含：Docker集群管理、自动化部署、调度中心管理（FSchedule2.x)、日志采集、链路追踪、健康检查

前后端分离：**vue**、**go**、**mysql**、**clickhouse**、**docker swarm**

## 目录说明
**application**：应用层，同时也是API处理函数
**domain**：领域层（部份逻辑在：../lbCore/domain/lbc/）
**infrastructure**：基础设施层（部份逻辑在：../lbCore/infrastructure/lbcRepository/）
**interfaces/job**：本地定时任务
**wwwsrc**：前端vue
**/FOPS/wwwsrc/src/views/fops**: 前端页面
**/FOPS/wwwsrc/src/router/route.ts**: 前端菜单路由
**/FOPS/wwwsrc/src/api/fops/index.ts**: 定义后端api接口

## 依赖项目
**/farseer-go/**: 基础框架
**/farseer-go/docker**: docker client sdk(自己写的http或docker exec交互)
**/farseer-go/utils/exec/**: exec shell command client(自己写的shell交互)