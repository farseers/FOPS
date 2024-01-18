<template>
  <div class="layout-padding" style="position: relative;">
    <el-card shadow="hover">
      <el-header style="padding: 0">
        <el-select v-model="state.clusterId" placeholder="请选择集群" class="ml10" @change="onClusterChange">
          <el-option v-for="item in state.clusterData" :key="item.Id" :label="item.Name" :value="item.Id"></el-option>
        </el-select>
        <el-button size="default" type="success" class="ml10" @click="onOpenAdd('add')"><el-icon><ele-FolderAdd /></el-icon>新增应用</el-button>
        <el-button size="default" type="info" class="ml10" @click="onClearDockerImage('add')"><el-icon><ele-Delete /></el-icon>清除None镜像</el-button>
        <el-button size="default" type="danger" class="ml10" @click="onAllBuild()"><el-icon><ele-SwitchButton /></el-icon>全部构建</el-button>
      </el-header>
      <!--应用列表-->
      <el-container>
        <el-main style="padding: 0">
          <el-space wrap style="align-items: unset;">
            <el-card shadow="hover" v-for="(v, k) in state.tableData.data" :key="k" style="width: 280px;">
              <template #header>
                <div class="card-header" style="height: 20px;">
                  <el-tag size="default">{{ v.AppName }}</el-tag>
                  <el-tag v-if="v.IsHealth" size="small" type="success" style="margin-left: 5px">健康</el-tag>
                  <el-tag v-else-if="v.ActiveInstance!=null && v.ActiveInstance.length > 0" size="small" type="warning" style="margin-left: 5px">不健康</el-tag>
                  <el-tag v-else size="small" type="danger" style="margin-left: 5px">未运行</el-tag>
                  <el-tooltip content="实例数量/副本数量" slot="label">
                    <el-tag size="small" style="margin-left: 5px">{{v.ActiveInstance.length}}/{{ v.DockerReplicas }}</el-tag>
                  </el-tooltip>
                  <el-button class="button" size="small" @click="onOpenEdit('edit', v)" type="warning" style="float:right;position: relative;">修改</el-button>
                </div>
              </template>
                <el-button size="small" type="success"  @click="showFsLogLevel(2,v.AppName)" style="float:right;position: relative;margin-left: 5px">日志</el-button>
                <el-button size="small" @click="onRestartDocker(v)" type="warning" style="float:right;position: relative;"><el-icon><ele-SwitchButton /></el-icon>重启</el-button>
                <div class="appItem" style="margin-bottom: 10px">仓库版本
                  <div class="appItem">
                    <el-tag v-if="v.DockerImage !=''" size="small">{{ v.DockerImage }}</el-tag>
                    <el-tag v-else size="small">未构建</el-tag>
                  </div>
                </div>
                <div class="appItem" style="margin-bottom: 10px">部署版本
                  <el-button v-if="v.DockerVer != v.ClusterVer.DockerVer" size="small" @click="onSyncDockerVer(v)" type="info" style="float:left;position: absolute;margin:-2px 0 0 5px;">同步镜像</el-button>
                  <div class="appItem">
                    <el-tag v-if="v.ClusterVer.DockerImage !=''" size="small">{{ v.ClusterVer.DockerImage }}</el-tag>
                    <el-tag v-else size="small">未发布</el-tag>
                  </div>
                </div>
                <div class="appItem" style="margin-bottom: 10px">部署时间
                  <span v-if="v.ClusterVer.DockerImage !=''">{{ v.ClusterVer.DeploySuccessAt }}</span>
                  <el-tag v-else size="small">未发布</el-tag>
                </div>
                <div class="appItem" style="margin-bottom: 10px">部署角色
                  <el-tag v-if="v.DockerNodeRole=='manager'" type="danger" size="small" style="margin-left: 5px">{{ v.DockerNodeRole }}</el-tag>
                  <el-tag v-else size="small" style="margin-left: 5px">{{ v.DockerNodeRole }}</el-tag>
                  <el-button size="small" @click="onBuildAdd(v)" type="danger" style="margin-left: 5px"><el-icon><ele-SwitchButton /></el-icon>构建</el-button>
                </div>
              <div class="appItem" style="margin-bottom: 10px">日志
                <el-tooltip content="警告数量" slot="label">
                  <el-tag @click="showFsLogLevel(3,v.AppName)" v-if="v.LogWaringCount > 0" type="warning" size="small" style="margin-left: 5px;cursor: pointer" title="警告数量">{{ v.LogWaringCount }}</el-tag>
                  <el-tag v-else type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.LogWaringCount }}</el-tag>
                </el-tooltip>
                /
                <el-tooltip content="异常数量" slot="label">
                  <el-tag  @click="showFsLogLevel(4,v.AppName)" v-if="v.LogErrorCount > 0" type="danger" size="small" style="margin-left: 5px;cursor: pointer">{{ v.LogErrorCount }}</el-tag>
                  <el-tag v-else type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.LogErrorCount }}</el-tag>
                </el-tooltip>
              </div>
              <div class="appItem" style="margin-bottom: 10px">任务
                <el-tooltip content="成功数量" slot="label">
                  <el-tag @click="showTask(5,v.AppName)" v-if="v.TaskSuccessCount > 0" type="warning" size="small" style="margin-left: 5px;cursor: pointer" title="警告数量">{{ v.TaskSuccessCount }}</el-tag>
                  <el-tag v-else type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.TaskSuccessCount }}</el-tag>
                </el-tooltip>
                /
                <el-tooltip content="失败数量" slot="label">
                  <el-tag  @click="showTask(4,v.AppName)" v-if="v.TaskFailCount > 0" type="danger" size="small" style="margin-left: 5px;cursor: pointer">{{ v.TaskFailCount }}</el-tag>
                  <el-tag v-else type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.TaskFailCount }}</el-tag>
                </el-tooltip>
              </div>
            </el-card>
          </el-space>
        </el-main>
        <el-aside width="550px">
          <el-card>
            <h3 style="padding: 5px;">构建队列</h3>
            <template v-if="state.tableLogData.data.length > 0">
              <el-table  :data="state.tableLogData.data" v-loading="state.tableLogData.loading" style="width: 100%;background: #ffffff;">
                <el-table-column prop="Id" label="编号" width="70" />
                <el-table-column prop="AppName" label="应用名称" ></el-table-column>
                <el-table-column label="状态" width="90" show-overflow-tooltip>
                  <template #default="scope">
                    <el-tag v-if="scope.row.Status==0" size="small" type="info">未开始</el-tag>
                    <el-tag v-else-if="scope.row.Status==1" size="small" type="warning">构建中</el-tag>
                    <el-tag v-if="scope.row.Status==2 && scope.row.IsSuccess == true" size="small" type="success">成功</el-tag>
                    <el-tag v-else-if="scope.row.Status==2 && scope.row.IsSuccess == false" size="small" type="danger">失败</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="FinishAt" width="170" label="完成时间"></el-table-column>
                <el-table-column label="操作" width="80">
                  <template #default="scope">
                    <el-button v-if="scope.row.Status!=0" size="small" type="success" @click="showLog(scope.row)">日志</el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-pagination
                  @size-change="onHandleSizeLogChange"
                  @current-change="onHandleCurrentLogChange"
                  class="mt15"
                  :pager-count="5"
                  :page-sizes="[10, 20, 30]"
                  v-model:current-page="state.tableLogData.param.pageNum"
                  background
                  v-model:page-size="state.tableLogData.param.pageSize"
                  layout="total, sizes, prev, pager, next, jumper"
                  :total="state.tableLogData.total"
              >
              </el-pagination>
            </template>
            <el-empty v-else description="暂无数据"></el-empty>
          </el-card>
        </el-aside>
      </el-container>
    </el-card>

  <appDialog ref="appDialogRef" @refresh="getTableData()" @showOverlay="onShowOverlay()" @hideOverlay="onHideOverlay()" />
  <appAddDialog ref="appAddDialogRef" @refresh="getTableData()" @showOverlay="onShowOverlay()" @hideOverlay="onHideOverlay()" />
    <logDialog ref="logDialogRef"  />
    <taskDialog ref="taskDialogRef"  />
  <el-dialog title="构建日志" v-model="state.logDialogIsShow" style="width: 80%;height: 85%;top:20px;margin-bottom: 50px">
    <el-card shadow="hover" class="layout-padding-auto" style="background-color:#393d49;overflow: auto;">
      <pre style="color: #fff;background-color:#393d49;height: 100%;" v-html="state.logContent"></pre>
    </el-card>
  </el-dialog>

    <div v-if="state.showOverlay" class="overlay">
      <div class="overlay-content">
        <img :src="Image" style="width: 200px" alt="Image">
      </div>
    </div>
  </div>
