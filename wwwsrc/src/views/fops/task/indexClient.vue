<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
				<el-table-column prop="Id" label="序号" width="200" />
				<el-table-column label="客户端">
          <template #default="scope">
          <div style="float: left;padding-right: 10px;padding-top: 5px">
            <el-tag v-if="scope.row.Status==0">刚上线</el-tag>
            <el-tag v-if="scope.row.Status==1">接受调度</el-tag>
            <el-tag v-if="scope.row.Status==2">无法调度</el-tag>
            <el-tag v-if="scope.row.Status==3">拒绝调度</el-tag>
            <el-tag v-if="scope.row.Status==4">离线</el-tag>
          </div>
          <div style="float: left">
              <span>{{scope.row.Name}}</span><br>
              <span>{{scope.row.Ip}}:{{scope.row.Port}}</span>
          </div>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="250" show-overflow-tooltip>
          <template #default="scope">
            <span>激活：{{scope.row.ActivateAt}}</span><br>
            <span>调度：{{scope.row.ScheduleAt}}</span>
          </template>
        </el-table-column>
        <el-table-column label="队列数量" show-overflow-tooltip>
          <template #default="scope">
            <span>队列数量：{{scope.row.QueueCount}}</span><br>
            <span>工作数量：{{scope.row.WorkCount}}</span><br>
            <span>错误次数：{{scope.row.ErrorCount}}</span>
          </template>
        </el-table-column>
        <el-table-column label="系统数据" show-overflow-tooltip>
          <template #default="scope">
            <span>CPU百分比：{{scope.row.CpuUsage}}</span><br>
            <span>内存百分比：{{scope.row.MemoryUsage}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="Jobs" label="Jobs" show-overflow-tooltip>
          <template #default="scope">
            <el-tag
                v-for="(tag, index) in scope.row.Jobs"
                :key="index">
              【{{ tag.Name }}-{{tag.Ver}}】
            </el-tag>
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
</style>
