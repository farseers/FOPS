<template>
    <div class="w100">
        <div class="conlyRow">
            <div v-for="item, index in state.tableData" :key="index.toString() + 'conly1'" class="conlyCol">
                <el-card :class="item.IsHealth ? 'conlyCard' : 'conlyCard conly_w'">
                    <div class="name">
                        <span>{{ item.IP }} | {{ item.NodeName }} <img v-show="item.OS == 'linux'" :src="linux" alt="" /></span>
                    </div>
                    <div v-show="item.IsMaster">
                        <el-tag type="danger" size="small">manager</el-tag> {{ item.Memory }}
                    </div>
                    <div v-show="!item.IsMaster">
                      <el-tag size="small">worker</el-tag> {{ item.Memory }}
                    </div>
                    <div>{{ item.OS }} {{ item.Architecture }} | {{ item.EngineVersion }}</div>
                    <div>
                        <el-tag effect="dark" size="small" style="margin-right: 5px;"
                            :type="item.Status == 'Ready' ? 'success' : 'danger'">{{ item.Status }}</el-tag>
                        <el-tag effect="dark" size="small" 
                        :type="item.Availability == 'Active' ? 'success' : 'danger'">
                        {{ item.Availability }}
                    </el-tag>
                    </div>
                    <div>CPU ：<b>{{ item.CpuUsagePercent }}</b>% ({{ item.CPUs }}核)</div>
                    <div>内存：<b>{{ item.MemoryUsagePercent }}</b>% / <b>{{ item.MemoryUsage }}</b> MB</div>
                  <!-- <div class="line" v-show="item.Label && item.Label.length>0"></div> -->
                    <el-tag class="ks" v-for="row, j in item.Label" :key="index.toString() + j.toString() + 'conly2'">
                        <div>{{ row.Name }} = {{ row.Value }}</div>
                    </el-tag>
                </el-card>
            </div>
        </div>
    </div>
</template>

<script setup name="fopsTaskTimeOut">
import { reactive, onMounted,defineExpose   } from 'vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from "/@/api/fops";
import linux from '/@/assets/linux.png';
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
    tableData: []
});
// 初始化表格数据
const getData = () => {
    var param = {
    }
    // 请求接口
    serverApi.ColonyNodeList(param).then(function (res) {
        if (res.Status) {
            state.tableData = res.Data;
        } else {
            ElMessage.warning(res.StatusMessage);
        }
    })
};
onMounted(() => {
	getData()
});
defineExpose({
    getData,
});
</script>

<style scoped lang="scss">
.conlyRow {
    flex-wrap: wrap;
    display: flex !important;
    min-height: 160px;
    line-height: 22px;
}

.conlyCol {
     padding: 5px;
    box-sizing: border-box;
    width: 180px;
}

.conlyCard {
    background: var(--el-color-success-light-9);
    //border: 1px dotted var(--el-color-primary);
    :deep(.el-card__body) {
        padding: 10px 5px;
	}
    :deep(.el-tag--large){
        padding: 0 10px;
        height: 26px;
        --el-icon-size: 16px;
    }
    .el-card__body>div {
    text-align: center;
    font-size: 12px
}
    .name {
        font-weight: 700;
        span {
            display: inline-block;
            position: relative;
        }

        img {
            position: absolute;
            top: -5px;
            right: -30px;
            width: 20px;
            margin-left: 10px
        }
    }

    .line {
        width: 100%;
        height: 1px;
        background-color: var(--el-color-info);
    }

    .ks {
        text-align: center;
        border-radius: 5px;
        width: 100%;
        margin-bottom: 5px;

        div {
            padding: 3px;
        }
    }
}

.conly_w {
    background: var(--el-color-danger-light-8);
    border: 1px dotted var(--el-color-danger);
}
</style>
