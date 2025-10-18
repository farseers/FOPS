<template>
    <div class="w100">
        <div class="conlyRow">
            <div v-for="item, index in state.tableData" :key="index.toString() + 'ecfy'" class="conlyCol">
                <el-card :class="item.IsHealth ? 'conlyCard' : 'conlyCard conly_w'">
                    <div class="name">
                        <el-tag size="default">{{ item.AppName }}</el-tag>
                        <el-tooltip content="实例数量/副本数量" slot="label">
                            <el-tag  @click="showDockerTag(item,1)" v-if="item.IsHealth" size="small" style="margin-left: 5px;cursor: pointer;">{{ item.DockerInstances }}/{{ item.DockerReplicas }}</el-tag>
                            <el-tag  @click="showDockerTag(item,2)" v-else size="small" type="danger" style="margin-left: 5px;cursor: pointer;">{{ item.DockerInstances }}/{{ item.DockerReplicas }}</el-tag>
                        </el-tooltip>
                    </div>
                    <div style="display: flex;">
                      <el-tooltip content="删除服务" slot="label">
                        <el-icon style="margin-left: 12px;cursor: pointer;color: #f56c6c;font-size: 18px" @click="onDeleteDocker(v)"><ele-CircleCloseFilled /></el-icon>
                      </el-tooltip>
                      <el-tooltip content="重启服务" slot="label" v-if="item.DockerReplicas > 0">
                          <el-icon style="margin-left: 20px;cursor: pointer;color: #F56C6C;font-size: 18px;" @click="onRestartDocker(item)"><ele-Refresh /></el-icon>
                      </el-tooltip>
                      <el-tooltip content="容器日志" slot="label" v-if="item.DockerReplicas > 0">
                          <el-icon style="margin-left: 20px;cursor: pointer;color: #409EFF;font-size: 18px;"  @click="showDockerLog(item.AppName)"><ele-Reading /></el-icon>
                       </el-tooltip>
                      <el-tooltip content="应用日志" slot="label" v-if="item.IsSys === false">
                          <el-icon style="margin-left: 20px;cursor: pointer;color: #409EFF;font-size: 18px;" @click="showFsLogLevel(2, item.AppName)"><ele-Document /></el-icon>
                      </el-tooltip>
                    </div>
                    <div>应用日志
                        <el-tooltip content="警告数量" slot="label">
                            <el-tag v-if="item.LogWaringCount > 0" @click="showFsLogLevel(3, item.AppName)"
                                type="warning" size="small" style="margin-left: 5px;cursor: pointer">{{ item.LogWaringCount }}</el-tag>
                            <el-tag v-else @click="showFsLogLevel(3, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.LogWaringCount }}</el-tag>
                        </el-tooltip>
                        /
                        <el-tooltip content="异常数量" slot="label">
                            <el-tag v-if="item.LogErrorCount > 0" @click="showFsLogLevel(4, item.AppName)" type="danger"
                                size="small" style="margin-left: 5px;cursor: pointer">{{ item.LogErrorCount }}</el-tag>
                            <el-tag v-else @click="showFsLogLevel(4, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.LogErrorCount }}</el-tag>
                        </el-tooltip>
                    </div>
                    <div>调度任务
                        <el-tooltip content="成功数量" slot="label">
                            <el-tag v-if="item.TaskSuccessCount > 0" @click="showTask(2, item.AppName)" type="success"
                                size="small" style="margin-left: 5px;cursor: pointer">{{ item.TaskSuccessCount
                                }}</el-tag>
                            <el-tag v-else @click="showTask(2, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.TaskSuccessCount }}</el-tag>
                        </el-tooltip>
                        /
                        <el-tooltip content="失败数量" slot="label">
                            <el-tag v-if="item.TaskFailCount > 0" @click="showTask(3, item.AppName)" type="danger"
                                size="small" style="margin-left: 5px;cursor: pointer">{{ item.TaskFailCount }}</el-tag>
                            <el-tag v-else @click="showTask(3, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.TaskFailCount }}</el-tag>
                        </el-tooltip>
                    </div>
                  <div class="progress_cs">
                    <el-tag type="info" size="small">CPU</el-tag>
                   <span class="progress_sp">
                            <el-progress 
                            :text-inside="true" 
                            class="custom-progress" 
                            :stroke-width="18" 
                            :color="state.customColorsCpu" 
                            :percentage="item.CpuUsagePercent">
                             <span>{{item.CpuUsagePercent || 0 }}%</span>
                        </el-progress>
                       </span>
                   </div>
                  <div class="progress_cs">
                    <el-tag type="info" size="small">内存</el-tag>
                     <span class="progress_sp">
                            <el-progress 
                            :text-inside="true" 
                            class="custom-progress" 
                            :stroke-width="18" 
                            :color="state.customColors" 
                            :percentage="item.MemoryUsagePercent">
                            <span>{{item.MemoryUsagePercent || 0 }}% {{ item.MemoryUsage }}MB</span>
                        </el-progress>
                       </span> 
                      </div>
                </el-card>
            </div>
        </div>
        <div v-if="state.showOverlay" class="overlay">
    <div class="overlay-content">
      <img :src="Image" style="width: 200px" alt="Image">
    </div>
  </div>
            <dockerDialog ref="dockerDialogRef"/>
            <taskDialog ref="taskDialogRef"  />
            <logDialog ref="logDialogRef"  />
            <editAppNum ref="editAppNumRef" @refresh="getData()"/>
    </div>