</template>

<script setup lang="ts" name="fopsApp">

import {defineAsyncComponent, reactive, onMounted, ref, nextTick, watch, onUnmounted} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import Image from '/@/assets/loading.gif';
// var idPre = document.getElementById('idPre');
// idPre.scrollIntoView(false); // 滚动到底部

// 引入 api 请求接口
const serverApi = fopsApi();

// 引入组件
// 修改弹窗
const appDialog = defineAsyncComponent(() => import('/@/views/fops/app/dialog.vue'));
// 添加弹窗
const appAddDialog = defineAsyncComponent(() => import('/@/views/fops/app/addDialog.vue'));
// 任务组日志
const taskDialog= defineAsyncComponent(() => import('/src/views/fops/task/taskAppDialog.vue'));
// 日志
const logDialog = defineAsyncComponent(() => import('/src/views/fops/log/logV2Dialog.vue'));
const logDialogRef = ref();
// 定义变量内容
const appDialogRef = ref();
const appAddDialogRef = ref();
const taskDialogRef=ref();
const state = reactive({
  logDialogIsShow:false,
  logContent:'',
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			pageNum: 1,
			pageSize: 12,
		},
	},tableLogData: {
    data: [],
    total: 0,
    loading: false,
    param: {
      pageNum: 1,
      pageSize: 12,
    },

  },
  appName:"",
  logId:0,
  clusterId:0,
  clusterData:[],
  showOverlay:false,
  statTask:[],
});

