<template>
  <div>
    <el-dialog :title="'构建 ' + state.appName + '('+state.workflowsName+')'" v-model="state.isShowDialog" width="600px" :close-on-click-modal="true">
      <div class="dialog-content">
        <!-- 分支选择区域 -->
        <div class="section-card branch-section">
          <div class="section-header">
            <i class="el-icon-branch"></i>
            <span class="section-title">应用配置</span>
            <el-select placeholder="历史构建清单" v-model="state.manifestSelect" class="ml10" style="width: 450px;" size="default" @change="manifestSelectChange" clearable>
              <el-option v-for="item in state.buildManifestList" :key="item.GitCommitId"
                :label="' 分支: ' + item.GitBranch + '        ' + item.DockerImage + '       ' + item.CreateAt"
                :value="item">
                <!-- 自定义下拉选项的样式 -->
                <span style="float: left">分支: {{ item.GitBranch }}</span>
                <span style="float: right; color: var(--el-text-color-secondary); font-size: 13px;padding-left: 30px;">{{ item.CreateAt }}</span>
                <span style="float: right; color: var(--el-text-color-secondary); font-size: 13px;">{{ item.DockerImage }}</span>
              </el-option>
            </el-select>
          </div>
          <div class="section-body">
            <el-radio-group v-model="state.appBranchName" class="branch-radio-group">
              <el-radio v-for="item in state.appBranchList" :key="item.value" :label="item.value" size="default" class="branch-radio-item" @click="appBranchNameChange">
                {{ item.value }}
              </el-radio>
            </el-radio-group>
          </div>
        </div>

        <!-- 框架配置区域 -->
        <div class="section-card framework-section">
          <div class="section-header">
            <i class="el-icon-setting"></i>
            <span class="section-title">框架配置</span>
            <el-checkbox v-model="state.enableBackDefaultBranch" class="framework-checkbox">
              <span class="checkbox-label">匹配失败时退回到默认分支</span>
            </el-checkbox>
          </div>
          <div class="section-body">
            <el-table :data="state.appFrameworkList" style="width: 100%;"
              :cell-style="{ backgroundColor: '#fef5f5', padding: '5px 20px' }"
              :header-cell-style="{ backgroundColor: '#f5f7fa', padding: '5px 20px', height: '30px' }">
              <!-- <el-table-column prop="Id" label="编号" width="80"/> -->
              <el-table-column prop="Name" label="框架" width="150"></el-table-column>
              <!-- <el-table-column prop="Hub" label="托管地址" show-overflow-tooltip></el-table-column> -->
              <el-table-column prop="Branch" label="分支">
                <template #default="scope">
                  <!-- <span>{{ scope.row.IsAutoUpdate ? scope.row.Branch : scope.row.Branch}}</span> -->
                  <!-- <el-input v-model="scope.row.Branch" clearable style="height: 30px;"></el-input> -->
                  <el-autocomplete v-if="state.isShowDialog" class="inline-input" style="width: 100%;"
                    ref="autoCompleteRef" clearable v-model="scope.row.Branch" :debounce="0"
                    :fetch-suggestions="frameworkSearch" @focus="handleFrameworkFocus(scope.row)" placeholder="请输入分支名称"
                    :trigger-on-focus="true" @keyup.enter.native="onSubmit"></el-autocomplete>
                </template>
              </el-table-column>
            </el-table>
            <!-- <div class="checkbox-hint">勾选后将自动更新到最新版本</div> -->
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="onCancel" size="default">取 消</el-button>
          <el-button type="primary" @click="onSubmit" size="default">确 认</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="cropper">
import { ref, reactive } from 'vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from "/@/api/fops";
const emit = defineEmits(['refresh']);
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
  isShowDialog: false,
  enableBackDefaultBranch: false,
  appBranchList: [] as Array<{ value: string }>,
  appBranchName: '',
  workflowsName: '',
  appName: '',
  appFrameworkList: [{
    Name: '',
    Branch: '',
  }],
  buildManifestList: [],
  manifestSelect: [],
});

