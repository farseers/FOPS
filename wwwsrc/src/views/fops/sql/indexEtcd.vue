<template>
	<div class="system-user-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
        <label>TraceId</label>
        <el-input class="ml5" size="default" v-model="state.traceId" placeholder="链路ID" clearable style="max-width: 165px;"> </el-input>
        <label class="ml5">应用</label>
        <el-select class="ml5" style="max-width: 110px;" size="small" v-model="state.appName">
          <el-option label="全部" value=""></el-option>
          <el-option v-for="item in state.appData" :label="item.AppName" :value="item.AppName" ></el-option>
        </el-select>
        <label class="ml10">执行端IP</label>
        <el-input class="ml5" size="default" v-model="state.appIp" placeholder="执行端IP" clearable style="max-width: 120px;"></el-input>
        <label class="ml10">监控KEY</label>
        <el-input class="ml5" size="default" v-model="state.key" placeholder="监控KEY" clearable style="max-width: 180px;"> </el-input>
        <label class="ml10">leaseID</label>
        <el-input class="ml5" size="default" v-model="state.leaseID" placeholder="leaseID" clearable style="max-width: 180px;"> </el-input>
        <label class="ml10">耗时最高</label>
        <el-select class="ml5" v-model="state.startMin" placeholder="往前推N分钟的数据" style="max-width: 120px;" size="default">
          <el-option label="全部" :value="0"></el-option>
          <el-option label="1小时" :value="60"></el-option>
          <el-option label="30分钟" :value="30"></el-option>
          <el-option label="10分钟" :value="10"></el-option>
          <el-option label="5分钟" :value="5"></el-option>
          <el-option label="1分钟" :value="1"></el-option>
        </el-select>
        <label class="ml10">执行时间</label>
        <el-input class="ml5" size="default" v-model="state.searchUseTs" placeholder="执行时间大于毫秒的记录" style="max-width: 80px;"> </el-input> ms
        <el-checkbox v-model="state.onlyViewException" label="仅看异常" size="small" class="ml5" style="color:#ff5000;"/>
				<el-button size="default" type="primary" class="ml10" @click="onQuery">
					<el-icon>
						<ele-Search />
					</el-icon>
					查询
				</el-button>
        <el-button size="default" type="warning" class="ml5" @click="linkTraceDeleteSlow">
					<el-icon><ele-Delete /></el-icon>
					删除七天前数据
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
        <el-table-column width="180px" label="TraceID" show-overflow-tooltip>
          <template #default="scope">
            <span @click="onDetail(scope.row)">{{scope.row.TraceId}}</span>
          </template>
        </el-table-column>
        <el-table-column width="200px" label="应用" show-overflow-tooltip>
          <template #default="scope">
            <el-tag size="small">{{scope.row.AppName}} {{scope.row.AppIp}}</el-tag><br>
            {{scope.row.AppId}}
          </template>
        </el-table-column>
        <el-table-column width="120px" prop="UseDesc" label="执行耗时" show-overflow-tooltip>
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.UseTs > 100000000" type="danger">{{scope.row.UseDesc}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.UseTs > 50000000" type="warning">{{scope.row.UseDesc}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.UseTs > 1000000">{{scope.row.UseDesc}}</el-tag>
            <el-tag size="small" v-else type="success">{{scope.row.UseDesc}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="EtcdKey" label="监控KEY" show-overflow-tooltip></el-table-column>
        <el-table-column prop="EtcdLeaseID" label="LeaseID" show-overflow-tooltip></el-table-column>
        <el-table-column width="200px" label="异常" show-overflow-tooltip>
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.Exception!=null" type="danger">{{scope.row.Exception.ExceptionCallFile}}:{{scope.row.Exception.ExceptionCallLine}} {{scope.row.Exception.ExceptionCallFuncName}}</el-tag><br  v-if="scope.row.Exception!=null">
            <el-tag size="small" v-if="scope.row.Exception!=null" type="danger">{{scope.row.Exception.ExceptionMessage}}</el-tag>
            <el-tag size="small" v-else type="info">无</el-tag>
          </template>
        </el-table-column>
        <el-table-column width="180px" prop="CreateAt" label="执行时间" show-overflow-tooltip></el-table-column>
				<el-table-column label="操作" width="100">
					<template #default="scope">
						<el-button size="small" text type="primary" @click="onDetail(scope.row)">追踪</el-button>
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
    <detailDialog ref="detailDialogRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts" name="fopsTaskRunning">
import { defineAsyncComponent, reactive, onMounted, ref , watch} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import {friendlyJSONstringify} from "@intlify/shared";

// 引入 api 请求接口
const serverApi = fopsApi();
// 引入组件
const detailDialog = defineAsyncComponent(() => import('/src/views/fops/linkTrace/detailV2Dialog.vue'));


// 定义变量内容
const detailDialogRef = ref();
const state = reactive({
  appName:'',
  traceId:'',
  appIp:'',
  key:'',
  leaseID:0,
  startMin:0,
  searchUseTs:0,
  onlyViewException:false,
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			pageNum: 1,
      pageSize: 20,
		},
	},  appData:[],
});
// 监听 state.startMin 的变化
watch(() => state.startMin, (newValue, oldValue) => {
  console.log(`count 从 ${oldValue} 变为 ${newValue}`);
  getTableData()
});
watch(() => state.appName, (newValue, oldValue) => {
  console.log(`count 从 ${oldValue} 变为 ${newValue}`);
  getTableData()
});
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;

  var data={
    traceId:state.traceId,
    appName:state.appName,
    appIp:state.appIp,
    key:state.key,
    leaseID:state.leaseID.toString(),
    startMin:state.startMin.toString(),
    searchUseTs:state.searchUseTs.toString(),
    onlyViewException:state.onlyViewException,
    pageSize:state.tableData.param.pageSize.toString(),
    pageIndex:state.tableData.param.pageNum.toString(),
  }
  const params = new URLSearchParams(data).toString();
  // 请求接口
  serverApi.slowEtcd(params).then(function (res){
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
  row.tid = row.TraceId;
  detailDialogRef.value.openDialog(row);
}
const getAppData=()=>{
  serverApi.dropDownList({}).then(function (res){
    if (res.Status){
      state.appData=res.Data
    }else{
      state.appData=[]
    }
  })
}
// 删除用户
const onDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除：“${row.Name}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
      // 删除逻辑
      serverApi.taskDel({"TaskGroupId":row.Id}).then(function (res){
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
const linkTraceDeleteSlow = ()=>{
  ElMessageBox.confirm(`删除七天前的数据，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			serverApi.linkTraceDeleteSlow({dbName:'Etcd'}).then((res)=>{
				if(res.Status){
					onQuery();
					ElMessage.success('删除成功');
				}else{
				ElMessage.error(res.StatusMessage);
			}
			})
			
		})
		.catch(() => {});
}
const onQuery=()=>{
  state.tableData.param.pageNum = 1;
  getTableData();
}
// 页面加载时
onMounted(() => {
  getAppData();
  //getTableData();
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
