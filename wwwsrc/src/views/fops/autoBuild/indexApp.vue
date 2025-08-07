<template>
  <div class="layout-padding" style="position: relative;">
    <el-card shadow="hover">
      <!--应用列表-->
      <el-container>
        <el-main class="buildMain">
          <div :class="buildItemCls(v)" v-for="(v, k) in state.tableData.data" :key="k" >
            <el-card shadow="hover" style="flex: 1;" >
              <div class="title">{{ v.AppName }}</div>
              <div class="_divs">
                <ul :class="v.List.length>1?'_uls _uls1':'_uls'" v-for="(item, index) in v.List" :key="item.CommitId">
                  <li>
                    <el-tag size="small" style="margin-right: 3px;">{{ item.BranchName }}</el-tag>
                    <span style="margin-right: 3px;">git提交时间：{{ item.CommitAt }}</span>
                    <span>
                     自动构建： 
                  <el-switch
                    v-model="item.AutoBuild"
                    size="small"
                    width="50"
                    @change="setAutoBuild(item)"
                    inline-prompt
                    style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;padding: 2px;"
                    active-text="开启"
                    inactive-text="关闭"
                  />
                   </span>
                  </li>
                  <li v-if="item.CommitMessage"><span>git提交信息：{{ item.CommitMessage }}</span></li>
                  <li>
                    <el-tag size="small" style="margin-right: 3px;" type="success" v-show="item.BuildSuccess">成功</el-tag>
                    <el-tag size="small" style="margin-right: 3px;" type="danger" v-show="!item.BuildSuccess">失败</el-tag>
                    <span style="margin-right: 3px;">构建ID:{{ item.BuildId }}</span> 
                    <span>
                      失败次数:
                    <el-tag size="small" style="margin-right: 3px;" :type="item.BuildErrorCount > 0?'danger':'info'" >{{item.BuildErrorCount}}</el-tag>
                    </span>
                  
                  </li>
                  <li style="display: flex;align-items: center;padding: 5px 10px;">构建时间：{{ item.BuildAt }}</li>
                  <!-- v-if="item.BuildErrorCount == 3 && !item.BuildSuccess" -->
                  <li style="height: 20px;padding: 0;display: flex;align-items: center;">
                  <el-button  
                  v-show="item.BuildErrorCount >= 3 || !item.AutoBuild" 
                  size="small" @click="goBuild(item)" 
                  type="info" 
                  style="width:100%">
                  <el-icon class="iconfont icon-wenducanshu-05"></el-icon>
                  构建
                </el-button>
                </li>
                </ul>
              </div>
              
            </el-card>
          </div>
        
        </el-main>
        <el-aside width="460px">
          <el-card>

      <el-header style="padding: 0;--el-header-height:40px;float: right">
        <el-button size="small" type="info" class="ml10" @click="onClearDockerImage('add')"><el-icon><ele-Delete /></el-icon>清除镜像</el-button>
        <el-button size="small" type="danger" class="ml10" @click="onStopBuild(0)" style=""><el-icon><ele-SwitchButton /></el-icon>停止构建</el-button>
      </el-header>
            <h3 style="padding: 5px;">构建队列 </h3>
            <template v-if="state.tableLogData.data.length > 0">
              <el-table  :data="state.tableLogData.data" v-loading="state.tableLogData.loading" style="width: 100%;background: #ffffff;" :cell-style="{padding:'2px 0'}">
                <el-table-column prop="FinishAt" width="130" label="构建时间"></el-table-column>
                <el-table-column label="应用名称" show-overflow-tooltip>
                  <template #default="scope">
                    <el-tag v-if="scope.row.Status==0" size="small" type="info">未开始</el-tag>
                    <el-tag v-else-if="scope.row.Status==1" size="small" type="warning">构建中</el-tag>
                    <el-tag v-else-if="scope.row.Status==2 && scope.row.IsSuccess == true" size="small" type="success">成功</el-tag>
                    <el-tag v-else-if="scope.row.Status==2 && scope.row.IsSuccess == false" size="small" type="danger">失败</el-tag>
                    <el-tag v-else-if="scope.row.Status==3" size="small" type="info">取消</el-tag>
                    <el-tag size="small" type="info">{{ scope.row.WorkflowsName }}</el-tag>
                    <el-tag v-if="scope.row.BranchName !=''" size="small">{{ scope.row.BranchName }}</el-tag>
                    <span style="margin-left: 5px ">{{ scope.row.AppName }}:{{ scope.row.BuildNumber }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="80">
                  <template #default="scope">
                    <el-button v-if="scope.row.Status!=0" size="small" type="success" @click="showBuildLog(scope.row)">日志</el-button>
                    <el-button v-else size="small" type="danger" @click="onStopBuild(scope.row.Id)">停止</el-button>
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
  <taskDialog ref="taskDialogRef"  />
  <el-dialog title="构建日志" v-model="state.isShowBuildLogDialog" style="width: 80%;top:20px;margin-bottom: 50px;" class='initdialog__body'>
    <el-checkbox v-model="state.autoLog" style="margin-bottom: 5px;">自动刷新日志</el-checkbox>
    <div class="layout-padding-auto" style="background-color:#393d49;">
      <div ref="scrollableBuildLog" style="height: 100%;overflow-y: auto;">
        <pre style="color: #fff;background-color:#393d49;padding: 5px 0 5px 5px;" v-html="state.buildLogContent"></pre>
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
// 定义变量内容
const appDialogRef = ref();
const appAddDialogRef = ref();
const scrollableBuildLog = ref();
const state = reactive({
  isShowBuildLogDialog: false,
  buildLogContent: '',
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
      pageSize: 22,
    },

  },
  appName:"",
  buildLogId:0,
  showOverlay:false,
  statTask:[],
  autoLog:true,
});
const buildItemCls =(row:any)=>{
  // return 'buildItem'
  const List = row.List;
  if(List){
    const len = List.length;
    if(len <=1){
      return 'buildItem'
    }else if(len == 2){
       return 'buildItem w5'
    }else if(len == 3){
       return 'buildItem w3'
    }else{
       return 'buildItem w10'
    }
  }
}
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
  let param = { };
  // 获取应用列表
  serverApi.autobuildList(param).then(function (res){
    if (res.Status) {
      // console.log(res.Data)
      const arr =  res.Data;
      arr.sort((a:any, b:any) => a.List.length - b.List.length);
      state.tableData.data = arr;
      state.tableData.total = res.Data.length;
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
    buildType:1,
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
const setAutoBuild = (item:any)=>{
  const param = {
    "appName": item.AppName,  // 应用名称
    "branchName": item.BranchName, // 分支名称
    "isAuto": item.AutoBuild     // 开关状态
  }
  serverApi.setAutoBuild(param).then(async function(res){
          if(res.Status){
            // ElMessage.success("自动构建成功")
            // 刷新
            getTableData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
}
// 构建
const goBuild = (item: any) => {
  ElMessageBox.confirm(`请确认重新自动构建?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 提交数据
        var param={ "commitId": item.CommitId }
        serverApi.autobuildResetCommitId(param).then(async function(res){
          if(res.Status){
            ElMessage.success("自动构建成功")
            // 刷新
            getTableData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
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
  }
});

// 显示构建日志
const showBuildLog=(row:any)=>{
  state.buildLogId = row.Id
  serverApi.buildLog(state.buildLogId.toString()).then(function (res){
    state.buildLogContent = res
    state.isShowBuildLogDialog= true;
    setTimeout(() => {   //自动跳到底部 
        scrollableBuildLog.value.scrollTop = scrollableBuildLog.value.scrollHeight;
        }, 500)
    if(row.Status == 2){
      state.autoLog = false
    }else{
      state.autoLog= true
      clearInterval(intervalId);
      intervalId = setInterval(onShowLog, 500);
    }
   
    
  })
}


const onShowLog=()=>{
  serverApi.buildLog(state.buildLogId.toString()).then(function (res) {
    // 如果从接口获取到的内容与本地内容一样时，则不用滚动
   if(state.buildLogContent != res){
    state.buildLogContent = res;
    // 自动刷新日志
    // console.log(state.autoLog)
    if (state.autoLog ) {
      setTimeout(() => {   //自动跳到底部 
        scrollableBuildLog.value.scrollTop = scrollableBuildLog.value.scrollHeight;
        }, 500)
       
    }
   }
    
  })
}

const onShowOverlay=()=>{
  state.showOverlay=true
}
const onHideOverlay=()=>{
  state.showOverlay=false
}




// 停止构建
const onStopBuild=(rowId: any)=>{
  ElMessageBox.confirm(`请确认是否停止构建?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 提交数据
        var param={ "buildId": rowId ,"buildType":1}
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

let intervalLogId = null;
let intervalAppId = null;
// 页面加载时
onMounted(() => {
  getTableData();
  getTableLogData();
  intervalAppId = setInterval(getTableData, 3000);
  intervalLogId = setInterval(getTableLogData, 3000);
});
// 页面注销的时候
onUnmounted(()=>{
  clearInterval(intervalLogId);
  clearInterval(intervalAppId);
})
</script>
<style lang="scss">
.buildMain{
  padding: 0;overflow: hidden;
  display: flex;
  width: 100%;
  flex-wrap: wrap; 
  justify-content: flex-start;
  align-items: flex-start;
  .w5{
    width: 50%;
  }
  .w3{
    width: 33%;
  }
  .w10{
    width: 100%;
  }
}
.buildItem{
padding-right: 5px;
// padding-bottom: 10px;
box-sizing: border-box;
display: flex;
align-items: center; 
flex: 1 1 auto; 
.el-tag--small {
    padding: 0 7px;
    height: 20px;
    --el-icon-size: 12px;
}
.el-button--small{
  padding: 0 7px;
    height: 20px;
    --el-icon-size: 12px;
}
.el-card {
  background-color:#f9f9e3;
  margin-bottom: 5px;
}
.el-card__body{
  padding: 0;
  
}
.title{
  background-color: #545c64;
  padding: 10px 15px;
  color: #fff;
  font-size: 14px;
  box-sizing: border-box;
}
._divs{
  box-sizing: border-box;
  display: flex;
  flex-flow: wrap;
  ._uls1{
    border: 1px solid #ccc;
  }
  ._uls{
    font-size: 12px;
    align-items: center; 
    flex: 1 1 auto; 
    margin: 5px;
    padding: 5px;
    border-radius: 5px;
    box-sizing: border-box;
    width: 25%;
    li{
      padding: 5px 3px;
      display: flex;
      flex-wrap: wrap;
      align-items: center; 
      span{
        display: flex;
        flex-wrap: wrap;
        align-items: center; 
      }
    }
  }
}
ul{
  margin: 0;
  padding: 0;
  li{
    list-style: none;
    box-sizing: border-box;
  }
}
}
</style>