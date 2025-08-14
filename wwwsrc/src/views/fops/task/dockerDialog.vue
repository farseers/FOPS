<template>
    <div>
        <el-dialog title="容器日志" v-model="state.isShowDockerLogDialog" style="width: 80%;top:20px;margin-bottom: 50px;">
				<div style="display: flex;flex-flow: column;max-height: calc(90vh - 151px) !important;">
					<div>
					<el-tag size="default" :type="item.Id == state.dockerLog.Id ? '' : 'info'"
						@click="clickDockerLog(item)" v-for="item in state.dockerLogContent" :key="item.Id"
						style="cursor: pointer;margin:0 15px 5px 0">
						{{ item.Name }}（{{ item.Node }}）
					</el-tag>
				</div>
				<div style="margin: 5px 0;">
					<div>
						<el-tag size="small" type="success" style="margin-right:10px">{{ state.dockerLog.State
							}}</el-tag>
						<el-tag size="small" type="success" style="margin-right:10px">{{ state.dockerLog.StateInfo
							}}</el-tag>
						<el-tag size="small" type="success" style="margin-right:10px">{{ state.dockerLog.Image
							}}</el-tag>
					</div>
					<div style="color: #f56c6c;">{{ state.dockerLog.Error }}</div>
				</div>
				<div class="layout-padding-auto" style="background-color:#393d49;flex: 1;overflow: auto;" ref="scrollableDockerLog">
					<pre v-html="state.dockerLog.Log"
							style="color: #fff;background-color:#393d49;padding: 5px 0 5px 5px;"></pre>
				</div>
				</div>
				
			</el-dialog>
    </div>
</template>
<script setup name="dockerDialog">
import { reactive, defineExpose, ref } from 'vue';
import { fopsApi } from "/@/api/fops";
// 引入 api 请求接口
const serverApi = fopsApi();
const state = reactive({
    isShowDockerLogDialog: false, //容器日志
	dockerLogContent: [],//容器日志
	dockerLog: {
		Id: '',
		Name: '', Node: '', State: '', StateInfo: '', Error: '', Image: '',
	},//容器日志选中
});
const scrollableDockerLog = ref(); //容器日志
//点击容器日志选项
const clickDockerLog = (item) => {
    state.dockerLog = item
}
const openDockerLog = (appName) => { //容器日志
    serverApi.dockerLog({ "AppName": appName, "tailCount": 500 }).then(function (res) {
        state.dockerLogContent = res.Data;
        if (state.dockerLogContent && state.dockerLogContent.length > 0) {
            clickDockerLog(state.dockerLogContent[0])
        }
        state.isShowDockerLogDialog = true
        setTimeout(() => {   //自动跳到底部 
            scrollableDockerLog.value.scrollTop = scrollableDockerLog.value.scrollHeight;
        }, 500)
    })
}
defineExpose({
    openDockerLog
})
</script>
<style scoped>
</style>