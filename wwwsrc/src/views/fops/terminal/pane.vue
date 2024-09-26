<template>
  <el-dialog v-model="dialogVisible" 
  @close="close"
  @before-close="beforeClose"
  custom-class="fixed-dialog">
    <el-tabs v-model="activeName">
      <el-tab-pane name="first" label="SSH">
        <div style="text-align: center">
          <el-form ref="elForm" :model="form" status-icon :rules="rules" label-position="left" label-width="80px">
            <el-form-item label="Ip" prop="ip">
              <el-input v-model="form.ip" />
            </el-form-item>
            <el-form-item label="名称" prop="name">
              <el-input v-model="form.name" />
            </el-form-item>
            <el-form-item label="端口" prop="port">
              <el-input v-model="form.port" />
            </el-form-item>
            <el-form-item label="登录名" prop="user">
              <el-input v-model="form.user" />
            </el-form-item>
            <el-form-item label="登录密码" prop="pwd">
              <el-input v-model="form.pwd" />
            </el-form-item>
            <el-form-item v-if="!editId">
              <el-button type="success" @click="submitForm()">连接</el-button>
              <el-button @click="resetForm()">重置</el-button>
            </el-form-item>
            <el-form-item v-if="editId">
              <el-button type="success" @click="submitForm()">重新连接</el-button>
              <el-button @click="resetForm()">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>
      <el-tab-pane name="second" label="Terminal" v-if="secondShow">
        <div class="ssh-container" ref="terminal"></div>
      </el-tab-pane>
    </el-tabs>

  </el-dialog>
</template>

<script>

import { Session } from '/@/utils/storage';
import { ElMessage } from 'element-plus';
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { debounce } from 'lodash'
import { fopsApi } from "/@/api/fops";
const serverApi = fopsApi();
const packResize = (cols, rows) =>
  JSON.stringify({
    type: 'resize',
    cols: cols,
    rows: rows
  })
