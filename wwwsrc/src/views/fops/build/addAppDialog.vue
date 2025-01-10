<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="900px">
			<el-form ref="gitDialogFormRef" :model="state.ruleForm" size="default" label-width="120px">
				<el-row :gutter="35"><el-form-item label="应用名称">
          <el-input v-model="state.ruleForm.AppName" placeholder="请输入应用名称" style="max-width: 200px;margin-right: 5px"></el-input>
            名称需要与应用的AppName完全一致，才能检查健康状态
        </el-form-item>
        <el-form-item style="float: left" label="副本数量">
          <el-input v-model="state.ruleForm.DockerReplicas" type="number" placeholder="请输入副本数量"></el-input>
        </el-form-item>
          <el-form-item label="容器节点角色">
            <el-select v-model="state.ruleForm.DockerNodeRole" placeholder="请输入容器节点角色" class="ml10" style="max-width: 150px;" size="default">
              <el-option label="manager" value="manager"></el-option>
              <el-option label="worker" value="worker"></el-option>
              <el-option label="global" value="global"></el-option>
              <el-option label="不限制" value=""></el-option>
            </el-select>
          </el-form-item>
        <el-form-item label="容器参数">
          <el-input v-model="state.ruleForm.AdditionalScripts" type="textarea" placeholder="容器在创建时，附加的参数" clearable></el-input>
        </el-form-item>
        <el-form-item style="float: left" label="Cpu限制">
          <el-input v-model="state.ruleForm.LimitCpus" type="number" placeholder="请输入Cpu数量"></el-input>
        </el-form-item>
        <el-form-item label="内存限制">
          <el-input v-model="state.ruleForm.LimitMemory" placeholder="请输入内存"></el-input>
        </el-form-item>
          <el-form-item label="自动构建">
            <el-input v-model="state.ruleForm.UTWorkflowsName" placeholder="请输入工作流的文件名称" clearable></el-input>
          </el-form-item>
        <el-form-item label="Dockerfile">
          <el-input v-model="state.ruleForm.DockerfilePath" placeholder="请输入Dockerfile路径，默认为：./Dockerfile" clearable></el-input>
        </el-form-item>
        <el-form-item label="Git主仓库">
          <el-tag size="small">{{state.ruleForm.AppGitName}}</el-tag>
          <el-button type="primary" @click="onOpenGit(2)" size="default" style="margin-left: 5px;">添加Git</el-button>
        </el-form-item>
        <el-form-item label="依赖仓库">
          <el-button type="success" @click="onOpenGit(1)" size="default" style="margin-left: 5px;">添加依赖的仓库</el-button>
          <el-table :data="state.gitList" style="width: 100%">
            <el-table-column prop="Id" label="编号" width="60" />
            <el-table-column prop="Name" label="Git名称" show-overflow-tooltip width="120"></el-table-column>
            <el-table-column prop="Hub" label="托管地址" show-overflow-tooltip></el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="scope">
                <el-button size="small" text type="primary" @click="delGit(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
          <el-form-item label="工作流文件">
            在git仓库 <b>.fops/workflows/</b> 目录中定义yml后缀文件
          </el-form-item>
				</el-row>
			</el-form>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{ state.dialog.submitTxt }}</el-button>
				</span>
			</template>
		</el-dialog>
    <el-dialog title="Git列表" v-model="state.gitDialogIsShow" width="700px;" height="300px;">
      <el-card shadow="hover" class="layout-padding-auto">
        <div class="system-user-search mb15">
          <el-button size="default" type="success" class="ml10" @click="SureCheck()">
            <el-icon>
            </el-icon>
            确认选择
          </el-button>
        </div>
        <el-table ref="multipleTableRef" :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" :row-key="getRowKey" @selection-change="handleSelectionChange">
          <el-table-column type="selection" :reserve-selection="true" width="55"></el-table-column>
          <el-table-column prop="Id"  label="编号" width="60" />
          <el-table-column prop="Name" label="Git名称" show-overflow-tooltip></el-table-column>
          <el-table-column prop="Hub" label="托管地址" show-overflow-tooltip></el-table-column>
          <el-table-column prop="Branch" label="Git分支" show-overflow-tooltip></el-table-column>
        </el-table>
      </el-card>
    </el-dialog>
	</div>
