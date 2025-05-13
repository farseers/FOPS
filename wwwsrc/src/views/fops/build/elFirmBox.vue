<template>
	<div>
		<el-dialog title="提示" v-model="state.isShowDialog" width="400px">
			<div class="cropper-warp">
				<div style="margin-bottom: 5px;">请填写分支名称，并确认构建到本地!</div>
                 <el-autocomplete
                    class="inline-input"
                   style="width: 100%;"
                    clearable
                    v-model="state.inputValue"
                    :fetch-suggestions="querySearch"
                    placeholder="请输入分支名称"
                    :trigger-on-focus="true"
                    @select="handleSelect"
                   @keyup.enter.native="onSubmit"
                    ></el-autocomplete>
			</div>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">确 认</el-button>
				</span>
			</template>
		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="cropper">
import { reactive } from 'vue';
import { ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
const emit = defineEmits(['refresh']);
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
	isShowDialog: false,
     restaurants: [],
     inputValue:'',
     workflowsName:'',
     appName:''
});
const handleSelect = (item: any) => {
  console.log(item)
}
const createFilter = (queryString: string) => {
       return (restaurant:any) => {
          return (restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0);
        };
}
const querySearch = (queryString: string, cb: any) => {
  var restaurants = state.restaurants;
  console.log(restaurants)
        var results = queryString ? restaurants.filter(createFilter(queryString)) : restaurants;
        // 调用 callback 返回建议列表的数据
        cb(results);
}

// 打开弹窗
const openDialog = (row: any,workflowsName: any) => {
    state.inputValue = '';
	state.restaurants = [];
	state.isShowDialog = true;
    state.workflowsName = workflowsName;
    state.appName = row.AppName;
     let param = {
    "appName":  state.appName
  };
 
    serverApi.autobuildBranchList(param).then(function (res){
        if (res.Status) {
            const arr = res.Data.map((item : any)=>{
                return {value:item.BranchName}
            });
            state.restaurants = arr
        }
    })
};
// 关闭弹窗
const closeDialog = () => {
	state.isShowDialog = false;
};
// 取消
const onCancel = () => {
	closeDialog();
};
// 更换
const onSubmit = () => {
    if(state.inputValue){
        var param={
          "AppName" : state.appName,
          "WorkflowsName" : state.workflowsName,
          "branchName":state.inputValue
        }
        // console.log(param)
        serverApi.buildAdd(param).then(async function(res){
          if(res.Status){
            ElMessage.success("添加成功")
            // 刷新构建日志
             emit('refresh');
             closeDialog();
          }else{
            ElMessage.error(res.StatusMessage)
          }
        })
    }else{
         ElMessage.error('请输入分支名称')
    }
	   
};


// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss">
.inline-input{
    width: 100%;
}
</style>
