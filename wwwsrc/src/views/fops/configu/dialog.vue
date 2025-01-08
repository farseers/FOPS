<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.title" v-model="state.isShowDialog" width="500px">
			<el-form ref="gitDialogFormRef" size="default" label-width="120px">
				<el-form-item label="应用名称">
						<el-input v-model="state.AppName" placeholder="请输入应用名称" clearable></el-input>
                </el-form-item>
                <el-form-item label="键">
						<el-input v-model="state.Key" placeholder="请输入Key" clearable></el-input>
                </el-form-item>
                <el-form-item label="值">
						<el-input v-model="state.Value" placeholder="请输入Value" clearable></el-input>
                </el-form-item>
			</el-form>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">保存</el-button>
				</span>
			</template>
		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="fopsGitDialog">
import { reactive, ref } from 'vue';
import {fopsApi} from "/@/api/fops";
import {ElMessage } from 'element-plus';
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);
// 定义变量内容
const gitDialogFormRef = ref();
const state = reactive({
	isShowDialog:false,
    title:'',
    AppName: "",
    Key: "",
    Value: "",
    isEdit:false,
});

// 打开弹窗
const openDialog = ( row: any) => {
    // console.log(row)
  if(row){
    state.title = '编辑应用';
    state.AppName = row.AppName;
    state.Key = row.Key;
    state.Value = row.Value;
    state.isEdit = true;
  }else{
    state.AppName = '';
    state.Key = '';
    state.Value = '';
    state.title = '添加应用'
    state.isEdit = false;
  }
  state.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
	state.isShowDialog = false;
};
// 取消
const onCancel = () => {
	closeDialog();
};
// 提交
const onSubmit = () => {

  // 提交数据
  const param={
    AppName: state.AppName,
    Key: state.Key,
    Value: state.Value,
  }

	if (state.isEdit) {
    serverApi.configureUpdate(param).then(async function(res){
      if(res.Status){
        ElMessage.success("保存成功")
        closeDialog();
        emit('refresh');
      }else{
        ElMessage.error(res.StatusMessage)
      }
    })
  }else{
    serverApi.configureAdd(param).then(async function(res){
      if(res.Status){
        ElMessage.success("添加成功")
        closeDialog();
        emit('refresh');
      }else{
        ElMessage.error(res.StatusMessage)
      }
    })
  }
};
// 暴露变量
defineExpose({
	openDialog
});
</script>

