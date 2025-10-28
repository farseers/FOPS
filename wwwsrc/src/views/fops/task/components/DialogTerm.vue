<template>
    <div>
        <el-dialog v-model="dialogVisible" @close="close" custom-class="fixed-dialog">
            <div style="display: flex;">
                <el-form-item label="IP" style="width: 200px;margin-right: 5px;"><el-input
                        v-model="pRow.LoginIp" size='small' @keyup.enter="submitForm()"/></el-form-item>
                <el-form-item label="登录名" style="width: 180px;margin-right: 5px;"><el-input
                        v-model="pRow.LoginName" size='small' @keyup.enter="submitForm()"/></el-form-item>
                <el-form-item label="登录密码" style="width: 180px;margin-right: 5px;"><el-input
                        v-model="pRow.LoginPwd" size='small' @keyup.enter="submitForm()"/></el-form-item>
                <el-form-item label="端口号" style="width: 120px;margin-right: 5px;"><el-input
                        v-model="pRow.LoginPort" size='small' @keyup.enter="submitForm()"/></el-form-item>
                <el-form-item><el-button type="success" size='small' @click="submitForm()">连接</el-button></el-form-item>
            </div>
            <div style="flex: 1;">
                <InitTerm ref="initTerm" />
            </div>
        </el-dialog>
    </div>
</template>
<script>
import InitTerm from '/src/views/components/InitTerm.vue';
import { ElMessage } from 'element-plus';
import { fopsApi } from "/@/api/fops";
const serverApi = fopsApi();
const defaultRow = {
    LoginIp: '',
    LoginName: 'root',
    LoginPwd: '',
    LoginPort: 22
}
export default {
    name: 'DialogTerm',
    components: { InitTerm },
    data() {
        return {
            dialogVisible: false,
            pRow: {
                ...defaultRow
            }
        }
    },
    methods: {
        onOpenEditRole(row){
			const _this = this;
            //  console.log('onOpenEditRole',row)
             const LoginIp = row.LoginIp;
			serverApi.terminalClientInfo({LoginIp:LoginIp}).then((res)=>{
				if(res.Status){
                    // console.log(res)
                    const d = res.Data;
                    const { LoginName,LoginPort,LoginPwd } = d;
                     this.LoginName =LoginName;
                     if(LoginPwd){this.LoginPwd = LoginPwd;}
                     this.LoginPort = LoginPort;
                    if(res.StatusCode != 403){
                        _this.$refs.initTerm && _this.$refs.initTerm.inits(LoginIp)
                        // {
                        //     LoginName:LoginName,
                        //     LoginPort:LoginPort
                        // }
                    }
					
				}else{
				ElMessage.error(res.StatusMessage);
			}
			})
		},
        init(row) {
            this.pRow = { ...defaultRow }
            this.$refs.initTerm && this.$refs.initTerm.clearWs()
            this.pRow.LoginIp = row.IP;
            this.dialogVisible = true;
            this.onOpenEditRole(this.pRow)
            // submitForm();
        },
        submitForm() {
            if(this.pRow.LoginIp ){ // && this.pRow.LoginName && this.pRow.LoginPwd && this.pRow.LoginPort
                this.$refs.initTerm && this.$refs.initTerm.initStart(this.pRow)
            }
          
        },
        close() {
            this.$refs.initTerm && this.$refs.initTerm.clearWs()
            this.dialogVisible = false;
        }
    }
}
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
        display: flex;
        flex-flow: column;
    }
}
</style>