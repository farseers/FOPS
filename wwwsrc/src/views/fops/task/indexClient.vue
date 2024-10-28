<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
				<el-table-column label="客户端">
          <template #default="scope">
          <div style="float: left;margin-right: 10px;">
            <el-tag v-if="scope.row.Status==0">刚上线</el-tag>
            <el-tag v-if="scope.row.Status==1">接受调度</el-tag>
            <el-tag v-if="scope.row.Status==2">无法调度</el-tag>
            <el-tag v-if="scope.row.Status==3">拒绝调度</el-tag>
            <el-tag v-if="scope.row.Status==4">离线</el-tag>
          </div>
          <div style="float: left">
            <el-tag type="danger" size="small" v-if="scope.row.IsMaster" style="margin-right: 5px;">主</el-tag> {{scope.row.Name}} {{scope.row.Ip}}:{{scope.row.Port}} | {{scope.row.Job.Name}}（<span style="color:#4eb8ff">Ver:{{scope.row.Job.Ver}}</span>）
          </div>
          </template>
        </el-table-column>
        <el-table-column label="激活时间" width="250" show-overflow-tooltip>
          <template #default="scope">
            <span>{{scope.row.ActivateAt}}</span>
          </template>
        </el-table-column>
        <el-table-column label="调度时间" width="250" show-overflow-tooltip>
          <template #default="scope">
            <span v-if='scope.row.ScheduleAt != "0001-01-01 00:00:00" '>{{scope.row.ScheduleAt}}</span>
          </template>
        </el-table-column>
        <el-table-column label="队列数量" width="90" show-overflow-tooltip>
          <template #default="scope">
            <span>{{scope.row.QueueCount}}</span>
          </template>
        </el-table-column>
        <el-table-column label="工作数量" width="90" show-overflow-tooltip>
          <template #default="scope">
            <span>{{scope.row.WorkCount}}</span>
          </template>
        </el-table-column>
        <el-table-column label="错误次数" width="90" show-overflow-tooltip>
          <template #default="scope">
            <span>{{scope.row.ErrorCount}}</span>
          </template>
        </el-table-column>
			</el-table>
		</el-card>
	</div>
</template>

<script setup lang="ts" name="fopsTaskTimeOut">
import { defineAsyncComponent, reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";

// 引入 api 请求接口
const serverApi = fopsApi();


// 定义变量内容
const state = reactive({
  keyWord:'',
  enable:-1,
  taskStatus:-1,
  clientId:'',
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			pageNum: 1,
			pageSize: 10,
		},
	},
});

// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	var param={
  }
  // 请求接口
  serverApi.clientList(param).then(function (res){
    if (res.Status){
      state.tableData.data = res.Data;
      state.tableData.total = res.Data.length;
        state.tableData.loading = false;
    }else{
      state.tableData.data=[]
        state.tableData.loading = false;
    }

  })
};
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
.el-table__row .el-table__cell{
  padding: 0 0;
}
</style>
