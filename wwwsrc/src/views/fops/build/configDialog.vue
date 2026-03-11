<template>
  <div>
    <el-dialog title="配置管理" v-model="state.isShowDialog" style="width: 80%;top:20px;margin-bottom: 50px;">
        <el-form :model="state.ruleForm" size="default" label-width="120px">
          <el-form-item label="应用名称">
            <el-input v-model="state.ruleForm.appName" disabled></el-input>
          </el-form-item>
          <el-form-item label="配置版本">
            <el-space>
              <el-tag type="info">应用版本: {{ state.ruleForm.appConfigVer }}</el-tag>
              <el-tag :type="state.ruleForm.dockerConfigVer === '未创建' ? 'warning' : 'success'">Docker版本: {{ state.ruleForm.dockerConfigVer }}</el-tag>
            </el-space>
          </el-form-item>
          <el-form-item label="配置内容">
            <el-input v-model="state.ruleForm.content" type="textarea" :rows="500" placeholder="请输入配置内容（YAML格式）" style="font-family: 'Courier New', monospace; font-size: 13px;"></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="onCancel" size="default">取消</el-button>
            <el-button type="primary" @click="onSubmit" size="default" :loading="state.loading">保存</el-button>
          </span>
        </template>
    </el-dialog>
  </div>
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
    appConfigVer: 0,
    dockerConfigVer: '未创建',
  },
});

// 打开弹窗
const openDialog = (appName: string) => {
  state.ruleForm.appName = appName;
  state.isShowDialog = true;
  state.loading = true;

  // 获取配置内容
  serverApi.getConfig(appName).then((res: any) => {
    state.loading = false;
    if (res.Status) {
      // 处理新的响应格式
      const data = res.Data;
      state.ruleForm.content = data.Content || data;
      state.ruleForm.appConfigVer = data.AppConfigVer || 0;
      state.ruleForm.dockerConfigVer = data.DockerConfigVer || '未创建';
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
  line-height: 2;
  height: 1000px;
}
</style>
