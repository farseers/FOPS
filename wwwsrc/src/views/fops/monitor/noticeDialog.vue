<template>
    <div class="system-role-dialog-container">
        <el-dialog :title="title" v-model="isShowDialog" width="769px">
            <el-form-item label="通知类型">
              <el-select v-model="infoRow.NoticeType"  placeholder="请选择通知类型" >
                <el-option
                  v-for="item in typeList"
                  :key="item.NoticeType"
                  :label="item.NoticeTypeName"
                  :value="item.NoticeType"
                />
              </el-select>
              <!-- <el-input v-model="infoRow.NoticeType" /> -->
            </el-form-item>
            <el-form-item label="姓名"><el-input v-model="infoRow.Name" /></el-form-item>
            <el-form-item label="邮箱"><el-input v-model="infoRow.Email" /></el-form-item>
            <el-form-item label="号码"><el-input v-model="infoRow.Phone" /></el-form-item>
            <el-form-item label="接口Key"><el-input v-model="infoRow.ApiKey" /></el-form-item>
            <el-form-item label="是否启用">
            <el-switch v-model="infoRow.Enable" style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" active-text="启用" inactive-text="关闭" />
        </el-form-item>
            <el-form-item label="备注"><el-input v-model="infoRow.Remark" /></el-form-item>
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
import {  ElMessage } from 'element-plus';
const serverApi = fopsApi();
export default{
    data(){
        return {
            title:'',
            isShowDialog:false,
            infoRow:{
              "Id":null,
              "Name":'',
              "NoticeType":0,
              "Email":"",
              "Phone":"",
              "ApiKey":"",
              "Remark":"",
              Enable:true
            }
            
        }
    },
    props:{
      typeList:{
        type:Array,
        default:()=>{
          return []
        }
      }
    },
    methods:{
        onCancel(){
            this.isShowDialog = false;
            this.infoRow = {
              "Id":null,
              "Name":'',
              "NoticeType":0,
              "Email":"",
              "Phone":"",
              "ApiKey":"",
              "Remark":"",
              Enable:true
            }
        },
        onSubmit(){
            let param = {...this.infoRow};
            const daterange = param.daterange;
            if(daterange && daterange.length > 0 ){
                param.StartTime = daterange[0];
                param.EndTime = daterange[1];
            }
            serverApi.monitorSaveNotice(param).then(d=>{
                let { Status,StatusMessage } = d;
                if(Status){
                    this.onCancel()
                    this.$emit('search')
                }else{
                    ElMessage.error(StatusMessage)
                }
            })
        },
        info(id){
          if(id){
            this.title = '编辑通知用户';
            serverApi.monitorInfoNotice({id:id}).then(d=>{
                let { Data,Status,StatusMessage } = d;
                if(Status){
                  
                    this.infoRow = {...Data,daterange:[Data.StartTime,Data.EndTime]}
                    this.isShowDialog = true;
                }else{
                    ElMessage.error(StatusMessage)
                }
            })
          }else{
            this.title = '新增通知用户';
            this.infoRow = {
              "Id":null,
              "NoticeType":0,
              'Name':'',
              "Email":"",
              "Phone":"",
              "ApiKey":"",
              "Remark":"",
              Enable:true
            }
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