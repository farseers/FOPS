<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
        <el-select v-model="state.isApp" placeholder="Git类型" class="ml10" style="max-width: 150px;" size="small">
          <el-option label="全部" :value="-1"></el-option>
          <el-option label="框架" :value="0"></el-option>
          <el-option label="应用" :value="1"></el-option>
        </el-select>
<!--				<el-input size="default" placeholder="请输入用户名称" style="max-width: 180px"> </el-input>-->
<!--				<el-button size="default" type="primary" class="ml10">-->
<!--					<el-icon>-->
<!--						<ele-Search />-->
<!--					</el-icon>-->
<!--					查询-->
<!--				</el-button>-->
				<el-button size="default" type="success" class="ml10" @click="onOpenAdd('add')">
					<el-icon>
            <ele-FolderAdd />
					</el-icon>
					新增Git
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
				<el-table-column prop="Id" label="编号" width="60" />
				<el-table-column label="Git名称" width="200" show-overflow-tooltip>
          <template #default="scope">
            <el-tag v-if="scope.row.IsApp == true" size="small">{{scope.row.Name}}</el-tag>
            <el-tag v-else size="small" type="info">{{scope.row.Name}}</el-tag>
          </template>
        </el-table-column>
				<el-table-column prop="Hub" label="托管地址" show-overflow-tooltip></el-table-column>
				<el-table-column prop="Branch" label="Git分支" width="100" show-overflow-tooltip></el-table-column>
				<el-table-column prop="UserName" label="账户名称" width="180" show-overflow-tooltip></el-table-column>
        <el-table-column prop="Path" label="存储目录" width="220" show-overflow-tooltip></el-table-column>
        <el-table-column prop="PullAt" label="拉取时间" width="180" show-overflow-tooltip></el-table-column>
				<el-table-column label="操作" width="170">
					<template #default="scope">
						<el-button size="small" text type="primary" @click="onOpenEdit('edit', scope.row)">修改</el-button>
						<el-button size="small" text type="primary" @click="onRowDel(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
<!--			<el-pagination-->
<!--				@size-change="onHandleSizeChange"-->
<!--				@current-change="onHandleCurrentChange"-->
<!--				class="mt15"-->
<!--				:pager-count="5"-->
<!--				:page-sizes="[10, 20, 30]"-->
<!--				v-model:current-page="state.tableData.param.pageNum"-->
<!--				background-->
<!--				v-model:page-size="state.tableData.param.pageSize"-->
<!--				layout="total, sizes, prev, pager, next, jumper"-->
<!--				:total="state.tableData.total"-->
<!--			>-->
<!--			</el-pagination>-->
		</el-card>
		<GitDialog ref="gitDialogRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts" name="fopsGit">
import { defineAsyncComponent, reactive, onMounted, ref,watch} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";

// 引入 api 请求接口
const serverApi = fopsApi();

// 引入组件
const GitDialog = defineAsyncComponent(() => import('/@/views/fops/git/dialog.vue'));

// 定义变量内容
const gitDialogRef = ref();
const state = reactive({
  isApp:-1,
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
// 监听 state.isApp 的变化
watch(() => state.isApp, (newValue, oldValue) => {
  console.log(`count 从 ${oldValue} 变为 ${newValue}`);
  getTableData()
});
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	const data = [];
  // 请求接口
  serverApi.gitList({isApp:state.isApp}).then(function (res){
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
  gitDialogRef.value.openDialog(type);
};
// 打开修改用户弹窗
const onOpenEdit = (type: string, row: any) => {
  gitDialogRef.value.openDialog(type, row);
};

// 删除
const onRowDel = (row: any) => {
  ElMessageBox.confirm(`此操作将永久删除：“${row.Name}”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 删除逻辑
        serverApi.gitDel({"GitId":row.Id}).then(function (res){
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
