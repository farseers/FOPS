<template>
	<div class="system-role-container layout-padding">
		<div class="system-role-padding layout-padding-auto layout-padding-view">
			<div class="system-user-search mb15">
				<!-- <el-input v-model="state.tableData.param.search" size="default" placeholder="请输入名称" style="max-width: 180px"> </el-input> -->
				<el-button size="default" type="primary" class="ml10" @click="getTableData()">
					<el-icon>
						<ele-Search />
					</el-icon>
					查询
				</el-button>
				<el-button size="default" type="success" class="ml10" @click="onOpenAddRole()">
					<el-icon>
						<ele-FolderAdd />
					</el-icon>
					新增
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
				<el-table-column type="index" label="序号" width="60" />
				<el-table-column prop="Id" label="ID" show-overflow-tooltip  width="60"></el-table-column>
				<el-table-column prop="Name" label="名称" show-overflow-tooltip></el-table-column>
				<el-table-column prop="LoginIp" label="Ip" show-overflow-tooltip></el-table-column>
				<el-table-column prop="LoginName" label="登录名" show-overflow-tooltip></el-table-column>
				<!-- <el-table-column prop="LoginPwd" label="登录密码" show-overflow-tooltip></el-table-column> -->
				<el-table-column prop="LoginPort" label="端口"  width="80" show-overflow-tooltip></el-table-column>
				<el-table-column label="操作" width="150">
					<template #default="scope">
						<el-button size="small" plain type="success" @click="onOpenEditRole(scope.row)">终端</el-button>
						<el-button size="small" plain type="danger" @click="onRowDel(scope.row)">删除</el-button>
					</template>
				</el-table-column>
                
			</el-table>
			<el-pagination
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				:page-sizes="[10, 20, 30,50]"
				v-model:current-page="state.tableData.param.pageIndex"
				background
				v-model:page-size="state.tableData.param.pageSize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="state.tableData.total"
			>
			</el-pagination>
		</div>
       
		<PaneDialog ref="roleDialogRef"
		:dialogVisible="state.dialogVisible"
		@refresh="getTableData()" />
		
	</div>
</template>

<script setup lang="ts" name="systemRole">
import { defineAsyncComponent, reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';

// 引入组件
const PaneDialog = defineAsyncComponent(() => import('./pane.vue'));
import {fopsApi} from "/@/api/fops";
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const roleDialogRef = ref();
const state = reactive({
    dialogVisible:false,
    fullscreen:true,
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			search: '',
			pageIndex: 1,
			pageSize: 10,
		},
	},
});
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	serverApi.terminalClientList({...state.tableData.param}).then(function (res){
		state.tableData.loading = false;
      if (res.Status){
		const Data = res.Data;
		const { List,RecordCount } = Data;
		state.tableData.total = RecordCount
        state.tableData.data = List
      }else{
		ElMessage.error(res.StatusMessage);
      }
    })
};
// 打开新增角色弹窗
const onOpenAddRole = () => {
	roleDialogRef.value.open()
};
// 打开修改角色弹窗
const onOpenEditRole = (row: any) => {
	serverApi.terminalClientInfo({Id:row.Id}).then((res)=>{
		if(res.Status){
			roleDialogRef.value.edit(res.Data)
		}else{
		ElMessage.error(res.StatusMessage);
      }
	})
};
// 删除角色
const onRowDel = (row: any) => {
	ElMessageBox.confirm(`确定删除${row.Name}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			serverApi.terminalClientDel({Id:row.Id}).then((res)=>{
				if(res.Status){
					getTableData();
					ElMessage.success('删除成功');
				}else{
				ElMessage.error(res.StatusMessage);
			}
			})
			
		})
		.catch(() => {});
};
// 分页改变
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pageSize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.pageIndex = val;
	getTableData();
};
// 页面加载时
onMounted(() => {
	getTableData();
});
</script>

<style lang="scss">
.system-role-container {
	.system-role-padding {
		padding: 15px;
		.el-table {
			flex: 1;
		}
	}
}

</style>