</template>
<script setup name="Ecfy">
import { reactive, onMounted, defineExpose, ref,defineAsyncComponent } from 'vue';
import { ElMessage,ElMessageBox } from 'element-plus';
import { fopsApi } from "/@/api/fops";
import Image from '/@/assets/loading.gif';
const dockerDialog = defineAsyncComponent(() => import('/src/views/fops/task/dockerDialog.vue'));
const logDialog = defineAsyncComponent(() => import('/src/views/fops/log/logV2Dialog.vue'));
const taskDialog= defineAsyncComponent(() => import('/src/views/fops/task/taskAppDialog.vue'));
const editAppNum = defineAsyncComponent(() => import('/src/views/fops/build/editAppNum.vue'));
// 引入 api 请求接口
const serverApi = fopsApi();
const conlyTabs = ref(null)
const logDialogRef = ref();
const taskDialogRef = ref();
const dockerDialogRef = ref();
const editAppNumRef = ref();
// 定义变量内容
const state = reactive({
    tableData: [],
    statTask: [],
    showOverlay:false,
    isShowDockerLogDialog: false, //容器日志
	dockerLogContent: [],//容器日志
    clusterId:0,
	dockerLog: {
		Id: '',
		Name: '', Node: '', State: '', StateInfo: '', Error: '', Image: '',
	},//容器日志选中
   customColors:[
        {color: '#5cb87a', percentage: 50},
        {color: '#e6a23c', percentage: 80},
        {color: '#f56c6c', percentage: 100},
    ],
     customColorsCpu:[
        {color: '#5cb87a', percentage: 200},
        {color: '#e6a23c', percentage: 300},
        {color: '#f56c6c', percentage: 400},
    ]
});
const showDockerTag = (row,type) =>{
    // console.log(row,row.DockerReplicas,type)
  editAppNumRef.value.openDialog(row,type);
}
const showDockerLog = (AppName) => {
    dockerDialogRef.value.openDockerLog(AppName);
}
// 打开FS日志
const showFsLogLevel=(level,appName)=>{
  logDialogRef.value.openDialogLogLevel(level,appName,'应用日志');
}
// 任务组日志
const showTask=(st,appName)=>{
  taskDialogRef.value.openDialogApp(st,appName);
}
const  format =(c,t,s)=> {
  console.log(c,t,s)
        return  `${c}%`;
      }
const getData = () => {
    var param = {
        "ClusterId": state.clusterId,
        "IsSys": true,
    }
    // 获取应用列表
    serverApi.appsList(param).then(function (res) {
        if (res.Status) {
            state.tableData = res.Data;
        } else {
            ElMessage.warning(res.StatusMessage);
        }
    })
}

// 删除服务
const onDeleteDocker = (row) => {
  ElMessageBox.confirm(`请确认是否删除服务?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        state.showOverlay=true
        // 提交数据
        var param={
          "AppName" : row.AppName,
        }
        serverApi.appsServiceDel(param).then(async function(res){
          state.showOverlay=false
          if(res.Status){
            ElMessage.success("删除服务成功")
            // 刷新应用界面
            getTableData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        }).catch(() => {
          state.showOverlay=false});
      })
      .catch(() => {
        state.showOverlay=false});
};

const onRestartDocker = (row) => {
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
            getData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        }).catch(() => {
            state.showOverlay=false});
      })
      .catch(() => {
        state.showOverlay=false});
};
onMounted(() => {
    getData()
});
defineExpose({
    getData,
});
</script>
<style>
.progress_cs{
    display: flex;
    align-items: center;
}
.progress_sp{
    flex: 1;
    padding: 0 2px;
}
.custom-progress{
    padding: 0 2px;
}
.custom-progress .el-progress-bar__innerText {
  color: #000 !important; /* 例如设置为红色 */
}
.progress_cs .el-progress-bar__inner{
    max-width: 100%;
}
</style>
<style scoped lang="scss">
.el-dialog__body {
		display: flex;
		flex-direction: column;
	}

	.layout-container .layout-padding-auto {
		flex: 1;
		overflow: auto;
	}

.conlyRow {
    flex-wrap: wrap;
    display: flex !important;
    min-height: 170px;
    line-height: 15px;
}

.conlyCol {
    padding: 5px;
    box-sizing: border-box;
    width: 190px;
}


.conlyCard {
    background-color: #f9f9e3;
    //border: 1px dotted var(--el-color-primary);
    :deep(.el-card__body) {
        padding: 10px 5px;
	}
    .el-card__body>div {
    text-align: left;
    margin-bottom: 5px;
    font-size: 12px;
}
    .name {
        font-weight: 700;
    }
}

.conly_w {
    background: var(--el-color-danger-light-8);
    border: 1px dotted var(--el-color-danger);
}

.layout-container .layout-padding-auto {
    flex: 1;
    overflow: auto;
}
.ecfy_btn{
    --el-button-size: 20px;
    padding: 3px 9px;
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