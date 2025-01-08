<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px">
			<el-form ref="gitDialogFormRef" :model="state.ruleForm" size="default" label-width="100px">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="任务名称">
							<el-input v-model="state.ruleForm.Name" placeholder="请输入任务名称" readonly></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="版本">
							<el-input v-model="state.ruleForm.Ver" placeholder="请输入版本" readonly></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="任务标题">
							<el-input v-model="state.ruleForm.Caption" placeholder="请输入任务标题" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="参数">
							<el-input type="textarea" maxlength="500" show-word-limit resize="none" :rows="5" class="textarea-box" v-model="state.ruleForm.Data" placeholder="请输入参数" clearable></el-input>
						</el-form-item>
					</el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="开始时间">
              <el-date-picker v-model="state.ruleForm.StartAt" type="datetime" placeholder="请选择开始时间" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%"
              ></el-date-picker>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="下次执行时间">
              <el-date-picker v-model="state.ruleForm.NextAt" type="datetime" placeholder="请选择下次执行时间" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%"></el-date-picker>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="设置Cron">
              <el-input v-model="state.ruleForm.Cron" placeholder="请输入Cron时间格式" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="是否启用">
              <el-switch v-model="state.ruleForm.IsEnable" inline-prompt active-text="启用" inactive-text="禁用"></el-switch>
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

<script setup lang="ts" name="fopsTaskEditDialog">
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
    Name: '', // 任务名称
    Ver: 1, // 版本
    Caption: '', // 任务标题
    Data: '', // 仓库名称
    StartAt: '', // 开始时间
    NextAt: '', // 下次执行时间
    Cron: '', // 时间
    IsEnable: false, // 是否启用
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
		state.dialog.title = '修改任务';
		state.dialog.submitTxt = '修 改';

    // 详情
    var url='/basicapi/taskGroup/info-'+row.Name
    serverApi.taskGroupInfo(url).then(function (res){
      if (res.Status){
        var row=res.Data
        // 绑定数据
        state.ruleForm.Name=row.Name
        state.ruleForm.Ver=row.Ver
        state.ruleForm.Caption=row.Caption
        state.ruleForm.Data=friendlyJSONstringify(row.Data)
        state.ruleForm.StartAt=row.StartAt
        state.ruleForm.NextAt=row.NextAt
        state.ruleForm.Cron=row.Cron
        state.ruleForm.IsEnable=row.IsEnable
      }
    })

	} else {
		state.dialog.title = '新增任务';
		state.dialog.submitTxt = '新 增';

    state.ruleForm.Name=""
    state.ruleForm.Ver=0
    state.ruleForm.Caption=""
    state.ruleForm.Data=''
    state.ruleForm.StartAt=""
    state.ruleForm.NextAt=""
    state.ruleForm.Cron=""
    state.ruleForm.IsEnable=false
		// 清空表单，此项需加表单验证才能使用
		// nextTick(() => {
		// 	gitDialogFormRef.value.resetFields();
		// });
	}
	state.dialog.isShowDialog = true;
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

  // const originalObject = JSON.parse();
  // const modifiedObject = {};
  //
  // // 循环遍历原始对象，为每个键和值加上双引号并构建新的对象
  // for (const key in originalObject) {
  //   if (originalObject.hasOwnProperty(key)) {
  //     const quotedKey = `"${key}"`;
  //     const quotedValue = `${originalObject[key]}`;
  //     modifiedObject[quotedKey] = quotedValue;
  //   }
  // }

  // 提交数据
  var param={
    "Name":state.ruleForm.Name,
    "Ver":state.ruleForm.Ver,
    "Caption":state.ruleForm.Caption,
    "Data":JSON.parse(state.ruleForm.Data),
    "StartAt":state.ruleForm.StartAt,
    "NextAt":state.ruleForm.NextAt,
    "Cron":state.ruleForm.Cron,
    "IsEnable":state.ruleForm.IsEnable
  }

	if (state.dialog.type === 'add') {
    // serverApi.taskUpdate(param).then(function (res){
    //   if(res.Status){
    //     ElMessage.success("添加成功")
    //     closeDialog();
    //     emit('refresh');
    //   }else{
    //     ElMessage.error(res.StatusMessage)
    //   }
    // })

  }else if (state.dialog.type=='edit'){
    serverApi.taskUpdate(param).then(function (res){
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

// 暴露变量
defineExpose({
	openDialog,
});
</script>
