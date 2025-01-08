<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px">
			<el-form ref="gitDialogFormRef" :model="state.ruleForm" size="default" label-width="100px">
				<el-row :gutter="35">
					<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="跟踪ID">
							<el-input v-model="state.ruleForm.TraceId" placeholder="跟踪ID" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="应用ID">
							<el-input v-model="state.ruleForm.AppId" placeholder="应用ID" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="应用名称">
							<el-input v-model="state.ruleForm.AppName" placeholder="应用名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="应用IP">
							<el-input v-model="state.ruleForm.AppIp" placeholder="应用IP" clearable></el-input>
						</el-form-item>
					</el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="上级应用名称">
              <el-input v-model="state.ruleForm.ParentAppName" placeholder="上级应用名称" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="开始时间（微秒）">
              <el-input v-model="state.ruleForm.StartTs" placeholder="开始时间（微秒）" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="结束时间（微秒）">
              <el-input v-model="state.ruleForm.EndTs" placeholder="结束时间（微秒）" clearable></el-input>
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="耗时（微秒）">
              <el-input v-model="state.ruleForm.UseTs" placeholder="耗时（微秒）" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="链路类型">
              <el-input v-model="state.ruleForm.TraceType" placeholder="链路类型" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="异常信息">
              <el-tag size="small">ExceptionCallFile:{{state.ruleForm.Exception.ExceptionCallFile}}</el-tag>
              <el-tag size="small">ExceptionCallLine:{{state.ruleForm.Exception.ExceptionCallLine}}</el-tag>
              <el-tag size="small">ExceptionCallFuncName:{{state.ruleForm.Exception.ExceptionCallFuncName}}</el-tag>
              <el-tag size="small">ExceptionIsException:{{state.ruleForm.Exception.ExceptionIsException}}</el-tag>
              <el-tag size="small">ExceptionMessage:{{state.ruleForm.Exception.ExceptionMessage}}</el-tag>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="输出日志">
              <el-input v-model="state.ruleForm.List" placeholder="输出日志" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="请求的域名">
              <el-input v-model="state.ruleForm.WebDomain" placeholder="请求的域名" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="请求的地址">
              <el-input v-model="state.ruleForm.WebPath" placeholder="请求的地址" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="请求方式">
              <el-input v-model="state.ruleForm.WebMethod" placeholder="请求方式" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="请求类型">
              <el-input v-model="state.ruleForm.WebContentType" placeholder="请求类型" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="状态码">
              <el-input v-model="state.ruleForm.WebStatusCode" placeholder="状态码" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="header">
              <el-tag size="small">{{friendlyJSONstringify(state.ruleForm.WebHeaders)}}</el-tag>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="入参">
              <el-input v-model="state.ruleForm.WebRequestBody" placeholder="入参" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="出参">
              <el-input v-model="state.ruleForm.WebResponseBody" placeholder="出参" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="客户端请求IP">
              <el-input v-model="state.ruleForm.WebRequestIp" placeholder="客户端请求IP" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="消费服务">
              <el-input v-model="state.ruleForm.ConsumerServer" placeholder="消费服务" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="消费队列名称">
              <el-input v-model="state.ruleForm.ConsumerQueueName" placeholder="消费队列名称" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="消费路由key">
              <el-input v-model="state.ruleForm.ConsumerRoutingKey" placeholder="消费路由key" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="任务名称">
              <el-input v-model="state.ruleForm.TaskName" placeholder="任务名称" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="任务组ID">
              <el-input v-model="state.ruleForm.TaskGroupId" placeholder="任务组ID" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="任务ID">
              <el-input v-model="state.ruleForm.TaskId" placeholder="任务ID" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="WatchKey">
              <el-input v-model="state.ruleForm.WatchKey" placeholder="WatchKey" clearable></el-input>
            </el-form-item>
          </el-col>

				</el-row>
			</el-form>

		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="fopsLinkDetailDialog">
import { reactive, ref } from 'vue';
import {fopsApi} from "/@/api/fops";
import { ElMessageBox, ElMessage } from 'element-plus';
import {friendlyJSONstringify} from "@intlify/shared";
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);

// 定义变量内容
const gitDialogFormRef = ref();
const state = reactive({
	ruleForm: {
    TraceId:0,
    AppId:0,
    AppName:'',
    AppIp:'',
    ParentAppName:'',
    StartTs:0,
    EndTs:0,
    UseTs:0,
    TraceType:0,
    Exception:'',
    List:'',
    WebDomain:'',
    WebPath:'',
    WebMethod:'',
    WebContentType:'',
    WebStatusCode:0,
    WebHeaders: {},
    WebRequestBody:'',
    WebResponseBody:'',
    WebRequestIp:'',
    ConsumerServer:'',
    ConsumerQueueName:'',
    ConsumerRoutingKey:'',
    TaskName:'',
    TaskGroupId:0,
    TaskId:0,
    WatchKey:'',
	},
	dialog: {
		isShowDialog: false,
		type: '',
		title: '',
		submitTxt: '',
	},
});

// 打开弹窗
const openDialog = (row2: any) => {
  //state.ruleForm = row;
  state.dialog.title = '链路追踪详情';
  //state.dialog.submitTxt = '修 改';
  // 详情
  serverApi.linkTraceInfo(row2.tid).then(function (res){
    if (res.Status){
      // 绑定数据
      state.ruleForm=res.Data
    }
  })
	state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};

const getStatusDesc=(status:number)=>{
  switch (status){
    case 0:
      return "未开始"
    case 1:
      return "调度中"
    case 2:
      return "调度失败"
    case 3:
      return "执行中"
    case 4:
      return "失败"
    case 5:
      return "成功"
  }
  return ""
}
// 取消
const onCancel = () => {
	closeDialog();
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style>
textarea{
  height: 220px;
}
</style>
