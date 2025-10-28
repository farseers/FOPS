<template>
  <el-dialog v-model="dialogVisible" 
  @close="close"
  custom-class="fixed-dialog">
    <el-tabs v-model="activeName">
      <el-tab-pane name="first" label="SSH">
        <div style="text-align: center">
          <el-form ref="elForm" :model="form" status-icon :rules="rules" label-position="left" label-width="80px">
            <el-form-item label="名称" prop="name">
              <el-input v-model="form.name" />
            </el-form-item>
            <el-form-item label="Ip">
              <el-form-item prop="ip" style="width: 300px;">
              <el-input v-model="form.ip" />
            </el-form-item>
            <el-form-item label="端口" style="padding-left: 10px;width: 60px;" prop="port">
              <el-input v-model="form.port" />
            </el-form-item>
            </el-form-item>
           
            <el-form-item label="登录名" prop="user">
              <el-input v-model="form.user" />
            </el-form-item>
            <el-form-item label="登录密码" prop="pwd">
              <el-input v-model="form.pwd" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="save()">保存</el-button>
              <el-button type="success" @click="submitForm()"  v-if="editIp">重新连接</el-button>
              <el-button @click="resetForm()">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>
      <el-tab-pane name="second" label="Terminal">
          <InitTerm ref="initTerm"/>
      </el-tab-pane>  
    </el-tabs>

  </el-dialog>
</template>

<script>
import InitTerm from '/src/views/components/InitTerm.vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from "/@/api/fops";
const serverApi = fopsApi();
const packResize = (cols, rows) =>
  JSON.stringify({
    type: 'resize',
    cols: cols,
    rows: rows
  })
  const defaultRow= {
        user: 'root',
        pwd: '',
        ip: '',
        name: '',
        port: '22',
  }
export default {
  name: 'termPane',
  created() {

  },
  components:{InitTerm},
  data() {
    var validate = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(''));
      } else {
        callback();
      }
    };
    return {
      dialogVisible: false,
      secondShow:false,
     
      rules: {
        ip: [
          { validator: validate, trigger: 'blur' }
        ],
        name: [
          { validator: validate, trigger: 'blur' }
        ],
        port: [
          { validator: validate, trigger: 'blur' }
        ],
        user: [
          { validator: validate, trigger: 'blur' }
        ],
        pwd: [
          { validator: validate, trigger: 'blur',validator:this.validatorPwd }
        ],
      },
      activeName: 'first',
      first: true,
      term: null,
      fitAddon: null,
      ws: null,
      form: {
        ...defaultRow
      },
      editIp: null,
      inputBuffer:'',//输入的字符
     
    }
  },
  mounted() {

  },
  methods: {
    validatorPwd(){
          if(this.editIp){
            return true
          }else{
            if(this.form.pwd){
              return true
            }else{
              return false
            }
          }
    },
   
    open() { //新增
      this.activeName = 'first';
      this.form = {
        ...defaultRow
      }
      this.editIp = null;
      this.secondShow = false;
      this.clearWs()
      this.$refs['elForm'] && this.$refs['elForm'].resetFields();
      this.dialogVisible = true;
    },
    edit(row) { //修改
      this.activeName = 'second';
      this.form = {
        user: row.LoginName,
        pwd: row.LoginPwd,
        ip: row.LoginIp,
        name: row.Name,
        port: row.LoginPort,
      }
      this.editIp = row.LoginIp;
      this.secondShow = true;
      this.dialogVisible = true;
      this.initWs()
    },
    close() {
      this.clearWs()
      this.form = {
       ...defaultRow
      }
      this.dialogVisible = false;
      this.secondShow = false;
      this.activeName = 'first';
      this.editIp = null;
      this.$refs['elForm'].resetFields();
    },
    save_back(res,fn) {
      if (res.Status) {
        if(fn){
          fn && fn() //链接
        }else{
          this.close()
        }
       
      } else {
        ElMessage.error(res.StatusMessage);
      }
    },
    save(fn) {
      const _this = this;
      this.$refs['elForm'].validate((valid) => {
        if (valid) {
          let port = _this.form.port *1;
          if(isNaN(port)){port = 0}
          let param = {
            LoginName: _this.form.user,
            LoginIp: _this.form.ip,
            Name: _this.form.name,
            LoginPort: port,
          }
          if(_this.form.pwd){
            param.LoginPwd = _this.form.pwd
          }
          if (_this.editIp) {
            param.LoginIp = _this.editIp;
            serverApi.terminalClientUpdate(param).then(function (res) {
              _this.$emit('refresh')
              _this.save_back(res,fn)
             
            })
          } else {
            serverApi.terminalClientAdd(param).then(function (res) {
              if(res.Status){
                _this.$emit('refresh')
                _this.close()
              }
            })
          }

        } else {
          ElMessage.error('请填写完整');
          return false;
        }
      });
    },
    clearWs(type){
      if(type){ //清楚后 重连
        this.$refs.initTerm&&this.$refs.initTerm.clearInit(this.editIp)
      }else{
        this.$refs.initTerm&&this.$refs.initTerm.clearWs()
      }
        
    },
    initWs() {
     setTimeout(()=>{
      this.$refs.initTerm&&this.$refs.initTerm.init(this.editIp)
     },300)
      
    },
    submitForm() { //链接
      this.$refs['elForm'].validate((valid) => {
        if (valid) {
          this.save(()=>{//保存后跳到链接
            this.secondShow = true;
            this.activeName = 'second';
            this.clearWs(true)
          }) 
        } else {
          ElMessage.error('请填写完整');
          return false;
        }
      });
    },
    resetForm() {
      this.$refs['elForm'].resetFields();
    },
    
  }
}
//
</script>
<style lang="scss">
.fixed-dialog {
  width: 90%;
  height: 90%;
  display: flex;
  flex-flow: column;

  .el-dialog__body {
    flex: 1;
    max-height: none !important;

    .el-tabs {
      height: 100%;
      display: flex;
      flex-flow: column;

      .el-tabs__content {
        flex: 1;

        .el-tab-pane {
          height: 100%;

        }
      }
    }

  }
}
</style>