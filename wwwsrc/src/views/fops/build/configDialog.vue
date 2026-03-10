<template>
  <el-dialog title="配置管理" v-model="state.isShowDialog" width="80%" top="50px">
    <el-form :model="state.ruleForm" size="default" label-width="100px">
      <el-form-item label="应用名称">
        <el-input v-model="state.ruleForm.appName" disabled></el-input>
      </el-form-item>
      <el-form-item label="配置内容">
        <el-input
          v-model="state.ruleForm.content"
          type="textarea"
          :rows="20"
          placeholder="请输入配置内容（YAML格式）"
          style="font-family: 'Courier New', monospace; font-size: 13px;"
        ></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="onCancel" size="default">取消</el-button>
        <el-button type="primary" @click="onSubmit" size="default" :loading="state.loading">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="configDialog">
import { reactive } from 'vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from '/@/api/fops';

const serverApi = fopsApi();

const state = reactive({
  isShowDialog: false,
  loading: false,
  ruleForm: {
    appName: '',
    content: '',
  },
});

// 打开弹窗
const openDialog = (appName: string) => {
  state.ruleForm.appName = appName;
  state.isShowDialog = true;
  state.loading = true;

  // 获取配置内容
  serverApi.getConfig({ appName }).then((res: any) => {
    state.loading = false;
    if (res.Status) {
      state.ruleForm.content = res.Data;
    } else {
      ElMessage.error(res.StatusMessage || '获取配置失败');
    }
  }).catch(() => {
    state.loading = false;
    ElMessage.error('获取配置失败');
  });
};

// 取消
const onCancel = () => {
  state.isShowDialog = false;
};

// 提交
const onSubmit = () => {
  if (!state.ruleForm.content.trim()) {
    ElMessage.warning('配置内容不能为空');
    return;
  }

  state.loading = true;
  serverApi.saveConfig({
    appName: state.ruleForm.appName,
    content: state.ruleForm.content,
  }).then((res: any) => {
    state.loading = false;
    if (res.Status) {
      ElMessage.success('保存成功');
      state.isShowDialog = false;
    } else {
      ElMessage.error(res.StatusMessage || '保存失败');
    }
  }).catch(() => {
    state.loading = false;
    ElMessage.error('保存失败');
  });
};

// 暴露方法
defineExpose({
  openDialog,
});
</script>

<style scoped lang="scss">
:deep(.el-textarea__inner) {
  font-family: 'Courier New', Consolas, Monaco, monospace;
  font-size: 13px;
  line-height: 1.5;
}
</style>