// 初始化表格数据
const getTableData = () => {
  // 任务日志统计列表
  taskLogStat()

	state.tableData.loading = true;
	const data = [];
  // 请求接口
  serverApi.appsList({}).then(function (res){
    if (res.Status){
      for (let i = 0; i < res.Data.length; i++) {
        var item=res.Data[i]
        var taskFailCount=state.statTask.filter(t=>t.Status==4&&t.ClientName==item.AppName)
        var taskSuccessCount=state.statTask.filter(t=>t.Status==5&&t.ClientName==item.AppName)
        if(taskFailCount.length>0)
        {
          item.TaskFailCount=taskFailCount[0].Count
        }else{
          item.TaskFailCount=0
        }
        if(taskSuccessCount.length>0)
        {
          item.TaskSuccessCount=taskSuccessCount[0].Count
        }else{
          item.TaskSuccessCount=0
        }
        data.push(item)
      }
      state.tableData.data =data;
      state.tableData.total = data.length;
      state.tableData.loading = false;


    }else{
      state.tableData.data=[]
      state.tableData.loading = false;
    }
  })
};





const getTableLogData = () => {
  state.tableLogData.loading = true;
  const data = {
    appName:"",//state.appName
    pageIndex:state.tableLogData.param.pageNum,
    pageSize:state.tableLogData.param.pageSize,
  };
  // 请求接口
  serverApi.buildList(data).then(function (res){
    if (res.Status){
      state.tableLogData.data = res.Data.List;
      state.tableLogData.total = res.Data.RecordCount;
    }else{
      state.tableLogData.data=[]
    }
    state.tableLogData.loading = false;
  })
};
const getTableClusterData = () => {
  state.tableData.loading = true;
  const data = [];
  // 请求接口
  serverApi.clusterList({}).then(function (res){
    if (res.Status){
      var lst=[]
      for (let i = 0; i < res.Data.length; i++) {
        var item=res.Data[i]
        if (i==0){
          state.clusterId=item.Id;
        }
        item.Name=item.Name+" - "+item.DockerName
        lst.push(item)
      }
      state.clusterData = lst;
    }else{
      state.tableData.data=[]
    }
  })
};

