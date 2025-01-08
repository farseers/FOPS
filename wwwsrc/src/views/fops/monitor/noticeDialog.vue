<template>
  <div class="system-role-dialog-container">
    <el-dialog :title="title" v-model="isShowDialog" width="769px">
      <el-form ref="ruleFormRef" :model="infoRow" :rules="rules">
       
        <el-form-item>
          <el-form-item label="通知类型" prop='NoticeType' style="flex: 1;">
          <el-select v-model="infoRow.NoticeType" placeholder="请选择通知类型" style="flex: 1;">
            <el-option v-for="item in typeList" :key="item.Key" :label="item.Value" :value="item.Key" />
          </el-select>
          <!-- <el-input v-model="infoRow.NoticeType" /> -->
        </el-form-item>
          <el-form-item label="姓名" prop='Name' style="width: 250px;padding-left: 5px;"><el-input v-model="infoRow.Name" /></el-form-item>
          <el-form-item label="号码" prop='Phone' style="width: 250px;padding-left: 5px;"><el-input v-model="infoRow.Phone" /></el-form-item>
        </el-form-item>
        <el-form-item label="接口Key" prop='ApiKey'><el-input v-model="infoRow.ApiKey" /></el-form-item>
        <el-form-item>
          <el-form-item label="邮箱" style="width: 300px;padding-right: 5px;"><el-input v-model="infoRow.Email" /></el-form-item>
          <el-form-item label="是否启用" style="flex: 1;">
            <el-switch v-model="infoRow.Enable" style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
              active-text="启用" inactive-text="停用" />
          </el-form-item>
        </el-form-item>


        <el-form-item label="备注"><el-input v-model="infoRow.Remark" /></el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="onCancel" size="default">取 消</el-button>
          <el-button type="primary" @click="onSubmit" size="default">保 存</el-button>
        </span>
      </template>
    </el-dialog>


  </div>
</template>

<script>

import { fopsApi } from "/@/api/fops";
import { ElMessage } from 'element-plus';
const serverApi = fopsApi();
export default {
  data() {
    return {
      rules: {
        NoticeType: [{ required: true, trigger: 'change', message: '请选择通知类型' }],
        Name: [{ required: true, trigger: 'blur', message: '请输入姓名' }],
        Phone: [{ required: true, trigger: 'blur', message: '请输入号码' }],
        ApiKey: [{ required: true, trigger: 'blur', message: '请输入接口Key' }],
      },
      title: '',
      isShowDialog: false,
      infoRow: {
        "Id": null,
        "Name": '',
        "NoticeType": 0,
        "Email": "",
        "Phone": "",
        "ApiKey": "",
        "Remark": "",
        Enable: true
      }

    }
  },
  props: {
    typeList: {
      type: Array,
      default: () => {
        return []
      }
    }
  },
  methods: {
    onCancel() {
      this.isShowDialog = false;
      this.init()
    },
    onSubmit() {
      this.$refs.ruleFormRef && this.$refs.ruleFormRef.validate((valid) => {
        if (valid) {
          let param = { ...this.infoRow };
          const daterange = param.daterange;
          if (daterange && daterange.length > 0) {
            param.StartTime = daterange[0];
            param.EndTime = daterange[1];
          }
          serverApi.monitorSaveNotice(param).then(d => {
            let { Status, StatusMessage } = d;
            if (Status) {
              this.onCancel()
              this.$emit('search')
            } else {
              ElMessage.error(StatusMessage)
            }
          })
        }
      })

    },
    init() {
      this.infoRow = {
        "Id": null,
        "NoticeType": 0,
        'Name': '',
        "Email": "",
        "Phone": "",
        "ApiKey": "",
        "Remark": "",
        Enable: true
      }
      this.$refs.ruleFormRef && this.$refs.ruleFormRef.resetFields()
    },
    info(id) {
      this.init()
      if (id) {
        this.title = '编辑通知用户';
        serverApi.monitorInfoNotice({ id: id }).then(d => {
          let { Data, Status, StatusMessage } = d;
          if (Status) {
            this.infoRow = { ...Data, daterange: [Data.StartTime, Data.EndTime] }
            this.isShowDialog = true;
          } else {
            ElMessage.error(StatusMessage)
          }
        })
      } else {
        this.title = '新增通知用户';

        this.isShowDialog = true;
      }

    }
  }
}
</script>
<!--
<template>
  <el-form :model="form" label-width="120px">
      <el-form-item label="Activity name">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="Activity zone">
        <el-select v-model="form.region" placeholder="please select your zone">
          <el-option label="Zone one" value="shanghai" />
          <el-option label="Zone two" value="beijing" />
        </el-select>
      </el-form-item>
      <el-form-item label="Activity time">
        <el-col :span="11">
          <el-date-picker
            v-model="form.date1"
            type="date"
            placeholder="Pick a date"
            style="width: 100%"
          />
        </el-col>
        <el-col :span="2" class="text-center">
          <span class="text-gray-500">-</span>
        </el-col>
        <el-col :span="11">
          <el-time-picker
            v-model="form.date2"
            placeholder="Pick a time"
            style="width: 100%"
          />
        </el-col>
      </el-form-item>
      <el-form-item label="Instant delivery">
        <el-switch v-model="form.delivery" />
      </el-form-item>
      <el-form-item label="Activity type">
        <el-checkbox-group v-model="form.type">
          <el-checkbox label="Online activities" name="type" />
          <el-checkbox label="Promotion activities" name="type" />
          <el-checkbox label="Offline activities" name="type" />
          <el-checkbox label="Simple brand exposure" name="type" />
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label="Resources">
        <el-radio-group v-model="form.resource">
          <el-radio label="Sponsor" />
          <el-radio label="Venue" />
        </el-radio-group>
      </el-form-item>
      <el-form-item label="Activity form">
        <el-input v-model="form.desc" type="textarea" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">Create</el-button>
        <el-button>Cancel</el-button>
      </el-form-item>
    </el-form>
  </template>
  
  <script lang="ts" setup>
  import { reactive } from 'vue'
  
  // do not use same name with ref
  const form = reactive({
    name: '',
    region: '',
    date1: '',
    date2: '',
    delivery: false,
    type: [],
    resource: '',
    desc: '',
  })
  
  const onSubmit = () => {
    console.log('submit!')
  }
  </script>
   -->