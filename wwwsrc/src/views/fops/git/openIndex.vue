<template>
	<div class="system-user-container layout-padding">
    <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-button size="default" type="success" class="ml10" @click="SureCheck()">
					<el-icon>
            <ele-FolderAdd />
					</el-icon>
					确认选择
				</el-button>
			</div>
			<el-table ref="multipleTable" :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" :row-key="getRowKey" @selection-change="handleSelectionChange">
        <el-table-column type="selection" :reserve-selection="true" width="55"></el-table-column>
				<el-table-column prop="Id"  label="编号" width="60" />
				<el-table-column prop="Name" label="Git名称" show-overflow-tooltip></el-table-column>
				<el-table-column prop="Hub" label="托管地址" show-overflow-tooltip></el-table-column>
				<el-table-column prop="Branch" label="Git分支" show-overflow-tooltip></el-table-column>
			</el-table>
		</el-card>
      </el-dialog>
	</div>
</template>

<script setup lang="ts" name="fopsGit">
import { defineAsyncComponent, reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";

// 定义子组件向父组件传值/事件
const emit = defineEmits<{
  (event: 'selectItem', items: []): void
}>()

// 引入 api 请求接口
const serverApi = fopsApi();

// 定义变量内容
const gitDialogRef = ref();
const state = reactive({
  SelectItem:[],
  dialog: {
    isShowDialog: false,
    type: '',
    title: '',
    submitTxt: '',
  },
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
// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
// 取消
const onCancel = () => {
  closeDialog();
};
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	const data = [];
  // 请求接口
  serverApi.gitList({}).then(function (res){
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
const getRowKey=(row:any)=>{
  return row.Id;
}
const handleSelectionChange=(val:any)=> {
  console.log(val)
  if(val.length==0){return;}
  state.SelectItem.push(val.Id)
  console.log(state.SelectItem)
}

// 确认选择
const SureCheck=()=>{

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
.el-checkbox__inner{
  border: 1px solid #666;
}
</style>