</template>

<script setup lang="ts" name="fopsAppDialog">
import {reactive, ref, onMounted, getCurrentInstance, nextTick} from 'vue';
import {fopsApi} from "/@/api/fops";
import {ElMessageBox, ElMessage, ElTable} from 'element-plus';
import {friendlyJSONstringify} from "@intlify/shared";
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh','showOverlay','hideOverlay']);
const { proxy } = getCurrentInstance() as any;
// 定义变量内容
const gitDialogFormRef = ref();
const multipleTableRef = ref<InstanceType<typeof ElTable>>();
const state = reactive({
	ruleForm: {
    AppName:'', //应用名称
    DockerVer: '', // 镜像版本
    LocalClusterVer: { // 集群版本
      ClusterId: 0,
      DockerImage: '',
    },
    AppGit: 0, // 应用的源代码
    AppGitName: '', // 应用的源代码
    FrameworkGits:[], // 依赖的框架源代码
    DockerfilePath: '', // Dockerfile路径
    IsHealth:false, // 是否健康
    DockerInstances:0, // 实例数量
    DockerReplicas:1,// 副本数量
    DockerNodeRole:'',// 容器节点角色 manager or worker
    AdditionalScripts:'',// 多行内容，用多行文本框
    WorkflowsYmlPath:'',// 工作流定义的路径,
    UTWorkflowsName:'',// UT工作流名称（文件的名称）
    LimitCpus:0,        // Cpu核数限制
    LimitMemory:'',      // 内存限制
	},
  gitList:[],
  SelectItem:[],
  gitDialogIsShow:false,
	dialog: {
		isShowDialog: false,
		type: '',
		title: '',
		submitTxt: '',
	},
  tableData: {
    data: [],
    total: 0,
    loading: false,
    param: {
      pageNum: 1,
      pageSize: 10,
    },
  },
  gitType:1,
  isApp:-1,
});

// 打开弹窗
const openDialog = (type: string, row: any) => {
  state.dialog.type=type
  state.dialog.title = '新增应用';
  state.dialog.submitTxt = '新 增';
  // 清空表单，此项需加表单验证才能使用
  state.ruleForm.AppName=""
  state.ruleForm.DockerVer=""
  state.ruleForm.ClusterVer=""
  state.ruleForm.AppGit=0
  state.ruleForm.LimitCpus=0
  state.ruleForm.LimitMemory=0
  state.ruleForm.AppGitName=''
  state.ruleForm.FrameworkGits=[]
  state.ruleForm.DockerfilePath=""
  state.ruleForm.DockerReplicas=1
  state.ruleForm.DockerNodeRole=''
  state.ruleForm.AdditionalScripts=''
  state.ruleForm.WorkflowsYmlPath=''
  state.ruleForm.UTWorkflowsName=''
  state.SelectItem=[] // 清空
  state.tableData.data=[]
	state.dialog.isShowDialog = true;
};

