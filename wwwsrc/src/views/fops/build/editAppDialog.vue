<template>
  <div class="system-user-dialog-container">
    <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="1100px">
      <el-form ref="gitDialogFormRef" :model="state.ruleForm" size="default" label-width="120px">
        <el-row :gutter="35">
          <el-form-item label="应用名称">
            <el-input v-model="state.ruleForm.AppName" placeholder="请输入应用名称"
              style="max-width: 200px;margin-right: 5px"></el-input>
            <el-button v-if="state.dialog.type == 'edit'" @click="onDeleteApp" size="default" type="danger"
              style="margin-left: 5px;">删除应用</el-button>
          </el-form-item>
          <el-form-item label="仓库版本">
            <el-tag v-if="state.ruleForm.DockerImage != ''" size="small">{{ state.ruleForm.DockerImage }}</el-tag>
            <el-tag v-else size="small">未构建</el-tag>
          </el-form-item>
          <el-form-item label="集群版本">
            <el-input v-model="state.ruleForm.LocalClusterVer.DockerImage" placeholder="镜像名称"
              style="max-width: 300px;margin-right: 5px"></el-input>
            部署时间：
            <el-tag v-if="state.ruleForm.LocalClusterVer.DockerImage != ''" size="small"
              style="margin-right: 5px;">{{ state.ruleForm.LocalClusterVer.DeploySuccessAt }}</el-tag>
            <el-tag v-else size="small" style="margin-right: 5px;">未发布</el-tag>

            <el-tag v-if="state.ruleForm.IsHealth" size="small" type="success">健康</el-tag>
            <el-tag v-else-if="state.ruleForm.DockerInstances > 0" size="small" type="warning">不健康</el-tag>
            <el-tag v-else size="small" type="danger">未运行</el-tag>
          </el-form-item>
          <el-form-item style="float: left" label="副本数量">
            <el-input v-model="state.ruleForm.DockerReplicas" type="number" placeholder="请输入副本数量"></el-input>
          </el-form-item>
          <el-form-item label="容器节点角色">
            <el-select v-model="state.ruleForm.DockerNodeRole" placeholder="请输入容器节点角色" class="ml10"
              style="max-width: 150px;" size="default">
              <el-option label="manager" value="manager"></el-option>
              <el-option label="worker" value="worker"></el-option>
              <el-option label="global" value="global"></el-option>
              <el-option label="不限制" value=""></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="容器参数">
            <el-input v-model="state.ruleForm.AdditionalScripts" type="textarea" placeholder="容器在创建时，附加的参数"
              clearable></el-input>
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
            <el-input v-model="state.ruleForm.DockerfilePath" placeholder="请输入Dockerfile路径，默认为：./Dockerfile"
              clearable></el-input>
          </el-form-item>
          <el-form-item label="Git主仓库">
            <el-tag size="small">{{ state.ruleForm.AppGitName }}</el-tag>
            <el-button type="primary" @click="onOpenGit(2)" size="default" style="margin-left: 5px;">添加Git</el-button>
          </el-form-item>
          <el-form-item label="依赖仓库">
            <el-button type="success" @click="onOpenGit(1)" size="default" style="margin-left: 5px;">添加依赖的仓库</el-button>
            <el-table :data="state.gitList" style="width: 100%">
              <el-table-column prop="Id" label="编号" width="60" />
              <el-table-column prop="Name" label="Git名称" show-overflow-tooltip width="80"></el-table-column>
              <el-table-column prop="Hub" label="托管地址" show-overflow-tooltip></el-table-column>
              <el-table-column prop="CommitId" label="CommitId" width="200">
                <template #default="scope">
                  <el-input v-model="scope.row.CommitId" size="small" clearable placeholder="请输入CommitId" @change="onCommitIdChange(scope.row)"></el-input>
                  <!-- <span>{{ scope.row.CommitId || '未构建' }}</span> -->
                </template>
              </el-table-column>
              <!-- <el-table-column label="自动更新" width="60">
                <template #default="scope">
                  <el-switch v-model="scope.row.IsAutoUpdate" @change="onAutoUpdateChange(scope.row)" />
                </template>
              </el-table-column> -->
              <el-table-column label="操作" width="60">
                <template #default="scope">
                  <el-button size="small" text type="primary" @click="delGit(scope.row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
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

    <el-dialog title="框架列表" v-model="state.gitDialogIsShow" width="700px;" height="300px;">
      <el-card shadow="hover" class="layout-padding-auto">
        <div class="system-user-search mb15">
          <el-button size="default" type="success" class="ml10" @click="SureCheck()">
            <el-icon>
            </el-icon>
            确认选择
          </el-button>
        </div>
        <el-table ref="multipleTableRef" :data="state.tableData.data" v-loading="state.tableData.loading"
          style="width: 100%" :row-key="getRowKey" @selection-change="handleSelectionChange">
          <el-table-column type="selection" :reserve-selection="true" width="55"></el-table-column>
          <el-table-column prop="Id" label="编号" width="60" />
          <el-table-column prop="Name" label="Git名称" show-overflow-tooltip></el-table-column>
          <el-table-column prop="Hub" label="托管地址" show-overflow-tooltip></el-table-column>
          <el-table-column prop="Branch" label="Git分支" show-overflow-tooltip></el-table-column>
        </el-table>
      </el-card>
    </el-dialog>

  </div>
