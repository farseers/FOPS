<template>
	<div class="system-role-container layout-padding">
		<div class="system-role-padding layout-padding-auto layout-padding-view">
			<div class="system-user-search mb15">
				<!-- <el-input v-model="tableData.param.search" size="default" placeholder="请输入名称" style="max-width: 180px"> </el-input> -->
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
			<el-table :data="tableData.data" v-loading="tableData.loading" style="width: 100%">
				<el-table-column type="index" label="序号" width="60" />
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
				v-model:current-page="tableData.param.pageIndex"
				background
				v-model:page-size="tableData.param.pageSize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="tableData.total"
			>
			</el-pagination>
		</div>
       
		<PaneDialog ref="roleDialogRef"
		:dialogVisible="dialogVisible"
		@refresh="getTableData()" />
		
	</div>
</template>
<script>
import {fopsApi} from "/@/api/fops";
import { ElMessageBox, ElMessage } from 'element-plus';
import PaneDialog from './pane.vue';
const serverApi = fopsApi();
export default {
	name:'terminal',
	components:{PaneDialog},
	data(){
		return {
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
		}
	},
	mounted(){
		this.getTableData()
	},
	methods:{
		onOpenAddRole(){
			this.$refs.roleDialogRef && this.$refs.roleDialogRef.open()
		},
		getTableData(){
			this.tableData.loading = true;
			const _this = this;
			serverApi.terminalClientList({...this.tableData.param}).then(function (res){
				_this.tableData.loading = false;
			if (res.Status){
				const Data = res.Data;
				const { List,RecordCount } = Data;
				_this.tableData.total = RecordCount
				_this.tableData.data = List
			}else{
				ElMessage.error(res.StatusMessage);
			}
			})
		},
		onOpenEditRole(row){
			const _this = this;
			serverApi.terminalClientInfo({LoginIp:row.LoginIp}).then((res)=>{
				if(res.Status){
					_this.$refs.roleDialogRef && _this.$refs.roleDialogRef.edit(res.Data)
				}else{
				ElMessage.error(res.StatusMessage);
			}
			})
		},
	   onHandleSizeChange(val) {
		this.tableData.param.pageSize = val;
		this.getTableData();
		},
		// 分页改变
		onHandleCurrentChange (val){
			this.tableData.param.pageIndex = val;
			this.getTableData();
		},
		onRowDel(row){
			const _this = this;
			ElMessageBox.confirm(`确定删除${row.Name}”，是否继续?`, '提示', {
				confirmButtonText: '确认',
				cancelButtonText: '取消',
				type: 'warning',
			})
				.then(() => {
					serverApi.terminalClientDel({LoginIp:row.LoginIp}).then((res)=>{
						if(res.Status){
							_this.getTableData();
							ElMessage.success('删除成功');
						}else{
						ElMessage.error(res.StatusMessage);
					}
					})
					
				})
				.catch(() => {});
				}
	}
}
	
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
