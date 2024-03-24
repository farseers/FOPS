<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">

				<el-input size="default" v-model="state.keyWord" placeholder="请输入任务组名称" clearable style="max-width: 180px"> </el-input>
        <el-input size="default" v-model="state.taskId" placeholder="请输入任务ID" clearable style="max-width: 180px"  class="ml10"> </el-input>
        <el-select v-model="state.taskStatus" placeholder="请选择调度状态" class="ml10" @change="onStatusChange">
          <el-option label="全部" :value="-1"></el-option>
          <el-option label="成功" :value="5"></el-option>
          <el-option label="失败" :value="4"></el-option>
        </el-select>
				<el-button size="default" type="primary" class="ml10" @click="onQuery">
					<el-icon>
						<ele-Search />
					</el-icon>
					查询
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" class="mytable">
        <el-table-column prop="Id" label="任务ID" width="250">
          <template #default="scope">
            <div style="float:left;margin: 6px">
              <el-tag v-if="scope.row.Status==0" style="color:#7a7a7a">未开始</el-tag>
              <el-tag v-else-if="scope.row.Status==1">调度中</el-tag>
              <el-tag v-else-if="scope.row.Status==2" style="color:red">调度失败</el-tag>
              <el-tag v-else-if="scope.row.Status==3">执行中</el-tag>
              <el-tag v-else-if="scope.row.Status==4" style="color: red">失败</el-tag>
              <el-tag v-else-if="scope.row.Status==5" style="color:green">成功</el-tag>
            </div>
            <div style="float:left;;">
              <span title="任务ID">{{scope.row.Id}}</span><br>
              <span title="TraceId">{{scope.row.TraceId}}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="任务名称" >
          <template #default="scope">
            <span>{{scope.row.Caption}}</span><br>
            <span>{{scope.row.Name}}（<span style="color:#4eb8ff">Ver:{{scope.row.Ver}}</span>）</span>
          </template>
        </el-table-column>
        <el-table-column prop="StartAt" label="时间" width="210" show-overflow-tooltip>
          <template #default="scope">
            <span>开始: {{scope.row.StartAt}}</span><br>
            <span>完成: {{scope.row.RunAt}}</span>
          </template>
        </el-table-column>
        <el-table-column label="运行情况"  width="110" show-overflow-tooltip>
          <template #default="scope">
            <span>耗时: {{scope.row.RunSpeed}}</span><br>
            <span>进度: {{scope.row.Progress}}%</span>
          </template>
        </el-table-column>
        <el-table-column label="数据"  width="450">
          <template #default="scope">
            <span>{{friendlyJSONstringify(scope.row.Data)}}</span>
          </template>
        </el-table-column>
        <el-table-column label="客户端信息" width="180" show-overflow-tooltip>
          <template #default="scope">
            <div>
              <el-tag size="small">{{scope.row.Client.Name}} {{scope.row.Client.Ip}}:{{scope.row.Client.Port}}</el-tag>
            </div>
          </template>
        </el-table-column>

			</el-table>
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
    <editDialog ref="editDialogRef" @refresh="getTableData()" />
    <detailDialog ref="detailDialogRef" @refresh="getTableData()" />
    <taskDialog ref="taskDialogRef" @refresh="getTableData()" />
    <logDialog ref="logDialogRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts" name="fopsTask">
import {defineAsyncComponent, reactive, onMounted, ref, nextTick, watch} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import {friendlyJSONstringify} from "@intlify/shared";
import {time} from "echarts/core";

// 引入 api 请求接口
const serverApi = fopsApi();

// 引入组件
const editDialog = defineAsyncComponent(() => import('/src/views/fops/task/editGroupDialog.vue'));
const detailDialog = defineAsyncComponent(() => import('/src/views/fops/task/detailGroupDialog.vue'));
const taskDialog = defineAsyncComponent(() => import('/src/views/fops/task/taskDialog.vue'));
const logDialog = defineAsyncComponent(() => import('/src/views/fops/task/logDialog.vue'));


// 定义变量内容
const editDialogRef = ref();
const detailDialogRef = ref();
const taskDialogRef = ref();
const logDialogRef = ref();
const state = reactive({
  keyWord:'',
  appName:'',
  enable:-1,
  taskStatus:-1,
  clientId:'',
  taskId:'',
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			pageNum: 1,
			pageSize: 10,
		},
	},
  NowTime:new Date(),
  appData:[],
});

// 监听 state.taskStatus 的变化
watch(() => state.taskStatus, (newValue, oldValue) => {
  getTableData()
});

// 监听 state.enable 的变化
watch(() => state.enable, (newValue, oldValue) => {
  getTableData()
});

watch(() => state.appName, (newValue, oldValue) => {
  getTableData()
});

const getAppData=()=>{
  serverApi.appsList({}).then(function (res){
    if (res.Status){
      state.appData=res.Data
    }else{
      state.appData=[]
    }
  })
}

// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
  const params = new URLSearchParams();
  //params.append('clientName', state.appName);
  params.append('taskGroupName', state.keyWord);
  //params.append('enable', state.enable.toString());
  params.append('taskStatus', state.taskStatus.toString());
  //params.append('clientId', state.clientId);
  params.append('taskId', state.taskId);
  params.append('pageSize', state.tableData.param.pageSize.toString());
  params.append('pageIndex', state.tableData.param.pageNum.toString());

  // 请求接口
  serverApi.taskList(params.toString()).then(function (res){
    if (res.Status){
      state.tableData.data = res.Data.List;
      state.tableData.total = res.Data.RecordCount;
    }else{
      state.tableData.data=[]
    }
    state.tableData.loading = false;
  })
};

const compareTime=(nextAt:any)=>{
  var convertedTime = new Date(nextAt)
  return convertedTime.getTime() < new Date().getTime();
}
const onDetail=(row: any)=>{
  detailDialogRef.value.openDialog(row);
}
const onQuery=()=>{
  getTableData();
}
const onEdit=(type: string, row: any)=>{
  editDialogRef.value.openDialog(type, row);
}
const onTaskList=(row: any)=>{
  taskDialogRef.value.openDialog(row);
}
const onLog=(row: any)=>{
  logDialogRef.value.openDialog(row);
}
// 删除
const onDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除：“${row.Name}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
      // 删除逻辑
      serverApi.taskDel({"taskGroupName":row.Name}).then(function (res){
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
//启用停用
const onIsEnable=(row: any)=>{
  var setEnable=row.IsEnable
  var tips=""
  if(setEnable){
    setEnable=false
    tips="停用"
  }else{
    setEnable=true
    tips="启用"
  }

  ElMessageBox.confirm(`该任务即将：“${tips}”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 设置状态
        serverApi.taskGroupSetEnable({"taskGroupName":row.Name,"enable":setEnable}).then(function (res){
          if (res.Status){
            getTableData();
            if(setEnable){
              ElMessage.success('启用-成功');
            }else{
              ElMessage.success('停用-成功');
            }

          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
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
const onStatusChange=(value:number)=>{
  state.taskStatus=value
}
const onEnableChange=(value:number)=>{
  state.enable=value
}
// 页面加载时
onMounted(() => {
  // 等待下一次 DOM 更新后再执行代码
  nextTick(() => {
    getTableData();
    getAppData();
  });
});
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

.el-table tr td {
  /* 你的自定义样式 */
  padding: 0 0!important;
}
</style>
