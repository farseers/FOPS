<template>

  <div class="system-role-dialog-container">
    <el-dialog :title="title" v-model="isShowDialog" width="800px">
      <el-form ref="ruleFormRef">
          <el-form-item label="数据库类型">
            <el-radio-group v-model="BackupDataType">
            <el-radio :label="0">Mysql</el-radio>
            <el-radio :label="1">Clickhouse</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item style="display: flex;">
          <div style="flex: 1;padding-right: 10px;">
            <el-form-item label="主机"><el-input v-model="Host" placeholder="请输入主机"/></el-form-item>
          </div>
          <div style="flex: 1;padding-right: 10px;">
            <el-form-item label="端口"><el-input v-model="Port" placeholder="请输入端口" @input="onPort"/></el-form-item>
          </div>
        </el-form-item>
        <el-form-item style="display: flex;">
          <div style="flex: 1;padding-right: 10px;">
            <el-form-item label="用户名"><el-input v-model="Username" placeholder="请输入用户名"/></el-form-item>
          </div>
          <div style="flex: 1;padding-right: 10px;">
            <el-form-item label="密码"><el-input v-model="Password" placeholder="请输入密码"/></el-form-item>
          </div>
          <div style="flex: 1;padding-right: 10px;">
            <el-form-item label="Cron"><el-input v-model="Cron" placeholder="请输入Cron"/></el-form-item>
          </div>
        </el-form-item>
        <el-form-item  label="数据库">
            <div style="width: 100%;">
            <div style="flex: 1;padding-right: 10px;display: flex;align-items: center;">
            <el-form-item style="width: 100%;">
              <el-input v-model="baseTit" style="width: 200px;" clearable placeholder="请输入数据库"/>
              <el-button  style="margin-left: 10px;" type="success" @click="oAddBase">添加</el-button>
              <el-button  slot="reference" style="margin-left: 10px;" @click="baseCh" type="primary">查询</el-button>
            </el-form-item>
          </div>
          </div>
          <div style="width: 100%;display: flex;">
            <span v-show="addBases.length<=0" style="color: #909399;">请添加或选择数据库</span>
            <el-checkbox-group v-model="checkBase" size="small" style="margin-top: 5px;">
              <el-checkbox v-for="t,i in addBases" border  :label="t" :key="i">{{ t }}</el-checkbox>
            </el-checkbox-group>
          
         </div>
        </el-form-item>
        
        <el-form-item label="存储类型">
           <div style="width: 100%;">
            <el-radio-group v-model="StoreType">
              <el-radio :label="0">OSS</el-radio>
              <el-radio :label="1">本地目录</el-radio>
          </el-radio-group>
           </div>
           <div v-show="StoreType == 1" style="width: 100%;">
            <div style="flex: 1;padding-right: 10px;">
                <el-form-item label="目录"><el-input v-model="Directory" style="flex: 1;" placeholder="请输入目录"/></el-form-item>
              </div>
           </div>
           <div v-show="StoreType == 0" style="width: 100%;">
            <div style="width: 100%;display: flex;margin-bottom: 10px;">
              <div style="flex: 1;padding-right: 10px;">
                <el-form-item label="AccessKeyID"><el-input v-model="AccessKeyID" style="flex: 1;" placeholder="请输入AccessKeyID"/></el-form-item>
              </div>
              <div style="flex: 1;padding-right: 10px;">
                <el-form-item label="AccessKeySecret"><el-input v-model="AccessKeySecret" style="flex: 1;" placeholder="请输入AccessKeySecret"/></el-form-item>
              </div>
            
            </div>
            <div style="width: 100%;display: flex;margin-bottom: 10px;">
              <div style="flex: 1;padding-right: 10px;">
                <el-form-item label="BucketName"><el-input v-model="BucketName" style="flex: 1;" placeholder="请输入BucketName"/></el-form-item>
              </div>
              <div style="flex: 1;padding-right: 10px;">
                <el-form-item label="访问结点"><el-input v-model="Endpoint" style="flex: 1;" placeholder="请输入访问结点"/></el-form-item>
              </div>
            </div>
            <div style="width: 100%;display: flex;">
              <div style="flex: 1;padding-right: 10px;">
                  <el-form-item label="区域"><el-input v-model="Region" style="flex: 1;" placeholder="请输入区域"/></el-form-item>
                </div>
            </div>
           </div>
        </el-form-item>

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
      title: '编辑',
      isShowDialog: false,
      baseTit:'',//数据库输入
      "BackupDataType": 0,// 数据库类型：0 = Mysql, 1 = Clickhouse
      "Host": "", // 主机
      "Port": 0, // 端口
      "Username": "", // 用户名
      "Password": "", // 密码
      "Cron": "",  // Cron (字符串，文本框，如：0 0 0/1 * * ?）
      "StoreType": 0, // 存储类型：0 = OSS， 1= 本地目录
      "AccessKeyID": "", // AccessKeyID
      "AccessKeySecret": "",// AccessKeySecret
      "Endpoint": "", // 访问结点，如：https://oss-cn-hangzhou.aliyuncs.com
      "Region": "", // 区域,如：cn-hangzhou
      "BucketName": "" ,// BucketName
      "Directory": "", // 目录
      addBases:[],//选择的数据库
      checkBase:[],
      Id:null,
    }
  },
  methods: {
    oAddBase(){
      if(this.baseTit && !this.addBases.includes(this.baseTit)){
        this.checkBase.push(this.baseTit)
        this.addBases.push(this.baseTit)
        this.baseTit = '';
      }
    },
    onPort(){
      this.Port = this.Port *1;
    },
    changeStore(){
      this.Directory = '';
      this.AccessKeyID = '';
      this.AccessKeySecret = '';
      this.Endpoint = '';
      this.Region = '';
      this.BucketName = '';
    },
    baseCh(){ //数据库列表
      serverApi.backupData_getDatabaseList({
          "BackupDataType": this.BackupDataType *1,
          "Host": this.Host,
          "Port": this.Port,
          "Username": this.Username,
          "Password": this.Password
        }).then(d => {
            let { Status, StatusMessage,Data } = d;
            if (Status) {
              console.log(Data)
              this.baseData = [...Data]
            } else {
              ElMessage.error(StatusMessage)
            }
          })
    },
    onCancel() {
      this.Id = null;
      this.baseData = []
      this.BackupDataType = 0;
      this.baseTit = '';
      this.Host = '';
      this.Port = '';
      this.Username = '';
      this.Password = '';
      this.addBases = [];
      this.Cron = '';
      this.StoreType = 0;
      this.Directory = '';
      this.AccessKeyID = '';
      this.AccessKeySecret = '';
      this.Endpoint = '';
      this.Region = '';
      this.BucketName = '';
      this.isShowDialog = false;
     
      this.$emit('search')
    },
    onSubmit() {
      let baseData = this.checkBase;
      let StoreType = this.StoreType * 1;
      let StoreConfig = {};
      if(StoreType == 0){
        StoreConfig = {
          "AccessKeyID": this.AccessKeyID, // AccessKeyID
          "AccessKeySecret": this.AccessKeySecret,// AccessKeySecret
          "Endpoint": this.Endpoint, // 访问结点，如：https://oss-cn-hangzhou.aliyuncs.com
          "Region": this.Region, // 区域,如：cn-hangzhou
          "BucketName": this.BucketName ,// BucketName
        }
      }
      if(StoreType == 1){
        StoreConfig = {
          Directory:this.Directory
        }
        
      }
      let param = {
        Id:this.Id ,
        "BackupDataType": this.BackupDataType * 1,// 数据库类型：0 = Mysql, 1 = Clickhouse
        "Host": this.Host, // 主机
        "Port": this.Port * 1, // 端口
        "Username": this.Username, // 用户名
        "Password": this.Password, // 密码
        "Database": baseData, //数据库（多个）
        "Cron": this.Cron,  // Cron (字符串，文本框，如：0 0 0/1 * * ?）
        "StoreType": StoreType, // 存储类型：0 = OSS， 1= 本地目录
        "StoreConfig": JSON.stringify(StoreConfig)// 存储配置（根据存储类型0或1，对应的配置会不同） 列表页不展示
      }
          serverApi.backupData_add(param).then(d => {
            let { Status, StatusMessage } = d;
            if (Status) {
              this.onCancel()
            } else {
              ElMessage.error(StatusMessage)
            }
          })
    },
    info(id) {
      this.baseData = []
       this.title = '新增'
      if (id) {
        this.title = '编辑'
        serverApi.backupData_info({ id: id }).then(d => {
          let { Data, Status, StatusMessage } = d;
          if (Status) {
            const row = Data;
            this.Id = row.Id;
            this.BackupDataType = row.BackupDataType;
            this.Host = row.Host;
            this.Port = row.Port;
            this.Username = row.Username;
            this.Password = row.Password;
            this.Cron = row.Cron;
            this.StoreType = row.StoreType;
            this.BackupDataType = row.BackupDataType;
            this.BackupDataType = row.BackupDataType;
            if(row.StoreConfig){
              var StoreConfig = JSON.parse(row.StoreConfig)
              if(this.StoreType == 0){
                this.AccessKeyID=StoreConfig.AccessKeyID; // AccessKeyID
                this.AccessKeySecret= StoreConfig.AccessKeySecret;// AccessKeySecret
                this.Endpoint=StoreConfig.Endpoint; // 访问结点，如：https://oss-cn-hangzhou.aliyuncs.com
                this.Region=StoreConfig.Region; // 区域,如：cn-hangzhou
                this.BucketName=StoreConfig.BucketName;// BucketName
              }
              if(this.StoreType == 1){
                this.Directory = StoreConfig.Directory
              }
            }
            
            // Database
            this.checkBase = row.Database;
            this.addBases = row.Database;
            this.isShowDialog = true;
          } else {
            ElMessage.error(StatusMessage)
          }
        })
      } else {
        this.isShowDialog = true;
      }

    }
  }
}
</script>
<style scoped>
.baseCls{
    border: 1px solid #909399;
    color: #909399;
    background-color: #fff;
    border-radius: 5px;
    overflow: hidden;
    margin: 5px;
    font-size: 14px;
    border: 1px solid #409EFF;
    color: #409EFF;
    display: flex;
    height: 28px;
}
.baseCls span{
    display: flex;
    align-items: center;
    padding: 3px 5px;
}
.baseCls .s2{
  cursor: pointer;
  color: #409EFF;
 
}
.baseCls .s2:hover{
  color: #fff;
  background-color: #F56C6C;
}
</style>