// 打开弹窗
const openDialog = (row: any, workflowsName: any, branchName: any) => {
  // console.log(row,workflowsName,branchName)
  state.appBranchName = '';
  state.appBranchList = [];
  state.appFrameworkList = [];
  state.buildManifestList = [];
  state.manifestSelect = [];
  state.isShowDialog = true;
  state.enableBackDefaultBranch = true;
  state.workflowsName = workflowsName;
  state.appName = row.AppName;
  if (branchName) { state.appBranchName = branchName }
  let param = {
    "appName": state.appName
  };

  // 读取应用分支列表
  serverApi.buildBranchList(param).then(function (res) {
    if (res.Status) {
      const arr = res.Data.map((item: any) => {
        return { value: item.BranchName }
      });
      state.appBranchList = arr
    }
  })

  // 读取依赖库列表
  serverApi.appFrameworkList(param).then(function (res) {
    if (res.Status) {
      state.appFrameworkList = res.Data
    }
  })

  // 读取构建历史列表
  serverApi.BuildManifestList(param).then(function (res) {
    if (res.Status) {
      state.buildManifestList = res.Data
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

// 提交构建
const onSubmit = () => {
  if (state.appBranchName) {
    var param = {
      "AppName": state.appName,
      "WorkflowsName": state.workflowsName,
      "BranchName": state.appBranchName,
      "EnableBackDefaultBranch": state.enableBackDefaultBranch,
      "FrameworkList": state.appFrameworkList.map(row => {
        return {
          "FrameworkId": row.Id,
          "CommitId": row.Branch,
        };
      }),
    }

    // console.log(param)
    serverApi.buildAdd(param).then(async function (res) {
      if (res.Status) {
        ElMessage.success("添加成功")
        // 刷新构建日志
        emit('refresh');
        closeDialog();
      } else {
        ElMessage.error(res.StatusMessage)
      }
    })
  } else {
    ElMessage.error('请输入分支名称')
  }
};

// 1. 定义缓存变量
// 用来存储每个仓库对应的分支列表，避免重复请求
const branchCache = ref({});

// 2. 记录当前行
let currentEditRow = null;
const handleFrameworkFocus = (row) => {
  currentEditRow = row;
};

// 3. frameworkSearch 核心逻辑
const frameworkSearch = (queryString, cb) => {
  if (!currentEditRow) {
    cb([]);
    return;
  }

  // 生成唯一的缓存 Key，防止不同仓库的数据串掉
  const cacheKey = `git_${currentEditRow.Id}`;

  // === 情况 A：缓存里有数据 ===
  // 用户是在输入内容，或者第二次点击。我们要做的是“前端过滤”，而不是请求接口。
  const cachedList = branchCache.value[cacheKey];
  if (cachedList) {
    console.log('使用缓存数据进行前端过滤');
    cb(cachedList);
    return;
  }

  var param = {
    "gitId": currentEditRow.Id
  };

  serverApi.RemoteBranchList(param).then(res => {
    if (res.Status && res.Data) {
      const results = res.Data.map(item => {
        return {
          value: item.BranchName,
          CommitId: item.CommitId,
          CommitMessage: item.CommitMessage
        };
      });

      // 存入缓存！下次就不会再请求了
      branchCache.value[cacheKey] = results;

      // 返回数据
      cb(results);
    } else {
      cb([]);
    }
  }).catch(() => {
    cb([]);
  });
};

// 应用分支选择事件
const appBranchNameChange = () => {
  console.log(state.manifestSelect)
  if (state.enableBackDefaultBranch && !state.manifestSelect || JSON.stringify(state.manifestSelect) === '{}') {
      state.appFrameworkList.forEach(curItem => {
        curItem.Branch = state.appBranchName;
      });
  }
};

// 构建清单选择事件
const manifestSelectChange = (item: any) => {
  // 此时 item 是整行对象
  const param = {
    dockerImage: item.DockerImage
  };

  serverApi.BuildManifestDetail(param).then(res => {
    if (res.Status) {
      const apiData = res.Data;
      state.enableBackDefaultBranch = false

      // ================= 需求 1：自动选择 Radio =================
      // 判断返回的数据中是否有 App 和 GitBranch
      if (apiData.App && apiData.App.GitBranch) {
        // 直接赋值，el-radio-group 会自动选中 label 匹配的项
        state.appBranchName = apiData.App.GitBranch;
      }

      // ================= 需求 2：填充 Table 数据 =================
      if (apiData.Dependencies && apiData.Dependencies.length > 0) {
        // 遍历 API 返回的依赖列表
        apiData.Dependencies.forEach(dep => {
          // 1. 在当前表格数据中，查找 Name 匹配的行
          // 注意：这里假设 API 返回的 GitName 和表格里的 Name 是对应的
          const targetRow = state.appFrameworkList.find(row => row.Name === dep.GitName);

          // 2. 如果找到了，更新该行的 Branch 字段
          if (targetRow) {
            // 这里将 GitCommitId 赋值给 Branch 字段
            // 如果你需要的是 GitBranch 而不是 CommitId，请改为 dep.GitBranch
            targetRow.Branch = dep.GitCommitId;
          }
        });
      }
    }
  });
};

// 暴露变量
defineExpose({
  openDialog,
});
</script>

<style scoped lang="scss">
.dialog-content {
  padding: 10px 0;
}

.section-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 20px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.08);
  }

  &:last-child {
    margin-bottom: 0;
  }
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 6px;
  padding-bottom: 12px;
  border-bottom: 2px solid #e4e7ed;

  i {
    font-size: 18px;
    color: #409eff;
    margin-right: 8px;
  }

  .section-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}

.section-body {
  padding: 8px 0;
}

.branch-section {
  background: linear-gradient(135deg, #f5f7fa 0%, #f8f9fa 100%);
}

.branch-radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: flex-start;

  :deep(.el-radio) {
    width: 100%;
    margin-right: 0;
  }

  .branch-radio-item {
    background: white;
    padding: 20px 16px;
    border-radius: 6px;
    border: 1px solid #dcdfe6;
    transition: all 0.3s ease;
    width: 100%;

    &:hover {
      border-color: #409eff;
      box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
    }

    :deep(.el-radio__label) {
      font-size: 14px;
      color: #606266;
    }
  }
}

.framework-section {
  background: linear-gradient(135deg, #fef5f5 0%, #f8f9fa 100%);
}

.framework-checkbox {
  padding-left: 20px;

  :deep(.el-checkbox__label) {
    font-size: 14px;
  }

  .checkbox-label {
    font-weight: 500;
    color: #303133;
  }
}

.checkbox-hint {
  margin-top: 8px;
  margin-left: 24px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.inline-input {
  width: 100%;
}
</style>
