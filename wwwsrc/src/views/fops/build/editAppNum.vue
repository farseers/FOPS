<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.appName" v-model="state.isShowDialog" width="300px">
			<div>
                <el-form-item :label="state.title">
                    <el-input-number size="default" v-model="state.dockerReplicas" controls-position="right" :min="0" style="width: 100%;"></el-input-number>
                </el-form-item>
            </div>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">确 定</el-button>
				</span>
			</template>
		</el-dialog>
    <div v-if="state.showOverlay" class="overlay">
    <div class="overlay-content">
      <img :src="Image" style="width: 200px" alt="Image">
    </div>
  </div>
	</div>
</template>

<script setup  name="fopsAppDialog" lang="ts">
import {reactive, ref, onMounted} from 'vue';
import {fopsApi} from "/@/api/fops";
import {ElMessageBox, ElMessage, ElTable} from 'element-plus';
import Image from '/@/assets/loading.gif';
const emit = defineEmits(['refresh']);
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
	isShowDialog:false,
    title:'副本数量',
    appName:'',//应用名称
    dockerReplicas:0,
    showOverlay:false
});

// 打开弹窗
const openDialog = (row: any, type: any) => {
    // console.log(row,type)
    state.dockerReplicas = row.DockerReplicas || 0;
    state.appName = row.AppName;
    state.isShowDialog = true;
};

// 关闭弹窗
const closeDialog = () => {
    state.dockerReplicas =  0;
	state.isShowDialog = false;
};
// 取消
const onCancel = () => {
	state.isShowDialog = false;
    state.dockerReplicas =  0;
};

// 提交
const onSubmit = () => {

        // 删除逻辑
        state.showOverlay = true;
        serverApi.setReplicas({
  "appName": state.appName, // # 应用名称
  "DockerReplicas": state.dockerReplicas  //# 数量
}).then(function (res){
          if (res.Status){
            closeDialog();
            ElMessage.success('修改成功');
            emit('refresh');
          }else{
            ElMessage.error(res.StatusMessage)
          }
          state.showOverlay = false;
        }).catch(() => {
        state.showOverlay = false;
      });
};


// 页面加载时
onMounted(() => {

});
// 暴露变量
defineExpose({
	openDialog,
});
</script>
<style>
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10000;
}

.overlay-content {
  text-align: center;
  color: white;
}
</style>