export default {
  name: 'termPane',
  created() {

  },
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
          { validator: validate, trigger: 'blur' }
        ],
      },
      activeName: 'first',
      first: true,
      term: null,
      fitAddon: null,
      ws: null,
      form: {
        user: '',
        pwd: '',
        ip: '',
        name: '',
        port: '',
      },
      editId: null,
      inputBuffer:'',//输入的字符
      option: {
        cursorBlink: true,
        cursorStyle: 'underline', // 光标样式 'block' | 'underline' | 'bar'
        fontSize: 14, // 调整字体大小
        letterSpacing: 1, // 字符间距
        lineHeight: 1, // 行高
        fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
        theme: {
          background: '#181d28'
        },
        cols: 30 // 初始化的时候不要设置fit，设置col为较小值（最小为可展示initText初始文字即可）方便屏幕缩放
      }
    }
  },
  mounted() {

  },
  beforeDestroy() {
    // console.log('beforeDestroy')
   this.clearWs()
  },
  
  methods: {
    clearWs(){
      if (this.ws) {
        this.ws.close();
        this.ws = null; 
      }
      if(this.term) {
        this.term.dispose(); 
        this.term.clear();
        this.term = null; 
      }
  },
    open() { //新增
      this.activeName = 'first';
      this.form = {
        user: '',
        pwd: '',
        ip: '',
        name: '',
        port: '',
      }
      this.editId = null;
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
      this.editId = row.Id;
      this.secondShow = true;
      this.dialogVisible = true;
      this.initWs()
    },
    beforeClose(){
      this.clearWs()
    },
    close() {
      this.form = {
        user: '',
        pwd: '',
        ip: '',
        name: '',
        port: '',
      }
      this.dialogVisible = false;
      this.secondShow = false;
      this.activeName = 'first';
      this.editId = null;
      this.$refs['elForm'].resetFields();
    },
    save_back(res,fn) {
      if (res.Status) {
        if(fn){
          fn && fn() //链接
        }else{
          this.close()
        }
        this.$emit('refresh')
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
            LoginPwd: _this.form.pwd,
            LoginIp: _this.form.ip,
            Name: _this.form.name,
            LoginPort: port,
          }
          if (_this.editId) {
            param.Id = _this.editId;
            serverApi.terminalClientUpdate(param).then(function (res) {
              _this.save_back(res,fn)
            })
          } else {
            serverApi.terminalClientAdd(param).then(function (res) {
              _this.save_back(res,fn)
            })
          }

        } else {
          ElMessage.error('请填写完整');
          return false;
        }
      });
    },
    initWs() {
      setTimeout(()=>{
        if (this.ws) {
        this.ws.close();
        this.ws = null; 
      }
      if(this.term) {
        console.log(this.term)
        this.term.dispose(); 
        this.term.clear();
        this.term = null; 
      }
      this.initTerm()
      
    },300)
     
      
    },
    submitForm() { //链接
      this.$refs['elForm'].validate((valid) => {
        if (valid) {
          this.save(()=>{//保存后跳到链接
            this.activeName = 'second'
            this.secondShow = true;
            this.initWs()
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
    utf8_to_b64(rawString) {
      return btoa(unescape(encodeURIComponent(rawString)));
    },
    b64_to_utf8(encodeString) {
      return decodeURIComponent(escape(atob(encodeString)));
    },
    bytesHuman(bytes, precision) {
      if (!/^([-+])?|(\.\d+)(\d+(\.\d+)?|(\d+\.)|Infinity)$/.test(bytes)) {
        return '-'
      }
      if (bytes === 0) return '0';
      if (typeof precision === 'undefined') precision = 1;
      const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB', 'BB'];
      const num = Math.floor(Math.log(bytes) / Math.log(1024));
      const value = (bytes / Math.pow(1024, Math.floor(num))).toFixed(precision);
      return `${value} ${units[num]}`
    },
    isWsOpen() {
      return this.ws && this.ws.readyState === 1
    },
    initTerm() {
      this.term = new Terminal(this.option)
      this.fitAddon = new FitAddon()
      // this.term.loadAddon(this.fitAddon)
      this.fitAddon.activate(this.term)
      this.term.open(this.$refs.terminal)

      // this.attachAddon = new AttachAddon(this.socket)
      // this.fitAddon = new FitAddon()
      // this.attachAddon.activate(this.terminalElement)
      // this.fitAddon.activate(this.terminalElement)
      // this.fitAddon.fit() // 初始化的时候不要使用fit
      setTimeout(() => {
        this.fitAddon.fit()
        this.initSocket()
        this.onTerminalResize()
        this.onTerminalKeyPress()
      }, 300); // 必须延时处理
    },

    onTerminalKeyPress() {
     
    this.term.onKey(e => {
      if (e.domEvent.ctrlKey && e.domEvent.key === 'c') {
          this.term.clear();
      }
    });
   
      this.term.onData(data => {
        this.ws.send(JSON.stringify({
            Id:this.editId,
            Command: data
          }));
      })
    },
    // resize 相关
    resizeRemoteTerminal() {
      // const { cols, rows } = this.term
      // 调整后端终端大小 使后端与前端终端大小一致
      // this.isWsOpen() && this.ws.send(JSON.stringify({ Id: this.editId, Command: '' }))
    },
    onResize: debounce(function () {
      this.fitAddon.fit()
    }, 500),
    onTerminalResize() {
      window.addEventListener('resize', this.onResize)
      // this.term.onResize(this.resizeRemoteTerminal)
    },
    removeResizeListener() {
      window.removeEventListener('resize', this.onResize)
    },
    // socket
    initSocket() {
      let w_s = 'wss';
      let host = window.location.host;
      if (process.env.NODE_ENV === 'development') {
        w_s = 'ws';
        host = '192.168.1.195:8889'
      }
      this.term && this.term.write('连接中...\r\n')
      const token = `${Session.get('token')}`;
      const socketUrl = `${w_s}://${host}/terminal/ws/ssh?Authorization=${token}`;
      this.ws = new WebSocket(socketUrl, ['webssh'])
      this.onOpenSocket()
      this.onCloseSocket()
      this.onErrorSocket()
      this.onMessageSocket()
     
    },

    // 打开连接
    onOpenSocket() {
      this.ws.onopen = () => {
        // this.term.write("链接成功\r\n");
        const str = JSON.stringify({ Id: this.editId, Command: '' });
        this.ws.send(str);
        this.term && this.term.reset()
        // setTimeout(() => {
        //   this.resizeRemoteTerminal()
        // }, 500)
      }
    },

    // 关闭连接
    onCloseSocket() {
      this.ws.onclose = () => {
        if(this.term){
          this.term.write("未连接， 刷新后重连...\r\n");
            setTimeout(() => {
              this.initSocket();
            }, 3000)
        }
       
      }
    },
    // 连接错误
    onErrorSocket() {
      this.ws.onerror = () => {
        this.term.write('连接失败，请刷新！')
      }
    },
    // 接收消息
    onMessageSocket() {
      this.ws.onmessage = res => {
        const msg = JSON.parse(res.data)
        const term = this.term;
       if((typeof msg) == 'string'){
        if (this.first) {
          this.first = false
          term.reset()
          term.element && term.focus()
          // this.resizeRemoteTerminal()
        }
        const bMsg = this.b64_to_utf8(msg);
        term.write(bMsg)
       }
        
      }
    }
  }
}
//
</script>
<style lang="scss">
.ssh-container {
  height: 100%;
  border-radius: 4px;
  background: rgb(24, 29, 40);
  padding: 0px;
  color: rgb(255, 255, 255);
  overflow: auto;

  .xterm-scroll-area::-webkit-scrollbar-thumb {
    background-color: #b7c4d1;
    /* 滚动条的背景颜色 */
  }
}

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