// 打开FS日志
const showFsLogLevel=(level:any,appName:any)=>{
  logDialogRef.value.openDialogLogLevel(level,appName);
}
// 任务组日志
const showTask=(st:any,appName:any)=>{
  taskDialogRef.value.openDialogApp(st,appName);
}

const onClusterChange=(value:number)=>{
  state.clusterId=value
}
// 打开新增用户弹窗
const onOpenAdd = (type: string) => {
  appAddDialogRef.value.openDialog(type,null);
};

// 打开修改用户弹窗
const onOpenEdit = (type: string, row: any) => {
  appDialogRef.value.openDialog(type, row);
};


// 清除镜像
const onClearDockerImage = () => {
  ElMessageBox.confirm(`此操作将永久清除：“None镜像”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        state.showOverlay=true
        // 删除逻辑
        serverApi.dockerClearImage().then(function (res){
          if (res.Status){
            ElMessage.success('清除成功');
          }else{
            ElMessageBox.alert(res.StatusMessage,'Warning',{ type: 'warning',dangerouslyUseHTMLString: true})
          }
          state.showOverlay=false
        })
      })
      .catch(() => {});
};

// 分页改变
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pageSize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.pageNum = val;
	getTableData();
};
const onHandleSizeLogChange = (val: number) => {
  state.tableLogData.param.pageSize = val;
  getTableLogData();
};
const onShowBuildList=(row: any)=>{
  state.appName=row.AppName
  state.tableLogData.param.pageNum=1
  state.tableLogData.param.pageSize=10
  getTableLogData();
}
// 分页改变
const onHandleCurrentLogChange = (val: number) => {
  state.tableLogData.param.pageNum = val;
  getTableLogData();
};
// 定义定时器
let intervalId = null;
// 使用 watch 监听 state 中 count 属性的变化
watch(() => state.logDialogIsShow, (newValue, oldValue) => {
  //console.log(`count 从 ${oldValue} 变为 ${newValue}`);
  if(!newValue){
    clearInterval(intervalId);
  }else {
    intervalId = setInterval(onShowLog, 1000);
  }
});

const showLog=(row:any)=>{
  state.logId=row.Id
  serverApi.buildLog(state.logId.toString()).then(function (res){
    state.logContent=res
    state.logDialogIsShow=true
  })
}
const onShowLog=()=>{
  serverApi.buildLog(state.logId.toString()).then(function (res){
    state.logContent=res
  })
}
const onShowOverlay=()=>{
  state.showOverlay=true
}
const onHideOverlay=()=>{
  state.showOverlay=false
}
// 构建
const onBuildAdd = (row:any) => {
  ElMessageBox.confirm(`请确认是否添加构建?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 提交数据
        var param={
          "AppName":row.AppName,
          "ClusterId":state.clusterId,
        }
        serverApi.buildAdd(param).then(async function(res){
          if(res.Status){
            ElMessage.success("添加成功")
            // 刷新构建日志
            getTableLogData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
};

// 重启容器
const onRestartDocker = (row:any) => {
  ElMessageBox.confirm(`请确认是否重启容器?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        state.showOverlay=true
        // 提交数据
        var param={
          "AppName":row.AppName,
          "ClusterId":state.clusterId,
        }
        serverApi.restartDocker(param).then(async function(res){
          state.showOverlay=false
          if(res.Status){
            ElMessage.success("重启成功")
            // 刷新应用界面
            getTableData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {
        state.showOverlay=false});
};
// 同步版本
const onSyncDockerVer = (row:any) => {
  ElMessageBox.confirm(`请确认是否要同步仓库镜像到集群中?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        state.showOverlay=true
        // 提交数据
        var param={
          "appName":row.AppName,
          "clusterId":state.clusterId,
        }
        serverApi.syncDockerImage(param).then(async function(res){
          state.showOverlay=false
          if(res.Status){
            ElMessage.success("同步成功")
            // 刷新构建日志
            getTableData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {
        state.showOverlay=false});
};
// 全部构建
const onAllBuild=()=>{
  ElMessageBox.confirm(`请确认是否构建全部应用?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        for (let i = 0; i < state.tableData.data.length; i++) {
          var item=state.tableData.data[i]
          onBuildAddFunc(item)
        }
      })
      .catch(() => {});
}
const onBuildAddFunc = (row:any) => {
    // 提交数据
    var param={
      "AppName":row.AppName,
      "ClusterId":state.clusterId,
    }
    serverApi.buildAdd(param).then(async function(res){
      if(res.Status){
        ElMessage.success("添加成功")
        // 刷新构建日志
        getTableLogData()
      }else{
        ElMessage.error(res.StatusMessage)
      }
    })
};
const getGitArray=(lst:[])=>{
  var array=[]
  for (let i = 0; i < lst.length; i++) {
    serverApi.gitInfo({"gitId":lst[i]}).then(function (res){
      if (res.Status){
        array.push(res.Data.Name)
      }
    })
  }
  return array
}
const getGit=(val:number)=>{
  var array=[]
  serverApi.gitInfo({"gitId":val}).then(function (res){
      if (res.Status){
        array.push(res.Data.Name)
      }
    })
  return array
}

// 任务日志统计列表
const taskLogStat=()=>{
  serverApi.taskStatList("").then(function (res){
    if (res.Status){
      state.statTask=res.Data
    }
  })
}

let intervalLogId = null;
let intervalAppId = null;
// 页面加载时
onMounted(() => {
	getTableData();
  getTableLogData();
  getTableClusterData();
  intervalLogId = setInterval(getTableLogData, 3000);
  intervalAppId = setInterval(getTableData, 3000);
});
// 页面注销的时候
onUnmounted(()=>{
  clearInterval(intervalLogId);
  clearInterval(intervalAppId);
})
</script>

<style lang="scss">
.system-user-container {
	:deep(.el-card__body) {
		display: flex;
		flex-direction: column;
		flex: 1;
		overflow: auto;
		.el-table {
			flex: 1;
		}
	}
}
.flex-warp {
  display: flex;
  flex-wrap: wrap;
  align-content: flex-start;
  margin: 0 -5px;
  .flex-warp-item {
    padding: 5px;
    width: 298px;
    min-height: 360px;
    //border: #666 1px solid;
    .flex-warp-item-box {
      border: 1px solid var(--next-border-color-light);
      width: 100%;
      height: 100%;
      border-radius: 2px;
      display: flex;
      flex-direction: column;
      transition: all 0.3s ease;

      .item-img {
        width: 100%;
        height: 215px;
        overflow: hidden;
        img {
          transition: all 0.3s ease;
          width: 100%;
          height: 100%;
        }
      }
      .item-txt {
        flex: 1;
        padding: 15px;
        display: flex;
        flex-direction: column;
        overflow: hidden;
        .item-txt-title {
          margin: 10px!important;
          text-overflow: ellipsis;
          overflow: hidden;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          display: -webkit-box;
          color: #666666;
          transition: all 0.3s ease;
          &:hover {
            color: var(--el-color-primary);
            text-decoration: underline;
            transition: all 0.3s ease;
          }
        }
        .item-txt-other {
          flex: 1;
          align-items: flex-end;
          display: flex;
          .item-txt-msg {
            font-size: 12px;
            color: #8d8d91;
          }
          .item-txt-price {
            display: flex;
            justify-content: space-between;
            align-items: center;
            .font-price {
              color: #ff5000;
              .font {
                font-size: 22px;
              }
            }
          }
        }
      }
    }
  }
}
.appItem{
  margin: 10px;
}
.el-row{
  margin: 0!important;
  display: block!important;

}
.flex-warp {
  width: 100%!important;
}
.flex-warp-item{
  float: left;
  margin: 5px;
}
.el-dialog__body{
  height: 100%!important;
}

.el-card__header{
  background-color: #545c64;
}
.el-space__item .el-card__body{
  background-color: #f9f9e3;
}
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10000;
}

.overlay-content {
  text-align: center;
  color: white;
}
</style>
