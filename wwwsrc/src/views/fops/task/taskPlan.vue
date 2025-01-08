<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
      <el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" class="mytable">
        <el-table-column label="名称" style="line-height: 45px;height: 45px">
          <template #default="scope">
            <div style="float: left">
              <span>{{scope.row.Caption}}</span><br>
              <span>{{scope.row.Name}}（<span style="color:#4eb8ff">Ver:{{scope.row.Ver}}</span>）</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="Plan" label="计划时间" width="190" show-overflow-tooltip>
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.Plan.includes(`等待调度`)">{{scope.row.Plan}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Plan.includes(`等待`)" type="info">{{scope.row.Plan}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Plan.includes(`超时`)" type="danger">{{scope.row.Plan}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Plan.includes(`已执行`)" type="success">{{scope.row.Plan}}</el-tag>
            <br><span>{{scope.row.RunAt}}</span>
          </template>
        </el-table-column>
        <el-table-column label="数据">
          <template #default="scope">
            <span>{{friendlyJSONstringify(scope.row.Data)}}</span>
          </template>
        </el-table-column>
        <el-table-column label="客户端信息" width="180" show-overflow-tooltip>
          <template #default="scope">
            <div>
              <el-tag v-if="scope.row.Client.Name != ''" size="small">{{scope.row.Client.Name}} {{scope.row.Client.Ip}}:{{scope.row.Client.Port}}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button v-if="scope.row.Status !=0" size="small" text type="danger" @click="onKill(scope.row)">停止任务</el-button>
          </template>
        </el-table-column>
      </el-table>
		</el-card>
	</div>
</template>

<script setup lang="ts" name="fopsTaskRunning">
import {defineAsyncComponent, reactive, onMounted, ref, onUnmounted} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import {friendlyJSONstringify} from "@intlify/shared";

// 引入 api 请求接口
const serverApi = fopsApi();

// 定义变量内容
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
	const params = new URLSearchParams();
  params.append('top', '20');

  // 请求接口
  serverApi.taskPlanList(params.toString()).then(function (res){
    if (res.Status){
      state.tableData.data = res.Data;
    }else{
      state.tableData.data=[]
    }
    state.tableData.loading = false;
  })
};

// 停止任务
const onKill = (row: any) => {
  ElMessageBox.confirm(`准备停止任务：“${row.Name}”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 删除逻辑
        serverApi.killTask({"taskGroupName":row.Name}).then(function (res){
          if (res.Status){
            getTableData();
            ElMessage.success('停止成功');
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
};

let intervalTableDataId = null;

// 页面加载时
onMounted(() => {
  state.tableData.loading = true; // 首次加载时，需要
	getTableData();
  intervalTableDataId = setInterval(getTableData, 1000);
});

// 页面注销的时候
onUnmounted(()=>{
  clearInterval(intervalTableDataId);
})

</script>

<style lang="scss">
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
.el-table tr td {
  /* 你的自定义样式 */
  padding: 0 0!important;
}
</style>
