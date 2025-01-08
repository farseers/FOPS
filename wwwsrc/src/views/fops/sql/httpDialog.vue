<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="70%">
			<el-form ref="gitDialogFormRef" size="default" label-width="100px">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <ul class="custom-list">
              <li>
                {{state.ruleForm.CreateAt}}
                <el-tag size="small">{{state.ruleForm.StatusCode}}</el-tag> {{state.ruleForm.RequestIp}} <el-tag type="success" size="small">{{state.ruleForm.Method}}</el-tag>
                {{state.ruleForm.Url}}
                整体耗时：
                  <el-tag size="small" v-if="state.ruleForm.UseTs > 100000000" type="danger">{{state.ruleForm.UseDesc}}</el-tag>
                  <el-tag size="small" v-else-if="state.ruleForm.UseTs > 50000000" type="warning">{{state.ruleForm.UseDesc}}</el-tag>
                  <el-tag size="small" v-else-if="state.ruleForm.UseTs > 1000000">{{scope.row.UseDesc}}</el-tag>
                  <el-tag size="small" v-else type="success">{{state.ruleForm.UseDesc}}</el-tag>
              </li>
              <li><el-tag size="small">Headers：</el-tag>{{friendlyJSONstringify(state.ruleForm.Headers)}}</li>
              <li><el-tag size="small">入参：</el-tag>{{state.ruleForm.RequestBody}}</li>
              <li><el-tag size="small">出参：</el-tag>{{state.ruleForm.ResponseBody}}</li>
            </ul>
					</el-col>
				</el-row>
			</el-form>


		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="fopshttpDialog">
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
    CreateAt:'',
    Method:'',
    Url:'',
    Headers:{},
    RequestBody:'',
    ResponseBody:'',
    StatusCode:'',
    UseDesc:'',
    UseTs:0,
  },
  TraceId:'',
  totalTs:0,
  Rgba:'',
  AppId:0,
  AppIp:'',
  AppName:'',
	dialog: {
		isShowDialog: false,
		type: '',
		title: '',
		submitTxt: '',
	},
});

// 打开弹窗
const openDialog = (type:number,row: any) => {
  //state.ruleForm = row;
  state.dialog.title = '请求报文(TraceId：'+row.TraceId+')';
  //state.dialog.submitTxt = '修 改';
  //console.log(row2)
  if(type==1){
    state.TraceId=row.TraceId
    state.ruleForm=row
  }else{
    state.TraceId=row.TraceId
    state.ruleForm.CreateAt=row.CreateAt
    state.ruleForm.Url=row.WebPath
    state.ruleForm.Method=row.WebMethod
    state.ruleForm.Headers=row.WebHeaders
    state.ruleForm.RequestBody=row.WebRequestBody
    state.ruleForm.ResponseBody=row.WebResponseBody
    state.ruleForm.StatusCode=row.WebStatusCode
    state.ruleForm.UseDesc=row.UseDesc
  }
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

/* 基本样式 */
.custom-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

/* 每个列表项的样式 */
.custom-list li {
  padding: 10px 15px;
  margin-bottom: 8px;
  border-radius: 5px;
  background-color: #f2f2f2;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s ease;
}

/* 悬停时的样式 */
.custom-list li:hover {
  background-color: #e0e0e0;
}

</style>
