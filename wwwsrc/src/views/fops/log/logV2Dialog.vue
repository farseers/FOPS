<template>
	<div class="system-user-container layout-padding">
    <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="80%">
      <div style="display: flex;flex-flow: column;max-height: calc(90vh - 151px) !important;">
        <el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
        <label>TraceId</label>
        <el-input size="default" v-model="state.traceId" placeholder="链路ID" style="max-width: 180px"> </el-input>
        <label class="ml10">应用名称</label>
        <el-select class="ml5" style="max-width: 110px;" size="small" v-model="state.appName">
          <el-option label="全部" value=""></el-option>
          <el-option v-for="item in state.appData" :label="item.AppName" :value="item.AppName" ></el-option>
        </el-select>
        <label class="ml10">执行端IP</label>
        <el-input class="ml5" size="default" v-model="state.appIp" placeholder="执行端IP" style="max-width: 120px;"></el-input>
        <label class="ml10">日志内容</label>
        <el-input size="default" v-model="state.logContent" placeholder="日志内容" clearable style="max-width: 250px;padding-left: 5px"> </el-input>
        <label class="ml10">日志类型</label>
        <el-select v-model="state.logLevel" placeholder="日志类型" class="ml10" style="max-width: 110px;" size="small">
          <el-option label="全部" :value="-1"></el-option>
          <el-option label="Trace" :value="0"></el-option>
          <el-option label="Debug" :value="1"></el-option>
          <el-option label="Info" :value="2"></el-option>
          <el-option label="Warning" :value="3"></el-option>
          <el-option label="Error" :value="4"></el-option>
          <el-option label="Critical" :value="5"></el-option>
        </el-select>
				<el-button size="default" type="primary" class="ml10" @click="onQuery">
					<el-icon>
						<ele-Search />
					</el-icon>
					查询
				</el-button>
			</div>
      <el-card style="color: #fff;background-color:#393d49;height: 100%;line-height:35px;overflow: auto;" class="layout-padding-auto">
        <p v-for="(v, k) in state.tableData.data" :key="k">
          <span style="color: #9caf62">{{v.CreateAt}}</span>
          {{v.TraceId}}
          <el-tag size="small" style="margin-right: 5px;">{{v.AppName}} {{v.AppIp}}</el-tag>
          <el-tag v-if="v.LogLevel == 'Info'" size="small">{{v.LogLevel}}</el-tag>
          <el-tag v-else-if="v.LogLevel == 'Debug'" type="info" size="small">{{v.LogLevel}}</el-tag>
          <el-tag v-else-if="v.LogLevel == 'Warn'" type="warning" size="small">{{v.LogLevel}}</el-tag>
          <el-tag v-else-if="v.LogLevel == 'Error'" type="danger" size="small">{{v.LogLevel}}</el-tag>
          {{v.Content}}
        </p>
      </el-card>
			<el-pagination
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				:page-sizes="[10, 20, 30]"
				v-model:current-page="state.tableData.param.pageNum"
				background
				v-model:page-size="state.tableData.param.pageSize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="state.tableData.total"
			>
			</el-pagination>
		</el-card>
    <detailDialog ref="detailDialogRef" @refresh="getTableData()" />
      </div>
    </el-dialog>
	</div>
</template>

<script setup lang="ts" name="fopsLogList">
import { defineAsyncComponent, reactive, onMounted, ref, watch  } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import {friendlyJSONstringify} from "@intlify/shared";

// 引入 api 请求接口
const serverApi = fopsApi();
// 引入组件
const detailDialog = defineAsyncComponent(() => import('/src/views/fops/log/detailDialog.vue'));


// 定义变量内容
const detailDialogRef = ref();

const state = reactive({
  traceId:'',
  appName:'',
  appIp:'',
  logContent:'',
  logLevel:-1,
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			pageNum: 1,
			pageSize: 18,
		},
	},
    dialog: {
    isShowDialog: false,
        type: '',
        title: '',
        submitTxt: '',
  },appData:[],
});
// 监听 state.startMin 的变化
watch(() => state.logLevel, (newValue, oldValue) => {
  console.log(`count 从 ${oldValue} 变为 ${newValue}`);
  getTableData()
});
watch(() => state.appName, (newValue, oldValue) => {
  console.log(`count 从 ${oldValue} 变为 ${newValue}`);
  getTableData()
});
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;

  var data={
    appName:state.appName,
    appIp:state.appIp,
    traceId:state.traceId,
    logContent:state.logContent,
    logLevel:state.logLevel.toString(),
    pageSize:state.tableData.param.pageSize.toString(),
    pageIndex:state.tableData.param.pageNum.toString(),
  }
  const params = new URLSearchParams(data).toString();
  // 请求接口
  serverApi.logList(params).then(function (res){
    if (res.Status){
      state.tableData.data = res.Data.List;
      state.tableData.total = res.Data.RecordCount;
      state.tableData.loading = false;
    }else{
      state.tableData.data=[]
      state.tableData.loading = false;
    }
  })

};
const onDetail=(row: any)=>{
  detailDialogRef.value.openDialog(row);
}
const openDialog = (row: any) => {
  state.dialog.isShowDialog = true;
  state.traceId=row.tid
  getTableData()
  getAppData();
}
const openDialogAppName = (row: any) => {
  state.dialog.isShowDialog = true;
  state.appName=row.AppName
  getTableData()
  getAppData();
}
const openDialogLogLevel = (level: any,appName:any,title:any) => {
  state.dialog.isShowDialog = true;
  state.dialog.title = title || ''
  state.logLevel=level
  state.appName=appName
  getTableData()
  getAppData();
}
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
// 删除用户
const onDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除：“${row.Name}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
      // 删除逻辑
      serverApi.taskDel({"TaskGroupId":row.Id}).then(function (res){
        if (res.Status){
          getTableData();
          ElMessage.success('删除成功');
        }else{
          ElMessage.error(res.StatusMessage)
        }
      })
		})
		.catch(() => {});
};
const getAppData=()=>{
  serverApi.dropDownList({}).then(function (res){
    if (res.Status){
      state.appData=res.Data
    }else{
      state.appData=[]
    }
  })
}
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
const onQuery=()=>{
  getTableData();
}

// 页面加载时
onMounted(() => {

});
// 暴露变量
defineExpose({
  openDialog,
  openDialogAppName,
  openDialogLogLevel
});
</script>

<style scoped lang="scss">
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
</style>
