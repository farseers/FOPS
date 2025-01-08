<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-button size="default" type="success" class="ml10" @click="onOpenAdd('add')">
					<el-icon>
						<ele-FolderAdd />
					</el-icon>
					新增集群
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
				<el-table-column prop="Id" label="序号" width="60" />
				<el-table-column prop="Name" label="集群名称" show-overflow-tooltip></el-table-column>
				<el-table-column prop="IsLocal" label="本地" width="100" show-overflow-tooltip></el-table-column>
				<el-table-column prop="FopsAddr" label="集群地址" show-overflow-tooltip></el-table-column>
				<el-table-column prop="FScheduleAddr" label="调度中心" show-overflow-tooltip></el-table-column>
				<el-table-column prop="DockerNetwork" label="Docker网络" show-overflow-tooltip></el-table-column>
				<el-table-column prop="DockerHub" label="DockerHub地址" show-overflow-tooltip></el-table-column>
        <el-table-column prop="DockerUserName" label="账户名称" show-overflow-tooltip></el-table-column>
				<el-table-column label="操作" width="100">
					<template #default="scope">
						<el-button size="small" text type="primary" @click="onOpenEdit('edit', scope.row)">修改</el-button>
						<el-button size="small" text type="primary" @click="onRowDel(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
		<clusterDialog ref="clusterDialogRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts" name="fopsCluster">
import { defineAsyncComponent, reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";

// 引入 api 请求接口
const serverApi = fopsApi();

// 引入组件
const clusterDialog = defineAsyncComponent(() => import('/@/views/fops/cluster/dialog.vue'));

// 定义变量内容
const clusterDialogRef = ref();
const state = reactive({
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
  // 请求接口
  serverApi.clusterList({}).then(function (res){
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

// 打开新增用户弹窗
const onOpenAdd = (type: string) => {
  clusterDialogRef.value.openDialog(type);
};
// 打开修改用户弹窗
const onOpenEdit = (type: string, row: any) => {
  clusterDialogRef.value.openDialog(type, row);
};
// 删除用户
const onRowDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除：“${row.Name}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
      // 删除逻辑
      serverApi.clusterDel({"ClusterId":row.Id}).then(function (res){
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
