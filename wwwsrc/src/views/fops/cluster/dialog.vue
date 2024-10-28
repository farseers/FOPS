<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="600px">
			<el-form ref="gitDialogFormRef" :model="state.ruleForm" size="default" label-width="120px">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb22">
						<el-form-item label="集群名称">
							<el-input v-model="state.ruleForm.Name" placeholder="请输入集群名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="集群地址">
							<el-input v-model="state.ruleForm.FopsAddr" placeholder="请输入集群地址" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="调度中心地址">
							<el-input v-model="state.ruleForm.FScheduleAddr" placeholder="请输入调度中心地址，例如：https://fschedule.xxx.com/" clearable></el-input>
						</el-form-item>
					</el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="Docker网络">
              <el-input v-model="state.ruleForm.DockerNetwork" placeholder="请输入Docker网络名称" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="DockerHub地址">
              <el-input v-model="state.ruleForm.DockerHub" placeholder="请输入托管地址" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="DockerHub账户">
              <el-input v-model="state.ruleForm.DockerUserName" placeholder="请输入账户名称" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="DockerHub密码">
              <el-input v-model="state.ruleForm.DockerUserPwd" placeholder="请输入账户密码" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="本地集群">
              <el-switch v-model="state.ruleForm.IsLocal" inline-prompt active-text="是" inactive-text="否"></el-switch>
            </el-form-item>
          </el-col>
				</el-row>
			</el-form>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{ state.dialog.submitTxt }}</el-button>
				</span>
			</template>
		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="fopsClusterDialog">
import { reactive, ref } from 'vue';
import {fopsApi} from "/@/api/fops";
import { ElMessageBox, ElMessage } from 'element-plus';
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);

// 定义变量内容
const gitDialogFormRef = ref();
const state = reactive({
	ruleForm: {
    Id:0, //编号
    Name: '', // 集群名称
    FopsAddr: '', // 集群地址
    FScheduleAddr: '', // K8s连接配置
    DockerNetwork: '', // Docker网络
    DockerHub: '', // 托管地址
    DockerUserName: '', // 账户名称
    DockerUserPwd: '', // 账户密码
    IsLocal: false, // 本地集群
	},
	dialog: {
		isShowDialog: false,
		type: '',
		title: '',
		submitTxt: '',
	},
});

// 打开弹窗
const openDialog = (type: string, row: any) => {
  state.dialog.type=type
	if (type === 'edit') {
		state.ruleForm = row;
		state.dialog.title = '修改集群';
		state.dialog.submitTxt = '修 改';
    // 绑定数据
    state.ruleForm.Id=row.Id
    state.ruleForm.Name=row.Name
    state.ruleForm.FopsAddr=row.FopsAddr
    state.ruleForm.FScheduleAddr=row.FScheduleAddr
    state.ruleForm.DockerNetwork=row.DockerNetwork
    state.ruleForm.DockerHub=row.DockerHub
    state.ruleForm.DockerUserName=row.DockerUserName
    state.ruleForm.DockerUserPwd=row.DockerUserPwd
    state.ruleForm.IsLocal=row.IsLocal

	} else {
		state.dialog.title = '新增集群';
		state.dialog.submitTxt = '新 增';

    state.ruleForm.Id=0
    state.ruleForm.Name=""
    state.ruleForm.FopsAddr=""
    state.ruleForm.FScheduleAddr=""
    state.ruleForm.DockerNetwork=""
    state.ruleForm.DockerHub=""
    state.ruleForm.DockerUserName=""
    state.ruleForm.DockerUserPwd=""
    state.ruleForm.IsLocal=false
	}
	state.dialog.isShowDialog = true;
	getMenuData();
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};
// 取消
const onCancel = () => {
	closeDialog();
};
// 提交
const onSubmit = () => {
  // 提交数据
  var param={
    "Id": state.ruleForm.Id,
    "Name": state.ruleForm.Name,
    "FopsAddr":state.ruleForm.FopsAddr,
    "FScheduleAddr":state.ruleForm.FScheduleAddr,
    "DockerNetwork":state.ruleForm.DockerNetwork,
    "DockerHub":state.ruleForm.DockerHub,
    "DockerUserName":state.ruleForm.DockerUserName,
    "DockerUserPwd":state.ruleForm.DockerUserPwd,
    "IsLocal":state.ruleForm.IsLocal,
  }

	if (state.dialog.type === 'add') {
    serverApi.clusterAdd(param).then(function (res){
      if(res.Status){
        ElMessage.success("添加成功")
        closeDialog();
        emit('refresh');
      }else{
        ElMessage.error(res.StatusMessage)
      }
    })

  }else if (state.dialog.type=='edit'){
    serverApi.clusterEdit(param).then(function (res){
      if(res.Status){
        ElMessage.success("修改成功")
        closeDialog();
        emit('refresh');
      }else{
        ElMessage.error(res.StatusMessage)
      }
    })

  }
};
// 初始化部门数据
const getMenuData = () => {
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>
