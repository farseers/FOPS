<template>
		<el-dialog  v-model="state.dialog.isShowDialog" style="height: 90%;width: 1300px">
      <div class="system-user-container layout-padding" style="width: 100%;">
        <el-card shadow="hover" class="layout-padding-auto">
          <div class="system-user-search mb15">
            <span>任务组名称：{{state.dialog.title}}</span>
            <el-select v-model="state.taskStatus" placeholder="请选择调度状态" class="ml10" @change="onStatusChange">
              <el-option label="全部" :value="-1"></el-option>
              <el-option style="color:#7a7a7a" label="未开始" :value="0"></el-option>
              <el-option label="调度中" :value="1"></el-option>
              <el-option label="调度失败" :value="2"></el-option>
              <el-option label="执行中" :value="3"></el-option>
              <el-option label="失败" :value="4"></el-option>
              <el-option label="成功" :value="5"></el-option>
            </el-select>
            <el-button size="default" type="primary" class="ml10" @click="onQuery">
              <el-icon>
                <ele-Search />
              </el-icon>
              查询
            </el-button>
          </div>
          <el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%;">
            <el-table-column prop="Id" label="任务ID" width="180">
              <template #default="scope">
                <span title="任务ID">{{scope.row.Id}}</span><br>
                <span title="TraceId">{{scope.row.TraceId}}</span>
              </template>
            </el-table-column>
            <el-table-column prop="StartAt" label="时间" width="210" show-overflow-tooltip>
              <template #default="scope">
                <span>开始: {{scope.row.StartAt}}</span><br>
                <span>完成: {{scope.row.RunAt}}</span>
              </template>
            </el-table-column>
            <el-table-column label="运行情况"  width="110" show-overflow-tooltip>
              <template #default="scope">
                <span>耗时: {{scope.row.RunSpeed}}</span><br>
                <span>进度: {{scope.row.Progress}}%</span>
              </template>
            </el-table-column>
            <el-table-column label="数据"  width="450">
              <template #default="scope">
                <span>{{friendlyJSONstringify(scope.row.Data)}}</span>
              </template>
            </el-table-column>
            <el-table-column label="客户端信息"  width="180" show-overflow-tooltip>
              <template #default="scope">
                <span>{{scope.row.Client.Name}} {{scope.row.Client.Ip}}:{{scope.row.Client.Port}}</span><br>
              </template>
            </el-table-column>
            <el-table-column label="任务状态" width="120" show-overflow-tooltip>
              <template #default="scope">
                <el-tag v-if="scope.row.Status==0" style="color:#7a7a7a">未开始</el-tag>
                <el-tag v-else-if="scope.row.Status==1">调度中</el-tag>
                <el-tag v-else-if="scope.row.Status==2" style="color:red">调度失败</el-tag>
                <el-tag v-else-if="scope.row.Status==3">执行中</el-tag>
                <el-tag v-else-if="scope.row.Status==4">失败</el-tag>
                <el-tag v-else-if="scope.row.Status==5" style="color:green">成功</el-tag>
              </template>
            </el-table-column>
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
import {defineAsyncComponent, reactive, onMounted, ref, nextTick, watch} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import {friendlyJSONstringify} from "@intlify/shared";

// 引入 api 请求接口
const serverApi = fopsApi();

// 定义变量内容
const editDialogRef = ref();
const state = reactive({
  keyWord:'',
  enable:-1,
  taskStatus:-1,
  taskGroupName:'',
  tableData: {
    data: [],
    total: 0,
    loading: false,
    param: {
      pageNum: 1,
      pageSize: 10,
    },
  },
    dialog: {
      isShowDialog: false,
      type: '',
      title: '',
      submitTxt: '',
},
});

// 监听 state.taskStatus 的变化
watch(() => state.taskStatus, (newValue, oldValue) => {
  getTableData()
});

// 初始化表格数据
const getTableData = () => {
  state.tableData.loading = true;

  const params = new URLSearchParams();
  params.append('taskStatus', state.taskStatus.toString());
  params.append('taskGroupName', state.taskGroupName);
  params.append('pageSize', state.tableData.param.pageSize.toString());
  params.append('pageIndex', state.tableData.param.pageNum.toString());
  // 请求接口
  serverApi.taskList(params.toString()).then(function (res){
    if (res.Status){
      state.tableData.data = res.Data.List;
      state.tableData.total = res.Data.RecordCount;
    }else{
      state.tableData.data=[]
    }
    state.tableData.loading = false;
  })
};

// 打开弹窗
const openDialog = (row: any) => {
  state.taskGroupName = row.Name
  state.dialog.isShowDialog = true;
  state.dialog.title = row.Name + " " +row.Caption;
  getTableData();
};
// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
const onQuery=()=>{
  getTableData();
}

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

const onStatusChange=(value:number)=>{
  state.taskStatus=value
}

// 页面加载时
onMounted(() => {
  // 等待下一次 DOM 更新后再执行代码
  // nextTick(() => {
  //   getTableData();
  // });
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