const loadGit=(lst:any)=>{
  state.gitList=[]
  for (let i = 0; i < lst.length; i++) {
    serverApi.gitInfo({"gitId":lst[i]}).then(function (res){
      if (res.Status){
        state.gitList.push(res.Data)
      }else{
        state.gitList=[]
      }
    })
  }
}
const loadGitInfo=(id:any)=>{
    serverApi.gitInfo({"gitId":id}).then(function (res){
      if (res.Status){
        state.ruleForm.AppGitName= res.Data.Name
      }else{
        state.ruleForm.AppGitName=""
      }
    })
}
const delGit=(row:any)=>{
  state.ruleForm.FrameworkGits = state.ruleForm.FrameworkGits.filter(number => number !== parseInt(row.Id));
  loadGit(state.ruleForm.FrameworkGits)
}

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
  // 提交数据
  var param={
    "AppName":state.ruleForm.AppName,
    "AppGit":parseInt(state.ruleForm.AppGit),
    "FrameworkGits":state.ruleForm.FrameworkGits,
    "DockerfilePath":state.ruleForm.DockerfilePath,
    "DockerReplicas":parseInt(state.ruleForm.DockerReplicas),
    "DockerNodeRole":state.ruleForm.DockerNodeRole,
    "AdditionalScripts":state.ruleForm.AdditionalScripts,
    "WorkflowsYmlPath":state.ruleForm.WorkflowsYmlPath,
    "UTWorkflowsName":state.ruleForm.UTWorkflowsName,
    "LimitCpus":parseFloat(state.ruleForm.LimitCpus),
    "LimitMemory":state.ruleForm.LimitMemory,
  }
  emit('showOverlay');
	if (state.dialog.type === 'add') {
    var json=JSON.stringify(param)
    serverApi.appsAdd(json).then(function (res){
      if(res.Status){
        ElMessage.success("添加成功")
        closeDialog();
        emit('refresh');
      }else{
        ElMessage.error(res.StatusMessage)
      }
      emit('hideOverlay');
    })

  }else if (state.dialog.type=='edit'){
    serverApi.appsEdit(param).then(function (res){
      if(res.Status){
        ElMessage.success("修改成功")
        closeDialog();
        emit('refresh');
      }else{
        ElMessage.error(res.StatusMessage)
      }
      emit('hideOverlay');
    })

  }
};

const getTableData = (type:any) => {
  if (type==1) {
    state.isApp=0
    state.SelectItem=state.ruleForm.FrameworkGits // 清空
  }else{
    state.isApp=1
    var select=[]
    select.push(state.ruleForm.AppGit)
    state.SelectItem=select // 清空
  }
  // 请求接口
  serverApi.gitList({isApp:state.isApp}).then(function (res){
    if (res.Status){
      state.tableData.data = res.Data;
      state.tableData.total = res.Data.length;
      onloadSelect(type)
    }else{
      state.tableData.data=[]
    }
  })
};

const onloadSelect=(type:number)=>{
  // 清空选项
  state.tableData.data.forEach(function (item,index){
    setCurrent(item,false)
  })
  if (type==1){
    state.tableData.data.forEach(function (item,index){
      var rowArray=state.ruleForm.FrameworkGits.filter(t=>t==item.Id);
      if(rowArray.length>0)
      {
        setCurrent(item,true)
      }else{
        setCurrent(item,false)
      }
    })
  }else{
    state.tableData.data.forEach(function (item,index){
      if(state.ruleForm.AppGit==item.Id)
      {
        setCurrent(item,true)
      }else{
        setCurrent(item,false)
      }
    })
  }
}
const getRowKey=(row:any)=>{
  return row.Id;
}
const handleSelectionChange=(val:any)=> {
  console.log(val)
  if(val.length==0){return;}
  state.SelectItem=[] // 清空
  for (let i = 0; i < val.length; i++) {
    var item=val[i]
    if(item.IsApp&&state.isApp==1){
      state.SelectItem.push(item.Id)
    }
    if(!item.IsApp&&state.isApp==0){
      state.SelectItem.push(item.Id)
    }
  }
  console.log(state.SelectItem)
}
const onOpenGit=(type:any)=>{
  state.gitType=type
  state.gitDialogIsShow=true
  getTableData(type)
}

const setCurrent=(row:any,isSelect:boolean)=>{
  nextTick(()=>{
    proxy.$refs.multipleTableRef.toggleRowSelection(row,isSelect)
  })
}
// 确认选择
const SureCheck=()=>{
  if (state.gitType==1){
    state.ruleForm.FrameworkGits=state.SelectItem
    loadGit(state.ruleForm.FrameworkGits)
  }else{
    state.ruleForm.AppGit=state.SelectItem[0]
    loadGitInfo(state.ruleForm.AppGit)
  }
  state.gitDialogIsShow=false
}
// 页面加载时
onMounted(() => {

});
// 暴露变量
defineExpose({
	openDialog,
  delGit,
});
</script>

<style>
textarea{
  height: 220px;
}
</style>
