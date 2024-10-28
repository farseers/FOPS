<template>
    <div class="conlyTab" ref="conlyTabs">
        <h3 style="padding: 5px;">执行计划</h3>
        <div v-if="state.tableData.length > 0" class="dsv">
            <el-table :data="state.tableData" v-loading="state.loading" 
            :height = "state.screenHeight"
            size="small"
            style="width: 100%;">
                <el-table-column label="名称" >
                    <template #default="scope">
                        <div style="float: left">
                            <span>{{ scope.row.Caption }}</span><br>
                            <span>{{ scope.row.Name }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="Plan" label="计划时间" width="110" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag size="small" v-if="scope.row.Plan.includes(`等待调度`)">{{ scope.row.Plan }}</el-tag>
                        <el-tag size="small" v-else-if="scope.row.Plan.includes(`等待`)"
                            type="info">{{ scope.row.Plan }}</el-tag>
                        <el-tag size="small" v-else-if="scope.row.Plan.includes(`超时`)"
                            type="danger">{{ scope.row.Plan }}</el-tag>
                        <el-tag size="small" v-else-if="scope.row.Plan.includes(`已执行`)"
                            type="success">{{ scope.row.Plan }}</el-tag>

                    </template>
                </el-table-column>

            </el-table>
        </div>
        <el-empty v-else description="暂无数据"></el-empty>
    </div>
</template>
<script setup name="QueTab">
import { reactive, onMounted,defineExpose,ref} from 'vue';
import {  ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";

// 引入 api 请求接口
const serverApi = fopsApi();
const conlyTabs = ref(null)
// 定义变量内容
const state = reactive({
	tableData:[],
    loading:false,
    screenHeight:400
});

// 初始化表格数据
const getData = () => {
	const params = new URLSearchParams();
  params.append('top', '20');
  // 请求接口
  serverApi.taskPlanList(params.toString()).then(function (res){
    if (res.Status){
      state.tableData = res.Data;
    }else{
        ElMessage.warning(res.StatusMessage);
       state.tableData=[]
    }
    if(conlyTabs && conlyTabs.value){
        // console.log(conlyTabs.value.offsetHeight)
        state.screenHeight = conlyTabs.value.offsetHeight - 40
    }
  }).catch(()=>{
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
.conlyTab{
    height: 100%;
    display: flex;
    flex-flow: column;
    .dsv{
        flex: 1;
    }
}
</style>