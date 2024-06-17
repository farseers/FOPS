<template>
    <div v-loading="state.loading" class="w100">
        <div class="conlyRow">
            <div v-for="item, index in state.tableData" :key="index.toString() + 'ecfy'" class="conlyCol">
                <el-card :class="item.IsHealth ? 'conlyCard' : 'conlyCard conly_w'">
                    <div class="name">
                        <el-tag size="default">{{ item.AppName }}</el-tag>
                        <el-tooltip content="实例数量/副本数量" slot="label">
                            <el-tag v-if="item.IsHealth" size="small" style="margin-left: 5px">{{ item.DockerInstances
                                }}/{{ item.DockerReplicas }}</el-tag>
                            <el-tag v-else size="small" type="danger" style="margin-left: 5px">{{ item.DockerInstances
                                }}/{{ item.DockerReplicas }}</el-tag>
                        </el-tooltip>
                    </div>
                    <div>
                        <el-button class="ecfy_btn" size="small" type="primary" @click="showDockerLog(item.AppName)">容器日志</el-button>
                            <el-button class="ecfy_btn" size="small" type="success" @click="showFsLogLevel(2, item.AppName)">应用日志</el-button>
                    </div>
                    <div>应用日志
                        <el-tooltip content="警告数量" slot="label">
                            <el-tag v-if="item.LogWaringCount > 0" @click="showFsLogLevel(3, item.AppName)"
                                type="warning" size="small" style="margin-left: 5px;cursor: pointer">{{
        item.LogWaringCount }}</el-tag>
                            <el-tag v-else @click="showFsLogLevel(3, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.LogWaringCount }}</el-tag>
                        </el-tooltip>
                        /
                        <el-tooltip content="异常数量" slot="label">
                            <el-tag v-if="item.LogErrorCount > 0" @click="showFsLogLevel(4, item.AppName)" type="danger"
                                size="small" style="margin-left: 5px;cursor: pointer">{{ item.LogErrorCount }}</el-tag>
                            <el-tag v-else @click="showFsLogLevel(4, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.LogErrorCount }}</el-tag>
                        </el-tooltip>
                    </div>
                    <div>调度任务
                        <el-tooltip content="成功数量" slot="label">
                            <el-tag v-if="item.TaskSuccessCount > 0" @click="showTask(2, item.AppName)" type="success"
                                size="small" style="margin-left: 5px;cursor: pointer">{{ item.TaskSuccessCount
                                }}</el-tag>
                            <el-tag v-else @click="showTask(2, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.TaskSuccessCount }}</el-tag>
                        </el-tooltip>
                        /
                        <el-tooltip content="失败数量" slot="label">
                            <el-tag v-if="item.TaskFailCount > 0" @click="showTask(3, item.AppName)" type="danger"
                                size="small" style="margin-left: 5px;cursor: pointer">{{ item.TaskFailCount }}</el-tag>
                            <el-tag v-else @click="showTask(3, item.AppName)" type="info" size="small"
                                style="margin-left: 5px;cursor: pointer">{{ item.TaskFailCount }}</el-tag>
                        </el-tooltip>
                    </div>
                </el-card>
            </div>
        </div>
            <dockerDialog ref="dockerDialogRef"/>
            <taskDialog ref="taskDialogRef"  />
            <logDialog ref="logDialogRef"  />
    </div>
</template>
<script setup name="Ecfy">
import { reactive, onMounted, defineExpose, ref,defineAsyncComponent } from 'vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from "/@/api/fops";
const dockerDialog = defineAsyncComponent(() => import('/src/views/fops/task/dockerDialog.vue'));
const logDialog = defineAsyncComponent(() => import('/src/views/fops/log/logV2Dialog.vue'));
const taskDialog= defineAsyncComponent(() => import('/src/views/fops/task/taskAppDialog.vue'));
// 引入 api 请求接口
const serverApi = fopsApi();
const conlyTabs = ref(null)
const logDialogRef = ref();
const taskDialogRef = ref();
const dockerDialogRef = ref();
// 定义变量内容
const state = reactive({
    tableData: [],
    loading: false,
    statTask: [],
    isShowDockerLogDialog: false, //容器日志
	dockerLogContent: [],//容器日志
	dockerLog: {
		Id: '',
		Name: '', Node: '', State: '', StateInfo: '', Error: '', Image: '',
	},//容器日志选中
});

const showDockerLog = (AppName) => {
    dockerDialogRef.value.openDockerLog(AppName);
}
// 打开FS日志
const showFsLogLevel=(level,appName)=>{
  logDialogRef.value.openDialogLogLevel(level,appName,'应用日志');
}
// 任务组日志
const showTask=(st,appName)=>{
  taskDialogRef.value.openDialogApp(st,appName);
}

const getData = () => {
    var param = {
        "ClusterId": 0,
    }
    state.loading = true
    // 获取应用列表
    serverApi.appsList(param).then(function (res) {
        if (res.Status) {
            state.tableData = res.Data;
        } else {
            ElMessage.warning(res.StatusMessage);
        }
        state.loading = false
    }).catch(() => {
        state.loading = false
    })
}
onMounted(() => {
    getData()
});
defineExpose({
    getData,
});
</script>
<style scoped lang="scss">
.el-dialog__body {
		display: flex;
		flex-direction: column;
	}

	.layout-container .layout-padding-auto {
		flex: 1;
		overflow: auto;
	}

.conlyRow {
    flex-wrap: wrap;
    display: flex !important;
    min-height: 200px;
}

.conlyCol {
     padding: 5px;
    box-sizing: border-box;
    width: 180px;
}


.conlyCard {
    background-color: #f9f9e3;
    border: 1px dotted var(--el-color-primary);
    :deep(.el-card__body) {
        padding: 10px 5px;
	}
    .el-card__body>div {
    text-align: left;
    margin-bottom: 5px;
    font-size: 12px;
}
    .name {
        font-weight: 700;
    }
}

.conly_w {
    background: var(--el-color-danger-light-9);
    border: 1px dotted var(--el-color-danger);
}

.layout-container .layout-padding-auto {
    flex: 1;
    overflow: auto;
}
.ecfy_btn{
    --el-button-size: 20px;
    padding: 3px 9px;
}
</style>