<template>
	<div class="system-user-container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
    <div class="system-user-search mb15" style="display: flex; align-items: center; flex-wrap: wrap;">
    <label>客户端名称</label>
    <el-select class="ml10" style="max-width: 150px;margin-right: 5px;" size="small" v-model="state.appName" @change="onQuery">
      <el-option label="全部" value=""></el-option>
      <el-option v-for="item in state.appData" :label="item.AppName" :value="item.AppName"></el-option>
    </el-select>
    <el-input size="default" v-model="state.keyWord" placeholder="请输入任务名称" clearable style="max-width: 180px; margin-left: 10px;" @keyup.enter="onQuery"></el-input>
    <el-input size="default" v-model="state.clientId" placeholder="请输入客户端ID" clearable style="max-width: 180px; margin-left: 10px;" @keyup.enter="onQuery"></el-input>
    <el-select v-model="state.enable" placeholder="请选择运行状态" class="ml10" style="margin-left: 10px;" @change="onQuery">
      <el-option label="全部" :value="-1"></el-option>
      <el-option label="停止" :value="0"></el-option>
      <el-option label="启用" :value="1"></el-option>
    </el-select>
    <el-button size="default" type="primary" class="ml10" style="margin-left: 10px;" @click="onQuery">
      <el-icon>
      <ele-Search />
      </el-icon>
      查询
    </el-button>
    </div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" class="mytable">
				<el-table-column label="名称" style="line-height: 45px;height: 45px">
          <template #default="scope">
            <div style="float: left;padding-right: 10px;padding-top: 5px">
              <el-tag size="small" v-if="scope.row.Task.ScheduleStatus==0" type="info">未调度</el-tag>
              <el-tag size="small" v-else-if="scope.row.Task.ScheduleStatus==1">调度中</el-tag>
              <el-tag size="small" v-else-if="scope.row.Task.ScheduleStatus==2" type="success" style="color:green">调度成功</el-tag>
              <el-tag size="small" v-else-if="scope.row.Task.ScheduleStatus==3" type="danger">调度失败</el-tag>
              <br />
              <el-tag size="small" v-if="scope.row.Task.ExecuteStatus==0" type="info">未开始</el-tag>
              <el-tag size="small" v-else-if="scope.row.Task.ExecuteStatus==1">执行中</el-tag>
              <el-tag size="small" v-else-if="scope.row.Task.ExecuteStatus==2" type="success" style="color:green">成功</el-tag>
              <el-tag size="small" v-else-if="scope.row.Task.ExecuteStatus==3" type="danger">失败</el-tag>
            </div>
            <div style="float: left">
              <el-tag size="small" cursor="cursor" @click="onIsEnable(scope.row)" v-if="scope.row.IsEnable">启用</el-tag>
              <el-tag size="small" cursor="cursor" @click="onIsEnable(scope.row)" v-else type="info">停用</el-tag>
              <span style="margin-left: 5px">{{scope.row.Caption}}</span><br>
              <span>{{scope.row.Name}}（<span style="color:#4eb8ff">Ver:{{scope.row.Ver}}</span>）</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="StartAt" label="开始时间" width="170" show-overflow-tooltip>
                  <template #default="scope">
                    <span title="开始时间">{{scope.row.StartAt}}</span><br>
                    <span title="上次运行时间">{{scope.row.LastRunAt}}</span>
                  </template>
        </el-table-column>
            <el-table-column label="下次运行时间" width="170" show-overflow-tooltip>
              <template #default="scope">
                <span title="Cron表达式" >{{scope.row.ConsumerRoutingKeyon}}</span><br>
                <span style="color:red" title="下次运行时间" v-if="compareTime(scope.row.NextAt)" > {{scope.row.NextAt}}</span>
                <span title="下次运行时间" v-else> {{scope.row.NextAt}}</span>
              </template>
            </el-table-column>
            <el-table-column label="运行情况" width="160" show-overflow-tooltip>
              <template #default="scope">
                <span>耗时: {{scope.row.RunSpeedAvg}} ms</span><br>
                <span>运行: {{scope.row.RunCount}} 次</span>
              </template>
            </el-table-column>
            <el-table-column label="数据">
              <template #default="scope">
                <div>{{friendlyJSONstringify(scope.row.Data)}}</div>
                <span v-for="(item, index) in scope.row.Clients.slice(0, 3)" :key="index">
                  <el-tag v-if="item.IsMaster" size="small" style="margin-right: 5px;">主 {{item.Name}} {{item.Ip}}:{{item.Port}}</el-tag>
                  <el-tag v-else size="small" type="info" style="margin-right: 5px;">{{item.Name}} {{item.Ip}}:{{item.Port}}</el-tag>
                </span>
              </template>
            </el-table-column>
				<el-table-column label="操作" width="140">
					<template #default="scope">
            <el-button size="small" text type="danger" @click="onTaskList(scope.row)">历史</el-button>
            <el-button size="small" text type="danger" @click="onLog(scope.row)">日志</el-button><br />
            <el-button size="small" text type="warning" @click="onEdit('edit',scope.row)">修改</el-button>
            <el-button size="small" text type="success" @click="onExecuteNow(scope.row)">执行</el-button>
            <el-button size="small" text type="info" @click="onDel(scope.row)">删除</el-button>
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
    <editDialog ref="editDialogRef" @refresh="getTableData()" />
    <taskDialog ref="taskDialogRef" @refresh="getTableData()" />
    <logDialog ref="logDetailDialogRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts" name="fopsTask">