</template>

<script setup lang="ts" name="fopsAppDialog">
import { reactive, ref, onMounted, getCurrentInstance, nextTick } from 'vue';
import { fopsApi } from "/@/api/fops";
import { ElMessageBox, ElMessage, ElTable } from 'element-plus';
import { friendlyJSONstringify } from "@intlify/shared";
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh', 'showOverlay', 'hideOverlay']);
const { proxy } = getCurrentInstance() as any;
// 定义变量内容
const gitDialogFormRef = ref();
const multipleTableRef = ref<InstanceType<typeof ElTable>>();

// 定义 FrameworkItem 类型
interface FrameworkItem {
  FrameworkId: number;
  CommitId: string;
  //IsAutoUpdate: boolean;
}

const state = reactive({
  ruleForm: {
    ClusterId: 0, // 集群ID
    AppName: '', //应用名称
    DockerVer: '', // 镜像版本
    DockerImage: '', // Docker镜像
    LocalClusterVer: { // 集群版本
      ClusterId: 0,
      DockerImage: '',
    },
    AppGit: 0, // 应用的源代码
    AppGitName: '', // 应用的源代码
    FrameworkList: [] as FrameworkItem[], // 依赖的框架源代码
    DockerfilePath: '', // Dockerfile路径
    IsHealth: false, // 是否健康
    DockerInstances: 0, // 实例数量
    DockerReplicas: 1,// 副本数量
    DockerNodeRole: '',// 容器节点角色 manager or worker
    AdditionalScripts: '',// 多行内容，用多行文本框
    WorkflowsYmlPath: '',// 工作流定义的路径
    UTWorkflowsName: '',// UT工作流名称（文件的名称）
    LimitCpus: 0,        // Cpu核数限制
    LimitMemory: '',      // 内存限制

  },
  gitList: [],
  SelectItem: [],
  gitDialogIsShow: false,
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
  gitType: 1,
  isApp: -1,
});

// 打开弹窗
const openDialog = (type: string, row: any, clusterId: number) => {
  state.dialog.type = type
  state.ruleForm = row;
  state.ruleForm.ClusterId = clusterId;
  state.dialog.title = '修改应用';
  state.dialog.submitTxt = '修 改';
  // 请求数据
  serverApi.appsDetail({ AppName: row.AppName, clusterId: clusterId }).then(function (res) {
    if (res.Status) {
      row = res.Data
      // 绑定数据
      state.ruleForm.AppName = row.AppName
      state.ruleForm.DockerVer = row.DockerVer
      state.ruleForm.DockerImage = row.DockerImage || ''
      state.ruleForm.LocalClusterVer = row.LocalClusterVer
      state.ruleForm.AppGit = row.AppGit
      state.ruleForm.FrameworkList = row.FrameworkList || []
      state.ruleForm.DockerfilePath = row.DockerfilePath
      // 从 FrameworkList 中提取 ID 列表用于选择
      state.SelectItem = row.FrameworkList ? row.FrameworkList.map((item: any) => item.FrameworkId) : []
      state.ruleForm.IsHealth = row.IsHealth
      state.ruleForm.DockerReplicas = row.DockerReplicas
      state.ruleForm.DockerNodeRole = row.DockerNodeRole
      state.ruleForm.AdditionalScripts = row.AdditionalScripts
      state.ruleForm.WorkflowsYmlPath = row.WorkflowsYmlPath
      state.ruleForm.UTWorkflowsName = row.UTWorkflowsName
      state.ruleForm.LimitCpus = row.LimitCpus
      state.ruleForm.LimitMemory = row.LimitMemory
      // if (state.ruleForm.DockerNodeRole == "manager") {
      //   state.ruleForm.DockerNodeRoleInt=0
      // } else {
      //   state.ruleForm.DockerNodeRoleInt=1
      // }
      state.ruleForm.AppGitName = row.AppGitName
      //loadGitInfo(row.AppGit)
      // 加载git数据
      loadGit()
    }
  })
  state.dialog.isShowDialog = true;
};

