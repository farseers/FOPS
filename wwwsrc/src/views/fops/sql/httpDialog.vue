<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="70%">
			<el-form ref="gitDialogFormRef" size="default" label-width="100px">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <ul class="custom-list">
              <li>
                {{state.ruleForm.CreateAt}}
                <el-tag size="small">{{state.ruleForm.HttpStatusCode}}</el-tag> {{state.ruleForm.RequestIp}} <el-tag type="success" size="small">{{state.ruleForm.HttpMethod}}</el-tag>
                <el-tag v-if="state.ruleForm.HttpContentType !=null " type="info" size="small">{{state.ruleForm.HttpContentType}}</el-tag>
                <el-tag style="font-size: 14px;" @click="copyToClipboard(state.ruleForm.HttpUrl)">{{state.ruleForm.HttpUrl}}</el-tag>
                耗时：
                  <el-tag size="small" v-if="state.ruleForm.UseTs > 100000000" type="danger">{{state.ruleForm.UseDesc}}</el-tag>
                  <el-tag size="small" v-else-if="state.ruleForm.UseTs > 50000000" type="warning">{{state.ruleForm.UseDesc}}</el-tag>
                  <el-tag size="small" v-else-if="state.ruleForm.UseTs > 1000000">{{state.ruleForm.UseDesc}}</el-tag>
                  <el-tag size="small" v-else type="success">{{state.ruleForm.UseDesc}}</el-tag>
              </li>
              <li><el-tag size="small">Headers：</el-tag>
                <el-button size="small" @click="copyToClipboard(formatJson(friendlyJSONstringify(state.ruleForm.HttpHeaders)))" type="info">复制</el-button>
                <pre>{{formatJson(friendlyJSONstringify(state.ruleForm.HttpHeaders))}}</pre>
              </li>
              <div style="display: flex; gap: 10px; list-style: none; padding: 0;">
                <li style="flex: 1;"><el-tag size="small">入参：</el-tag>
                  <el-button size="small" @click="copyToClipboard(formatJson(state.ruleForm.HttpRequestBody))" type="info">复制</el-button>
                  <pre>{{formatJson(state.ruleForm.HttpRequestBody)}}</pre>
                </li>
                <li style="flex: 1;"><el-tag size="small">出参：</el-tag>
                  <el-button size="small" @click="copyToClipboard(formatJson(state.ruleForm.HttpResponseBody))" type="info">复制</el-button>
                  <pre>{{formatJson(state.ruleForm.HttpResponseBody)}}</pre>
                </li>
              </div>
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
    HttpMethod:'',
    HttpUrl:'',
    HttpHeaders:{},
    HttpRequestBody:'',
    HttpResponseBody:'',
    HttpStatusCode:'',
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
  if (type == 2) {
      row.HttpStatusCode = row.WebStatusCode
      row.RequestIp = row.WebRequestIp
      row.HttpMethod = row.WebMethod
      row.HttpContentType = row.WebContentType
      row.HttpUrl = row.WebPath
      row.HttpHeaders = row.WebHeaders
      row.HttpRequestBody = row.WebRequestBody
      row.HttpResponseBody = row.WebResponseBody
  }

  state.dialog.title = '请求报文(TraceId：'+row.TraceId+')';
  //state.dialog.submitTxt = '修 改';
  //console.log(row2)
  state.TraceId=row.TraceId
  state.ruleForm = row


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

const formatJson = (jsonStr) =>{
  try {
    return JSON.stringify(JSON.parse(jsonStr), null, 2);
  } catch (e) {
    return jsonStr;
  }
}
const copyToClipboard = (text) =>{
      try {
        navigator.clipboard.writeText(text);
        alert('复制成功')
      } catch (err) {
        alert('复制失败')
      }
}

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

/* 增加垂直滚动条 */
pre{
  max-height: 600px; 
  overflow-y: auto;
  padding: 10px;
  background: #f8f8f8;
  border: 1px solid #eee;
  border-radius: 4px;
}
</style>
