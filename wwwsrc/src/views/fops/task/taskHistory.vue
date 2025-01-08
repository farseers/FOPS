<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
        <el-input size="default" v-model="state.clientName" placeholder="请输入应用名称" clearable style="max-width: 180px"> </el-input>
				<el-input size="default" v-model="state.taskGroupName" placeholder="请输入任务组名称" clearable style="max-width: 180px"> </el-input>
        <el-input size="default" v-model="state.taskId" placeholder="请输入任务ID" clearable style="max-width: 180px"  class="ml10"> </el-input>
        <el-select v-model="state.scheduleStatus" placeholder="调度结果" class="ml10">
          <el-option label="全部" :value="-1"></el-option>
          <el-option label="未调度" :value="0"></el-option>
          <el-option label="调度中" :value="1"></el-option>
          <el-option label="调度成功" :value="2"></el-option>
          <el-option label="调度失败" :value="3"></el-option>
        </el-select>
        <el-select v-model="state.executeStatus" placeholder="执行结果" class="ml10" style="margin-left: 5px;">
          <el-option label="全部" :value="-1"></el-option>
          <el-option label="未开始" :value="0"></el-option>
          <el-option label="执行中" :value="1"></el-option>
          <el-option label="成功" :value="2"></el-option>
          <el-option label="失败" :value="3"></el-option>
        </el-select>
				<el-button size="default" type="primary" class="ml10" @click="onQuery">
					<el-icon>
						<ele-Search />
					</el-icon>
					查询
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" class="mytable">
        <el-table-column prop="Id" label="任务ID" width="180">
          <template #default="scope">
            <div style="float:left;;">
              <span title="任务ID">{{scope.row.Id}}</span><br>
              <span title="TraceId">{{scope.row.TraceId}}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="任务" width="400">
          <template #default="scope">
            <div style="float:left;margin: 6px">
              <el-tag size="small" v-if="scope.row.ExecuteStatus==0" type="info">未执行</el-tag>
              <el-tag size="small" v-else-if="scope.row.ExecuteStatus==1" type="success" style="color:green">执行中</el-tag>
              <el-tag size="small" v-else-if="scope.row.ExecuteStatus==2" type="success" style="color:green">成功</el-tag>
              <el-tag size="small" v-else-if="scope.row.ExecuteStatus==3" type="danger">失败</el-tag>
            </div>
            <span>{{scope.row.Caption}}</span><br>
            <span>{{scope.row.Name}}（<span style="color:#4eb8ff">Ver:{{scope.row.Ver}}</span>）</span>
          </template>
        </el-table-column>
        <el-table-column prop="StartAt" label="时间" width="210" show-overflow-tooltip>
          <template #default="scope">
            <span>开始: {{scope.row.StartAt}}</span><br>
            <span>完成: {{scope.row.FinishAt}}</span>
          </template>
        </el-table-column>
        <el-table-column label="运行情况"  width="130" show-overflow-tooltip>
          <template #default="scope">
            <span>耗时: <el-tag size="small" type="info">{{scope.row.RunSpeed}}</el-tag></span><br>
            <span>进度: <el-tag size="small" type="info">{{scope.row.Progress}}%</el-tag></span>
          </template>
        </el-table-column>
        <el-table-column label="数据" show-overflow-tooltip>
          <template #default="scope">
            <span>{{friendlyJSONstringify(scope.row.Data)}}</span><br />
            <el-tag size="small" v-if="scope.row.Remark!=''" type="danger">{{scope.row.Remark}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="客户端信息" width="180" show-overflow-tooltip>
          <template #default="scope">
              <el-tag size="small" v-if="scope.row.ScheduleStatus==0" type="info">未调度</el-tag>
              <el-tag size="small" v-else-if="scope.row.ScheduleStatus==1" type="success" style="color:green">调度中</el-tag>
              <el-tag size="small" v-else-if="scope.row.ScheduleStatus==2" type="success" style="color:green">调度成功</el-tag>
              <el-tag size="small" v-else-if="scope.row.ScheduleStatus==3" type="danger">调度失败</el-tag>
              <br />
              <el-tag size="small">{{scope.row.Client.Name}} {{scope.row.Client.Ip}}:{{scope.row.Client.Port}}</el-tag>
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

// 定义变量内容
const state = reactive({
  appName:'',
  enable:-1,
  scheduleStatus:-1,
  executeStatus:-1,
  taskGroupName:'',
  clientName:'',
  clientId:'',
  taskId:'',
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			pageNum: 1,
			pageSize: 19,
		},
	},
  NowTime:new Date(),
  appData:[],
});

watch(() => state.appName, (newValue, oldValue) => {
  getTableData()
});

const getAppData=()=>{
  serverApi.dropDownList({}).then(function (res){
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
  params.append('clientName', state.clientName);
  params.append('taskGroupName', state.taskGroupName);
  params.append('scheduleStatus', state.scheduleStatus.toString());
  params.append('executeStatus', state.executeStatus.toString());
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

const onQuery=()=>{
  getTableData();
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

watch(() => state.scheduleStatus, (newValue, oldValue) => {
  getTableData()
});

watch(() => state.executeStatus, (newValue, oldValue) => {
  getTableData()
});

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
