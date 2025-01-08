<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-input size="default" v-model="state.taskGroupName" placeholder="请输入任务组名称" clearable style="max-width: 180px"> </el-input>
        <el-select v-model="state.logLevel" placeholder="请选择日志等级" clearable class="ml10">
          <el-option label="全部" :value="-1"></el-option>
          <el-option label="Trace" :value="0"></el-option>
          <el-option label="Debug" :value="1"></el-option>
          <el-option label="Info" :value="2"></el-option>
          <el-option label="Warning" :value="3"></el-option>
          <el-option label="Error" :value="4"></el-option>
          <el-option label="Critical" :value="5"></el-option>
        </el-select>
        <el-input size="default" v-model="state.taskId" placeholder="请输入任务ID" clearable style="max-width: 180px"  class="ml10"> </el-input>
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
          {{v.TaskId}}
          <el-tooltip :content="v.Caption" slot="label">
            <el-tag size="small" style="margin-right: 5px;">{{v.Name}}</el-tag>
          </el-tooltip>
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
	</div>
</template>

<script setup lang="ts" name="fopsTaskRunning">
import { reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";

// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
  keyWord:'',
  taskGroupName:'',
  logLevel:-1,
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
});

// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;

  const params = new URLSearchParams();
  params.append('taskId', state.taskId);
  params.append('logLevel', state.logLevel.toString());
  params.append('taskGroupName', state.taskGroupName.toString());
  params.append('pageSize', state.tableData.param.pageSize.toString());
  params.append('pageIndex', state.tableData.param.pageNum.toString());

  // 请求接口
  serverApi.taskLogList(params.toString()).then(function (res){
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
	getTableData();
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
