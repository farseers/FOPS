<template>
    <div class="conly_header">
        <div class="c1">任务组数量：<span id="c1">0</span></div>
        <div class="c2">超时未运行的任务数量：<span id="c2">0</span></div>
        <div class="c3">今天失败数量：<span id="c3">0</span></div>
    </div>
</template>

<script setup name="Counts" lang="ts">
import { reactive, onMounted,defineExpose,nextTick   } from 'vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from "/@/api/fops";
import { CountUp } from 'countup.js';
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
    TaskGroupCount:0,
    TaskGroupUnRunCount:0,
    TodayFailCount:0
});
// 初始化表格数据
const getData = () => {
  // 任务组数量
  serverApi.statInfo({}).then(function(res){
    if(res.Status){
        if( state.TaskGroupCount != res.Data.TaskGroupCount){
                state.TaskGroupCount = res.Data.TaskGroupCount;
                nextTick(()=>{
                    new CountUp(document.querySelector('#c1') as HTMLDivElement, state.TaskGroupCount).start();
               })
                
            }
            if( state.TaskGroupUnRunCount != res.Data.TaskGroupUnRunCount){
                state.TaskGroupUnRunCount = res.Data.TaskGroupUnRunCount;
                nextTick(()=>{
                    new CountUp(document.querySelector('#c2') as HTMLDivElement, state.TaskGroupUnRunCount).start();
               })
                
            }
            if( state.TodayFailCount != res.Data.TodayFailCount){
                state.TodayFailCount = res.Data.TodayFailCount;
                nextTick(()=>{
                    new CountUp(document.querySelector('#c3') as HTMLDivElement, state.TodayFailCount).start();
               })
                
            }
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
.conly_header{
    display: flex;
    align-items: center;
    height: 100%;
    color: #fff;
    font-size: 14px;
    font-weight: 700;
    overflow: hidden;
    border-radius: 5px;
    div{
        padding: 10px;
        flex: 1
    }
    .c1{
        background: var(--el-color-primary);
    }
    .c2{
        background: var(--el-color-warning);
    }
    .c3{
        background: var(--el-color-danger);
    }
}
</style>