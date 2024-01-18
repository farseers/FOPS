<template>
		<el-dialog  v-model="state.dialog.isShowDialog" style="height: 80%;width: 70%">
      <div class="system-user-container layout-padding" style="width: 100%;">
        <el-card shadow="hover" class="layout-padding-auto">
          <div class="system-user-search mb15">
            <el-input size="default" v-model="state.clientName" placeholder="请输入应用名称" style="max-width: 180px"> </el-input>
            <el-input size="default" v-model="state.taskGroupName" placeholder="请输入任务组名称" style="max-width: 180px"> </el-input>
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
          <el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
            <el-table-column prop="TaskId" label="任务ID" width="180" />
            <el-table-column label="日志内容">
              <template #default="scope">
                <el-tag v-if="scope.row.LogLevel == 'Info'" size="small">{{scope.row.LogLevel}}</el-tag>
                <el-tag v-else-if="scope.row.LogLevel == 'Debug'" type="info" size="small">{{scope.row.LogLevel}}</el-tag>
                <el-tag v-else-if="scope.row.LogLevel == 'Warn'" type="warning" size="small">{{scope.row.LogLevel}}</el-tag>
                <el-tag v-else-if="scope.row.LogLevel == 'Error'" type="danger" size="small">{{scope.row.LogLevel}}</el-tag>
                <span v-else>{{scope.row.LogLevel}}</span>
                {{scope.row.Content}}
              </template>
            </el-table-column>
            <el-table-column prop="CreateAt" width="170" label="日志时间" show-overflow-tooltip></el-table-column>
            <!--				<el-table-column label="操作" width="100">-->
            <!--					<template #default="scope">-->
            <!--						<el-button size="small" text type="primary" @click="onDetail(scope.row)">详情信息</el-button>-->
            <!--            <el-button size="small" text type="primary" @click="onEdit('edit',scope.row)">修改</el-button>-->
            <!--            <el-button size="small" text type="primary" @click="onDel(scope.row)">删除</el-button>-->
            <!--					</template>-->
            <!--				</el-table-column>-->
          </el-table>
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
		</el-dialog>
</template>

<script setup lang="ts" name="fopsTask">
import { defineAsyncComponent, reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";

// 引入 api 请求接口
const serverApi = fopsApi();
// 引入组件


// 定义变量内容
const editDialogRef = ref();
const state = reactive({
  keyWord:'',
  taskGroupName:'',
  clientName:'',
  logLevel:-1,
  taskId:'',
  tableData: {
    data: [],
    total: 0,
    loading: false,
    param: {
      pageNum: 1,
      pageSize: 10,
    },
  },dialog: {
    isShowDialog: false,
    type: '',
    title: '',
    submitTxt: '',
  },
});

// 初始化表格数据
const getTableData = () => {
  state.tableData.loading = true;

  const params = new URLSearchParams();
  params.append('logLevel', state.logLevel.toString());
  params.append('taskGroupName', state.taskGroupName.toString());
  params.append('taskId', state.taskId.toString());
  params.append('clientName', state.clientName.toString());
  params.append('pageSize', state.tableData.param.pageSize.toString());
  params.append('pageIndex', state.tableData.param.pageNum.toString());

  // 请求接口
  serverApi.taskLogListClientName(params.toString()).then(function (res){
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
const onDetail=(row: any)=>{

}
const onEdit=(type: string,row:any)=>{
  editDialogRef.value.openDialog(type, row);
}
const openDialog = (row: any) => {
  state.clientName=row.AppName
  state.dialog.isShowDialog = true;
  getTableData();
};
 
// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
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
// 暴露变量
defineExpose({
  openDialog,
  closeDialog,
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
