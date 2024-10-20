<template>
    <div class="system-user-container layout-padding">
        <el-card shadow="hover" class="configu_m">
			<div class="system-user-search">
                <!-- <el-input placeholder="请输入内容" v-model="state.appName"  size="small" clearable style="width: 150px;"></el-input> -->
                <el-button type="primary" size="small" @click="getTableData()">刷新</el-button>
                <el-button type="success" plain size="small" @click="onOpenAdd()">添加</el-button>
			</div>
            <div class="flex1" ref="tableDiv">
                <el-table 
                :data="state.tableData" v-loading="state.loading"
                size="small"
                :maxHeight="state.height"
             style="width: 100%">
                <el-table-column type="index" label="序号" width="50" align="center"></el-table-column>
				<el-table-column prop="AppName" label="应用名称" min-width="100" ></el-table-column>
                <el-table-column prop="Key" label="键" min-width="150" ></el-table-column>
                <el-table-column prop="Value" label="值" min-width="220" ></el-table-column>
				<el-table-column prop="Ver" label="版本" min-width="100" ></el-table-column>
                <el-table-column label="操作" align="center"  min-width="150">
                    <template #default="scope">
                        <el-button @click="handleDel(scope.row)" type="text" size="small">删除</el-button>
                        <el-button @click="onOpenEdit(scope.row)" type="text" size="small">编辑</el-button>
                        <el-button @click="handleRollback(scope.row)" type="text" size="small">回滚</el-button>
                    </template>    
                </el-table-column>
			</el-table>
            </div>
			
		</el-card>
        <Dialog ref="DialogRef" @refresh="getTableData()" />
    </div>
</template>
<script setup lang="ts" name="configu">
import { defineAsyncComponent, reactive, onMounted, ref,watch} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
const Dialog = defineAsyncComponent(() => import('/@/views/fops/configu/dialog.vue'));
import {fopsApi} from "/@/api/fops";
const serverApi = fopsApi();
const DialogRef = ref();
const tableDiv = ref();
const state = reactive({
	tableData:[],
    loading:false,
    height:'500'
});
// 初始化表格数据
const getTableData = () => {
	state.loading = true;
	const data = [];
  // 请求接口
  serverApi.configureAllList().then(function (res){
    if (res.Status){
       state.tableData = res.Data;
        state.loading = false;
    }else{
      state.tableData=[]
      state.loading = false;
    }
  })

};//
const handleRollback = (row:any) =>{ //回滚
    ElMessageBox.confirm(`请确认是否回滚此项?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 提交数据
        const param={
          "appName" : row.AppName,
          "key":row.Key
        }
        serverApi.configureRollback(param).then(async function(res){
          if(res.Status){
            ElMessage.success("回滚成功")
            getTableData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
}
const handleDel = (row:any) =>{ //删除
    ElMessageBox.confirm(`请确认是否删除此项?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 提交数据
        const param={
          "appName" : row.AppName,
          "key":row.Key
        }
        serverApi.configureDelete(param).then(async function(res){
          if(res.Status){
            ElMessage.success("删除成功")
            getTableData()
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
}
// 打开新增用户弹窗
const onOpenAdd = () => {
    DialogRef.value.openDialog();
};
// 打开修改用户弹窗
const onOpenEdit = (row: any) => {
    DialogRef.value.openDialog(row);
};
// 页面加载时
onMounted(() => {
	getTableData();
   state.height = tableDiv.value.offsetHeight
});
</script>
<style>
.el-table .el-table__body tr:nth-child(odd) {
  background-color: #f2f2f2;
}
.configu_m{
    height: 100%;
   
}
.configu_m .el-card__body{
    height: 100%;
    display: flex;
    flex-flow: column;
}
.flex1{
    flex: 1;
}
</style>