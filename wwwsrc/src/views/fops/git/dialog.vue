<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px">
			<el-form ref="gitDialogFormRef" :model="state.ruleForm" size="default" label-width="120px">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="Git名称">
							<el-input v-model="state.ruleForm.Name" placeholder="请输入Git名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="托管地址">
							<el-input v-model="state.ruleForm.Hub" placeholder="请输入托管地址" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="Git分支">
							<el-input v-model="state.ruleForm.Branch" placeholder="请输入Git分支" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="账户名称">
							<el-input v-model="state.ruleForm.UserName" placeholder="请输入账户名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="账户密码">
						<el-input v-model="state.ruleForm.UserPwd" placeholder="请输入"  clearable></el-input> <!--	type="password"-->
						</el-form-item>
					</el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="存储目录">
              <el-input v-model="state.ruleForm.Path" placeholder="请输入存储目录" clearable></el-input>
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="是否为应用仓库">
              <el-switch v-model="state.ruleForm.IsApp" inline-prompt active-text="是" inactive-text="否"></el-switch>
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

<script setup lang="ts" name="fopsGitDialog">
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
    Name: '', // Git名称
    Hub: '', // 托管地址
    Branch: 'main', // Git分支
    UserName: '', // 账户名称
    UserPwd: '', // 账户密码
    Path: '', // 存储目录
    PullAt: '', // 拉取时间
    IsApp:false, // 是否为应用仓库

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
		state.dialog.title = '修改Git';
		state.dialog.submitTxt = '修 改';
    // 绑定数据
    state.ruleForm.Id=row.Id
    state.ruleForm.Name=row.Name
    state.ruleForm.Hub=row.Hub
    state.ruleForm.Branch=row.Branch
    state.ruleForm.UserName=row.UserName
    state.ruleForm.UserPwd=row.UserPwd
    state.ruleForm.Path=row.Path
    state.ruleForm.PullAt=row.PullAt
    state.ruleForm.IsApp=row.IsApp


	} else {
		state.dialog.title = '新增Git';
		state.dialog.submitTxt = '新 增';

    state.ruleForm.Id=0
    state.ruleForm.Name=""
    state.ruleForm.Hub=""
    state.ruleForm.Branch="main"
    state.ruleForm.UserName=""
    state.ruleForm.UserPwd=""
    state.ruleForm.Path=""
    state.ruleForm.PullAt=""
		// 清空表单，此项需加表单验证才能使用
		// nextTick(() => {
		// 	gitDialogFormRef.value.resetFields();
		// });
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
    "Id":state.ruleForm.Id,
    "Name":state.ruleForm.Name,
    "Hub":state.ruleForm.Hub,
    "Branch":state.ruleForm.Branch,
    "UserName":state.ruleForm.UserName,
    "UserPwd":state.ruleForm.UserPwd,
    "Path":state.ruleForm.Path,
    "IsApp":state.ruleForm.IsApp
  }

	if (state.dialog.type === 'add') {
    serverApi.gitAdd(param).then(async function(res){
      if(res.Status){
        ElMessage.success("添加成功")
        closeDialog();
        emit('refresh');
      }else{
        ElMessage.error(res.StatusMessage)
      }
    })
  }else if (state.dialog.type=='edit'){
    serverApi.gitEdit(param).then(async function(res){
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
  onSubmit
});
</script>

<style>
textarea{
  height: 220px;
}
</style>
