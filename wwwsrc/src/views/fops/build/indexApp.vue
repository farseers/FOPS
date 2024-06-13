<template>
  <div class="layout-padding" style="position: relative;">
    <el-card shadow="hover">
      <el-header style="padding: 0">
        <el-select v-model="state.clusterId" placeholder="请选择集群" class="ml10" @change="onClusterChange" style="width: 250px;">
          <el-option v-for="item in state.clusterData" :key="item.Id" :label="item.Name" :value="item.Id"></el-option>
        </el-select>
        <el-button size="default" type="success" class="ml10" @click="onOpenAdd('add')"><el-icon><ele-FolderAdd /></el-icon>新增应用</el-button>
        <el-button size="default" type="info" class="ml10" @click="onClearDockerImage('add')"><el-icon><ele-Delete /></el-icon>清除None镜像</el-button>
<!--        <el-button size="default" type="warning" class="ml10" @click="onAllBuild()"><el-icon><ele-SwitchButton /></el-icon>全部构建</el-button>-->
        <el-button size="default" type="danger" class="ml10" @click="onStopBuild()"><el-icon><ele-SwitchButton /></el-icon>停止构建</el-button>
      </el-header>
      <!--应用列表-->
      <el-container>
        <el-main style="padding: 0;overflow: hidden;">
          <el-space wrap style="align-items: unset;">
            <el-card shadow="hover" v-for="(v, k) in state.tableData.data" :key="k" style="width: 280px;">
              <template #header>
                <div class="card-header" style="height: 20px;">
                  <el-tag size="default" @click="onOpenEdit('edit', v)" style="cursor: pointer;">{{ v.AppName }}</el-tag>
                  <el-button size="small" type="warning" @click="onRestartDocker(v)" style="float:right;position: relative;"><el-icon><ele-SwitchButton /></el-icon>重启</el-button>
                  <el-tooltip content="实例数量/副本数量" slot="label">
                    <el-tag v-if="v.IsHealth" size="small" style="margin-left: 5px">{{v.DockerInstances}}/{{ v.DockerReplicas }}</el-tag>
                    <el-tag v-else size="small" type="danger" style="margin-left: 5px">{{v.DockerInstances}}/{{ v.DockerReplicas }}</el-tag>
                  </el-tooltip>
                </div>
              </template>
                <el-button size="small" type="success" @click="showFsLogLevel(2,v.AppName)" style="float:right;position: relative;margin-left: 5px">应用日志</el-button>
                <el-button size="small" type="primary" @click="showDockerLog(v.AppName)" style="float:right;position: relative;margin-left: 5px">容器日志</el-button>

                <div class="appItem" style="margin-bottom: 10px">仓库版本
                  <div class="appItem">
                    <el-tag v-if="v.DockerImage !=''" size="small">{{ v.DockerImage }}</el-tag>
                    <el-tag v-else size="small">未构建</el-tag>
                  </div>
                </div>
                <div class="appItem" style="margin-bottom: 10px">部署版本
                  <el-button v-if="v.DockerImage != v.ClusterVer.DockerImage" size="small" @click="onSyncDockerVer(v)" type="info" style="float:left;position: absolute;margin:-2px 0 0 5px;">同步镜像</el-button>
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
                </div>
              <div class="appItem" style="margin-bottom: 10px">应用日志
                <el-tooltip content="警告数量" slot="label">
                  <el-tag v-if="v.LogWaringCount > 0" @click="showFsLogLevel(3,v.AppName)" type="warning" size="small" style="margin-left: 5px;cursor: pointer">{{ v.LogWaringCount }}</el-tag>
                  <el-tag v-else @click="showFsLogLevel(3,v.AppName)" type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.LogWaringCount }}</el-tag>
                </el-tooltip>
                /
                <el-tooltip content="异常数量" slot="label">
                  <el-tag v-if="v.LogErrorCount > 0" @click="showFsLogLevel(4,v.AppName)" type="danger" size="small" style="margin-left: 5px;cursor: pointer">{{ v.LogErrorCount }}</el-tag>
                  <el-tag v-else @click="showFsLogLevel(4,v.AppName)" type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.LogErrorCount }}</el-tag>
                </el-tooltip>
              </div>
              <div class="appItem" style="margin-bottom: 10px">调度任务
                <el-tooltip content="成功数量" slot="label">
                  <el-tag v-if="v.TaskSuccessCount > 0" @click="showTask(2,v.AppName)" type="success" size="small" style="margin-left: 5px;cursor: pointer">{{ v.TaskSuccessCount }}</el-tag>
                  <el-tag v-else @click="showTask(2,v.AppName)" type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.TaskSuccessCount }}</el-tag>
                </el-tooltip>
                /
                <el-tooltip content="失败数量" slot="label">
                  <el-tag v-if="v.TaskFailCount > 0" @click="showTask(3,v.AppName)" type="danger" size="small" style="margin-left: 5px;cursor: pointer">{{ v.TaskFailCount }}</el-tag>
                  <el-tag v-else @click="showTask(3,v.AppName)" type="info" size="small" style="margin-left: 5px;cursor: pointer">{{ v.TaskFailCount }}</el-tag>
                </el-tooltip>
              </div>
              <el-button v-if="v.AppGit > 0" size="small" @click="onSyncWorkflows(v)" type="info" style="margin-left: 5px;width:100%"><el-icon><ele-SwitchButton /></el-icon>刷新工作流</el-button>
              <div class="appItem" style="margin-bottom: 10px">构建
                <el-button v-if="v.AppGit > 0" v-for="(item, index) in v.WorkflowsNames" size="small" @click="onBuildAdd(v,item)" type="danger" style="margin-left: 5px;margin-bottom: 5px;">{{item}}</el-button>
              </div>
            </el-card>
          </el-space>
        </el-main>
        <el-aside width="550px">
          <el-card>
            <h3 style="padding: 5px;">构建队列</h3>
            <template v-if="state.tableLogData.data.length > 0">
              <el-table  :data="state.tableLogData.data" v-loading="state.tableLogData.loading" style="width: 100%;background: #ffffff;" :cell-style="{padding:'2px 0'}">
                <el-table-column prop="FinishAt" width="170" label="构建时间"></el-table-column>
                <el-table-column label="应用名称" show-overflow-tooltip>
                  <template #default="scope">
                    <el-tag v-if="scope.row.Status==0" size="small" type="info">未开始</el-tag>
                    <el-tag v-else-if="scope.row.Status==1" size="small" type="warning">构建中</el-tag>
                    <el-tag v-if="scope.row.Status==2 && scope.row.IsSuccess == true" size="small" type="success">成功</el-tag>
                    <el-tag v-else-if="scope.row.Status==2 && scope.row.IsSuccess == false" size="small" type="danger">失败</el-tag>
                    <el-tag size="small" type="info">{{ scope.row.WorkflowsName }}</el-tag>
                    <span style="margin-left: 5px ">{{ scope.row.AppName }}:{{ scope.row.BuildNumber }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="80">
                  <template #default="scope">
                    <el-button v-if="scope.row.Status!=0" size="small" type="success" @click="showBuildLog(scope.row)">日志</el-button>
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
  <el-dialog title="构建日志" v-model="state.isShowBuildLogDialog" style="width: 80%;top:20px;margin-bottom: 50px;">
    <el-checkbox v-model="state.autoLog" style="margin-bottom: 5px;">自动刷新日志</el-checkbox>
    <div class="layout-padding-auto" style="background-color:#393d49;">
      <div ref="scrollableBuildLog" style="height: 100%;overflow-y: auto;">
        <pre style="color: #fff;background-color:#393d49;padding: 5px 0 5px 5px;" v-html="state.buildLogContent"></pre>
      </div>
    </div>
  </el-dialog>
  <el-dialog title="容器日志" v-model="state.isShowDockerLogDialog" style="width: 80%;top:20px;margin-bottom: 50px;">
      <div>
        <el-tag size="default" :type="item.Id==state.dockerLog.Id?'':'info'"
        @click="clickDockerLog(item)"
        v-for="item in state.dockerLogContent" 
        :key="item.Id" 
        style="cursor: pointer;margin:0 15px 5px 0">
        {{ item.Name }}（{{ item.Node }}）
      </el-tag>
      </div>
      <div style="margin: 5px 0;">
            <div> 
<!--            <span style="display: inline-block;margin-right:10px">{{ state.dockerLog.Name }} </span>-->
<!--            <span style="display: inline-block;margin-right:10px">{{ state.dockerLog.Node }} </span>-->
              <el-tag size="default">{{ state.dockerLog.State }}</el-tag>
              <el-tag size="default">{{ state.dockerLog.StateInfo }}</el-tag>
              <el-tag size="default">{{ state.dockerLog.Image }}</el-tag>
<!--            <span style="display: inline-block;margin-right:10px">{{ state.dockerLog.State }}</span>-->
<!--            <span style="display: inline-block;margin-right:10px">{{ state.dockerLog.StateInfo }}</span>-->
<!--            <span style="display: inline-block;margin-right:10px">{{ state.dockerLog.Image }}</span>-->
          </div>
          <div style="color: #f56c6c;">{{ state.dockerLog.Error }}</div>
      </div>
    <div class="layout-padding-auto" style="background-color:#393d49;">
      <div ref="scrollableDockerLog" style="height: 100%;overflow-y: auto;">
        <pre v-html="state.dockerLog.Log" style="color: #fff;background-color:#393d49;padding: 5px 0 5px 5px;" ></pre>
      </div>
    </div>
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
const appDialog = defineAsyncComponent(() => import('/src/views/fops/build/editAppDialog.vue'));
// 添加弹窗
const appAddDialog = defineAsyncComponent(() => import('/src/views/fops/build/addAppDialog.vue'));
// 任务组日志
const taskDialog= defineAsyncComponent(() => import('/src/views/fops/task/taskAppDialog.vue'));
// 日志
const logDialog = defineAsyncComponent(() => import('/src/views/fops/log/logV2Dialog.vue'));

const logDialogRef = ref();
// 定义变量内容
const appDialogRef = ref();
const appAddDialogRef = ref();
const taskDialogRef = ref();
const scrollableBuildLog = ref();
const scrollableDockerLog = ref();
const state = reactive({
  isShowBuildLogDialog: false,
  isShowDockerLogDialog: false,
  dockerLogContent: [],//容器日志
  dockerLog:{
    Id:'',
    Name:'',  Node:'',  State:'',  StateInfo:'',  Error:'',  Image:'',
  },//容器日志选中
  buildLogContent: '',
  buildLogContents: '',
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
      pageSize: 20,
    },

  },
  appName:"",
  buildLogId:0,
  clusterId:0,
  clusterData:[],
  showOverlay:false,
  statTask:[],
  autoLog:true,
});


// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	const data = [];
  var param={
    "ClusterId" : state.clusterId,
  }
  // 获取应用列表
  serverApi.appsList(param).then(function (res){
    if (res.Status){
      for (let i = 0; i < res.Data.length; i++) {
        let item = res.Data[i];
        item.TaskFailCount=0
        item.TaskSuccessCount=0
        let taskFailCount = state.statTask.filter(t => t.ExecuteStatus == 3 && t.ClientName == item.AppName);

        if(taskFailCount.length > 0) {
          item.TaskFailCount=taskFailCount[0].Count
        }

        let taskSuccessCount = state.statTask.filter(t => t.ExecuteStatus == 2 && t.ClientName == item.AppName);
        if(taskSuccessCount.length > 0) {
          item.TaskSuccessCount=taskSuccessCount[0].Count
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
  // 请求接口
  serverApi.clusterList({}).then(function (res){
    if (res.Status){
      var lst=[]
      for (let i = 0; i < res.Data.length; i++) {
        var item=res.Data[i]
        if (i==0) {
          state.clusterId=item.Id;
          getTableData();
        }
        item.Name=item.Name+" - "+item.FopsAddr
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
  getTableData()
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

const onHandleSizeLogChange = (val: number) => {
  state.tableLogData.param.pageSize = val;
  getTableLogData();
};

// 分页改变
const onHandleCurrentLogChange = (val: number) => {
  state.tableLogData.param.pageNum = val;
  getTableLogData();
};
// 定义定时器
let intervalId = null;
// 使用 watch 监听 state 中 count 属性的变化
watch(() => state.isShowBuildLogDialog, (newValue, oldValue) => {
  if(!newValue){
    clearInterval(intervalId);
  }else {
    intervalId = setInterval(onShowLog, 500);
  }
});

// 显示构建日志
const showBuildLog=(row:any)=>{
  state.buildLogId=row.Id
  serverApi.buildLog(state.buildLogId.toString()).then(function (res){
    state.buildLogContent = res
    state.isShowBuildLogDialog=true
  })
}
//点击容器日志选项
const clickDockerLog = (item:any)=>{
  state.dockerLog = item
}
// 显示容器日志
const showDockerLog=(appName:string)=>{
  serverApi.dockerLog({ "AppName": appName, "tailCount": 100 }).then(function (res){
    state.dockerLogContent = res.Data;
    if(state.dockerLogContent && state.dockerLogContent.length>0){
      clickDockerLog(state.dockerLogContent[0])
    }
    state.isShowDockerLogDialog = true
    setTimeout(()=>{   //自动跳到底部 
      scrollableDockerLog.value.scrollTop = scrollableDockerLog.value.scrollHeight;
    },500)
  })
}

const onShowLog=()=>{
  serverApi.buildLog(state.buildLogId.toString()).then(function (res) {
    // 如果从接口获取到的内容与本地内容一样时，则不用滚动
    let isChange= state.buildLogContents != res;
    state.buildLogContent = res;
    state.buildLogContents = res;
    // 自动刷新日志
    if (state.autoLog && isChange) {
        scrollableBuildLog.value.scrollTop = scrollableBuildLog.value.scrollHeight;
    }
  })
}

const onShowOverlay=()=>{
  state.showOverlay=true
}
const onHideOverlay=()=>{
  state.showOverlay=false
}
// 构建
const onBuildAdd = (row:any,workflowsName:any) => {
  ElMessageBox.confirm(`请确认是否添加构建`+ workflowsName +`?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 提交数据
        var param={
          "AppName" : row.AppName,
          "ClusterId" : state.clusterId,
          "WorkflowsName" : workflowsName,
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

// 刷新工作流文件
const onSyncWorkflows = (row:any) => {
  state.showOverlay=true
  let param = {
    "AppName": row.AppName
  };
  serverApi.syncWorkflows(param).then(async function(res){
    if (res.Status) {
      ElMessage.success("刷新成功")
    } else {
      ElMessage.error(res.StatusMessage)
    }
    // 刷新
    getTableData()
  })
  state.showOverlay=false
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
// 停止构建
const onStopBuild=()=>{
  ElMessageBox.confirm(`请确认是否停止构建?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 提交数据
        var param={ }
        serverApi.buildStop(param).then(async function(res){
          if(res.Status){
            ElMessage.success("成功停止")
            // 刷新构建日志
            getTableLogData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
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
let statCountAppId = null;
// 页面加载时
onMounted(() => {
  getTableClusterData();
  getTableLogData();
  taskLogStat();
  intervalAppId = setInterval(getTableData, 3000);
  intervalLogId = setInterval(getTableLogData, 3000);
  statCountAppId = setInterval(taskLogStat, 10000);
});
// 页面注销的时候
onUnmounted(()=>{
  clearInterval(intervalLogId);
  clearInterval(intervalAppId);
  clearInterval(statCountAppId);
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
  display: flex;
  flex-direction: column;
}
.layout-container .layout-padding-auto {
  flex: 1;
  overflow: auto;
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