const loadGit = () => {
  serverApi.gitList({ isApp: 0 }).then(function (res) {
    if (res.Status) {
      // 从 FrameworkList 中提取 ID 列表
      const SelectItem = state.ruleForm.FrameworkList.map((item: any) => item.FrameworkId);
      const arr = res.Data.filter((item: any) => SelectItem.includes(item.Id));
      // 将 git 信息与 FrameworkList 合并，添加 Name 等信息
      state.gitList = arr.map((gitItem: any) => {
        const framework = state.ruleForm.FrameworkList.find((f: any) => f.FrameworkId === gitItem.Id);
        return {
          ...gitItem,
          CommitId: framework?.CommitId || '',
          //IsAutoUpdate: framework?.IsAutoUpdate !== undefined ? framework.IsAutoUpdate : true
        };
      });
    } else {
      state.gitList = []
    }
  })
}
const loadGitInfo = (id: any) => {
  serverApi.gitInfo({ "gitId": id }).then(function (res) {
    if (res.Status) {
      state.ruleForm.AppGitName = res.Data.Name
    } else {
      state.ruleForm.AppGitName = ""
    }
  })
}
const delGit = (row: any) => {
  state.ruleForm.FrameworkList = state.ruleForm.FrameworkList.filter((item: any) => item.FrameworkId !== parseInt(row.Id));
  loadGit()
}

// 处理CommitId变化
const onCommitIdChange = (row: any) => {
  // 更新 FrameworkList 中对应项的 IsAutoUpdate 值
  const framework = state.ruleForm.FrameworkList.find((item: any) => item.FrameworkId === row.Id);
  if (framework) {
    framework.CommitId = row.CommitId;
  }
}

// // 处理自动更新开关变化
// const onAutoUpdateChange = (row: any) => {
//   // 更新 FrameworkList 中对应项的 IsAutoUpdate 值
//   const framework = state.ruleForm.FrameworkList.find((item: any) => item.FrameworkId === row.Id);
//   if (framework) {
//     framework.IsAutoUpdate = row.IsAutoUpdate;
//   }
// }

// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
// 取消
const onCancel = () => {
  closeDialog();
};