import {defineAsyncComponent, reactive, onMounted, ref, nextTick, watch} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import {friendlyJSONstringify} from "@intlify/shared";

// 引入 api 请求接口
const serverApi = fopsApi();

// 引入组件
const editDialog = defineAsyncComponent(() => import('/src/views/fops/task/editGroupDialog.vue'));
const taskDialog = defineAsyncComponent(() => import('/src/views/fops/task/taskDialog.vue'));
const logDialog = defineAsyncComponent(() => import('/src/views/fops/task/logDialog.vue'));


// 定义变量内容
const editDialogRef = ref();
const taskDialogRef = ref();
const logDetailDialogRef = ref();
const state = reactive({
  keyWord:'',
  appName:'',
  enable:1,
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
  NowTime:new Date(),
  appData:[],
});

// 监听 state.enable 的变化
watch(() => state.enable, (newValue, oldValue) => {
  getTableData()
});

watch(() => state.appName, (newValue, oldValue) => {
  getTableData()
});

const getAppData=()=>{
  serverApi.dropDownList({}).then(function (res){
    if (res.Status){
      state.appData=res.Data
    }else{
      state.appData=[]
    }
  })
}

// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
  const params = new URLSearchParams();
  params.append('clientName', state.appName);
  params.append('taskGroupName', state.keyWord);
  params.append('enable', state.enable.toString());
  params.append('taskStatus', state.taskStatus.toString());
  params.append('clientId', state.clientId);
  params.append('pageSize', state.tableData.param.pageSize.toString());
  params.append('pageIndex', state.tableData.param.pageNum.toString());

  // 请求接口
  serverApi.taskGroupList(params.toString()).then(function (res){
    if (res.Status){
      state.tableData.data = res.Data.List;
      state.tableData.total = res.Data.RecordCount;
    }else{
      state.tableData.data=[]
    }
    state.tableData.loading = false;
  })
};

const compareTime=(nextAt:any)=>{
  var convertedTime = new Date(nextAt)
  return convertedTime.getTime() < new Date().getTime();
}
const onQuery=()=>{
  getTableData();
}
const onEdit=(type: string, row: any)=>{
  editDialogRef.value.openDialog(type, row);
}
const onTaskList=(row: any)=>{
  taskDialogRef.value.openDialog(row);
}
const onLog=(row: any)=>{
  logDetailDialogRef.value.openDialog(row);
}

// 删除
const onDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除：“${row.Name}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	}).then(() => {
      // 删除逻辑
      serverApi.taskDel({"taskGroupName":row.Name}).then(function (res){
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

 const onExecuteNow = (row: any) => {
  ElMessageBox.confirm(`是否立即执行当前任务：“${row.Name}”？`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    serverApi.taskGroupExecuteNow({ "taskGroupName": row.Name }).then(function (res) {
    if (res.Status) {
      ElMessage.success('执行成功');
      getTableData();
    } else {
      ElMessage.error(res.StatusMessage);
    }
    });
  })
  .catch(() => {});
};

//启用停用
const onIsEnable=(row: any)=>{
  var setEnable=row.IsEnable
  var tips=""
  if(setEnable){
    setEnable=false
    tips="停用"
  }else{
    setEnable=true
    tips="启用"
  }

  ElMessageBox.confirm(`该任务即将：“${tips}”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        // 设置状态
        serverApi.taskGroupSetEnable({"taskGroupName":row.Name,"enable":setEnable}).then(function (res){
          if (res.Status){
            getTableData();
            if(setEnable){
              ElMessage.success('启用-成功');
            }else{
              ElMessage.success('停用-成功');
            }

          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
      })
      .catch(() => {});
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
const onEnableChange=(value:number)=>{
  state.enable=value
}
// 页面加载时
onMounted(() => {
  // 等待下一次 DOM 更新后再执行代码
  nextTick(() => {
    getTableData();
    getAppData();
  });
});
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
