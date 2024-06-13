<template>
	<div class="system-user-container layout-padding" style="overflow: auto;"  v-loading="state.loading">
		<el-row :gutter="24" class="conlyRow">
		<el-col :lg="4" :md="6" :sm="6" :xs="8" v-for="item,index in state.tableData" :key="index" class="conlyCol">
			<el-card :class="item.IsHealth?'conlyCard':'conlyCard conly_w'">
				<div class="name">
					<span>{{ item.NodeName }}<img v-show="item.OS == 'linux'" :src="linux" alt=""></span>
					
				</div>
				<div v-show="item.IsMaster">
					<el-tag  type="danger" size="small">manager</el-tag>
				</div>
				<div v-show="!item.IsMaster"><el-tag size="small">worker</el-tag></div>
				<div>{{ item.OS }} {{ item.CPUs }}核 {{ item.Memory }} {{ item.Architecture }}</div>
				<div>
					<el-tag effect="dark" size="small" style="margin-right: 5px;" :type="item.Status == 'Ready'?'success':'danger'">{{ item.Status }}</el-tag>
					<el-tag effect="dark" size="small" :type="item.Availability == 'Active'?'success':'danger'">{{ item.Availability }}</el-tag>
				</div>
				<div>{{ item.IP }}</div>
				<div>docker：{{ item.EngineVersion }}</div>
				<!-- <div class="line" v-show="item.Label && item.Label.length>0"></div> -->
				<el-tag class="ks" v-for="row,j in item.Label" :key="index.toString()+j.toString()">
					<div>{{ row.Name }} = {{ row.Value}}</div>
				</el-tag>
			</el-card>
		</el-col>
		</el-row>
	</div>
</template>

<script setup lang="ts" name="fopsTaskTimeOut">
import { defineAsyncComponent, reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import linux from '/@/assets/linux.png';
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
	loading:false,
	tableData:[]
});
// 初始化表格数据
const getTableData = () => {
	state.loading = true;
	var param={
  }
  // 请求接口
  serverApi.ColonyNodeList(param).then(function (res){
    if (res.Status){
      state.tableData = res.Data;
    }else{
		ElMessage.warning(res.StatusMessage);
    }
	state.loading = false;
  }).catch(()=>{
	state.loading = false;
  })
};

// 页面加载时
onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss">
.conlyRow{
	display: flex !important;
}
.conlyCol{
	min-width: 200px;
	margin-bottom: 15px;
}
.conlyCard .el-card__body > div{
	margin-bottom: 5px;
	text-align: center;
	font-size: 14px
}
.conlyCard {
	background-color: #f9f9e3;
	border: 1px dotted var(--el-color-primary);
	.name{
		font-weight: 700;
		span{
			display: inline-block;
			position: relative;
		}
		img{
			position: absolute;
			top: -10px;
			right: -30px;
			width: 20px;
			margin-left: 10px
		}
	}
	.line{
		width: 100%;
		height: 1px;
		background-color: var(--el-color-info);
	}
	.ks{
		text-align: center;
		border-radius: 5px;
		width: 100%;
		margin-bottom: 5px;
		div{
			padding: 3px;
		}
	}
}
.conly_w{
	background: var(--el-color-danger-light-9);
	border: 1px dotted var(--el-color-danger);
}
</style>
