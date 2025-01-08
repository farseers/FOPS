<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
				<el-table-column prop="Id" label="序号" width="200" />
				<el-table-column prop="Name" label="服务器节点名称" show-overflow-tooltip></el-table-column>
				<el-table-column prop="Ip" label="IP" show-overflow-tooltip></el-table-column>
				<el-table-column prop="Port" label="端口" show-overflow-tooltip></el-table-column>
				<el-table-column prop="IsLeader" label="是否主节点" show-overflow-tooltip>
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.IsLeader == true" type="danger">主节点</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ActivateAt" label="激活时间" show-overflow-tooltip></el-table-column>
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
  serverApi.serverNodeList(param).then(function (res){
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