// 删除应用
const onDeleteApp = () => {
  ElMessageBox.confirm(`此操作将永久删除应用：“${state.ruleForm.AppName}”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      // 删除逻辑
      serverApi.appsDel({ "appName": state.ruleForm.AppName }).then(function (res) {
        if (res.Status) {
          closeDialog();
          ElMessage.success('删除成功');
          emit('refresh');
        } else {
          ElMessage.error(res.StatusMessage)
        }
      })
    })
    .catch(() => { });
};

// 提交
const onSubmit = () => {
  // 提交数据
  var param = {
    "ClusterId": state.ruleForm.LocalClusterVer.ClusterId,
    "ClusterDockerImage": state.ruleForm.LocalClusterVer.DockerImage,
    "AppName": state.ruleForm.AppName,
    "AppGit": parseInt(state.ruleForm.AppGit),
    "FrameworkList": state.ruleForm.FrameworkList.map((item: any) => ({
      FrameworkId: item.FrameworkId,
      CommitId: item.CommitId || '',
      //IsAutoUpdate: item.IsAutoUpdate !== undefined ? item.IsAutoUpdate : true
    })),
    "DockerfilePath": state.ruleForm.DockerfilePath,
    "DockerReplicas": parseInt(state.ruleForm.DockerReplicas),
    "AdditionalScripts": state.ruleForm.AdditionalScripts,
    "WorkflowsYmlPath": state.ruleForm.WorkflowsYmlPath,
    "UTWorkflowsName": state.ruleForm.UTWorkflowsName,
    "LimitCpus": parseFloat(state.ruleForm.LimitCpus),
    "LimitMemory": state.ruleForm.LimitMemory,
    "DockerNodeRole": state.ruleForm.DockerNodeRole,
  }
  emit('showOverlay');
  serverApi.appsEdit(param).then(function (res) {
    if (res.Status) {
      ElMessage.success("修改成功")
      closeDialog();
      emit('refresh');
    } else {
      ElMessage.error(res.StatusMessage)
    }
    emit('hideOverlay');
  })
};

const getTableData = (type: any) => {
  if (type == 1) {
    state.isApp = 0
    // 从 FrameworkList 中提取 ID 列表
    state.SelectItem = state.ruleForm.FrameworkList.map((item: any) => item.FrameworkId)
  } else {
    state.isApp = 1
    var select = []
    select.push(state.ruleForm.AppGit)
    state.SelectItem = select // 清空
  }
  // 请求接口
  serverApi.gitList({ isApp: state.isApp }).then(function (res) {
    console.log(11111111)
    if (res.Status) {
      state.tableData.data = res.Data;
      state.tableData.total = res.Data.length;
      onloadSelect(type)
    } else {
      state.tableData.data = []
    }
  })
};

const onloadSelect = (type: number) => {
  // 清空选项
  state.tableData.data.forEach(function (item, index) {
    setCurrent(item, false)
  })
  if (type == 1) {
    // 从 FrameworkList 中提取 ID 列表
    const frameworkIds = state.ruleForm.FrameworkList.map((item: any) => item.FrameworkId);
    state.tableData.data.forEach(function (item, index) {
      var rowArray = frameworkIds.filter((id: any) => id == item.Id);
      if (rowArray.length > 0) {
        setCurrent(item, true)
      } else {
        setCurrent(item, false)
      }
    })
  } else {
    state.tableData.data.forEach(function (item, index) {
      if (state.ruleForm.AppGit == item.Id) {
        setCurrent(item, true)
      } else {
        setCurrent(item, false)
      }
    })
  }
}
const getRowKey = (row: any) => {
  return row.Id;
}
const handleSelectionChange = (val: any) => {
  console.log(val)
  if (val.length == 0) { return; }
  state.SelectItem = [] // 清空
  for (let i = 0; i < val.length; i++) {
    var item = val[i]
    if (item.IsApp && state.isApp == 1) {
      state.SelectItem.push(item.Id)
    }
    if (!item.IsApp && state.isApp == 0) {
      state.SelectItem.push(item.Id)
    }
  }
  console.log(state.SelectItem)
}
const onOpenGit = (type: any) => {
  state.gitType = type
  state.gitDialogIsShow = true
  getTableData(type)
}

const setCurrent = (row: any, isSelect: boolean) => {
  nextTick(() => {
    proxy.$refs.multipleTableRef.toggleRowSelection(row, isSelect)
  })
}
// 确认选择
const SureCheck = () => {
  if (state.gitType == 1) {
    // 更新 FrameworkList：保留已存在的项并更新，添加新选择的项
    const existingMap = new Map(state.ruleForm.FrameworkList.map((item: any) => [item.FrameworkId, item]));
    state.ruleForm.FrameworkList = state.SelectItem.map((id: any) => {
      const existing = existingMap.get(id);
      return existing || {
        FrameworkId: id,
        CommitId: '',
        //IsAutoUpdate: true
      };
    });
    loadGit()
  } else {
    state.ruleForm.AppGit = state.SelectItem[0]
    loadGitInfo(state.ruleForm.AppGit)
  }
  state.gitDialogIsShow = false
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
textarea {
  height: 220px;
}
</style>
