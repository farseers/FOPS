<template>
  <div>
    <el-dialog title="构建配置" v-model="state.isShowDialog" width="600px" :close-on-click-modal="true">
      <div class="dialog-content">
        <!-- 分支选择区域 -->
        <div class="section-card branch-section">
          <div class="section-header">
            <i class="el-icon-branch"></i>
            <span class="section-title">选择应用分支</span>
          </div>
          <div class="section-body">
            <el-radio-group v-model="state.inputValue" class="branch-radio-group">
              <el-radio v-for="item in state.restaurants" :key="item.value" :label="item.value" size="default"
                class="branch-radio-item">
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
            <el-checkbox v-model="state.updateFramework" class="framework-checkbox">
              <span class="checkbox-label">使用最新的框架版本</span>
            </el-checkbox>
          </div>
          <div class="section-body">
            
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
import { reactive } from 'vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from "/@/api/fops";
const emit = defineEmits(['refresh']);
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
  isShowDialog: false,
  restaurants: [] as Array<{ value: string }>,
  inputValue: '',
  workflowsName: '',
  appName: '',
  updateFramework: false
});

// 打开弹窗
const openDialog = (row: any, workflowsName: any, branchName: any) => {
  // console.log(row,workflowsName,branchName)
  state.inputValue = '';
  state.restaurants = [];
  state.isShowDialog = true;
  state.workflowsName = workflowsName;
  state.appName = row.AppName;
  state.updateFramework = false;
  if (branchName) { state.inputValue = branchName }
  let param = {
    "appName": state.appName
  };

  serverApi.autobuildBranchList(param).then(function (res) {
    if (res.Status) {
      const arr = res.Data.map((item: any) => {
        return { value: item.BranchName }
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
  if (state.inputValue) {
    var param = {
      "AppName": state.appName,
      "WorkflowsName": state.workflowsName,
      "branchName": state.inputValue,
      "updateFramework": state.updateFramework
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
  margin-bottom: 16px;
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
  padding-left: 10px;
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
