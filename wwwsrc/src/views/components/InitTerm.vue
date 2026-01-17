<template>
    <div class="ssh-container" ref="terminal"></div>
</template>
<script>
import { Session } from '/@/utils/storage';
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { debounce } from 'lodash'
export default {  //终端
    name: 'InitTerm',
    data() {
        return {
            timeClose: null,
            errorNum:0,//重连次数
            againFlag:true,//断开后是否重新链接
            fireLoginIp:'',//记录ip
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
            },
            sshByLogin:false, //走第二个接口
            initRow:{
                LoginIp: '',
                LoginName: '',
                LoginPwd: '',
                LoginPort: 22
            }
        }
    },
    mounted(){
        // window.addEventListener('keydown', this.handleKeyDown);
    },
    beforeDestroy() {
        // console.log('beforeDestroy')
        // window.removeEventListener('keydown', this.handleKeyDown);
        this.clearWs()
    },
    methods: {
        initStart(row){ //走第二个接口
            this.sshByLogin = true;
            this.againFlag = false; //不在重新链接
            this.initRow = {...row}
            // console.log('row',row.LoginIp)
            this.init(row.LoginIp)
           
        },
        clearInit(loginIp){ //断开 - 重连 
            this.sshByLogin = false;
            this.againFlag = false; //不在重新链接
            this.init(loginIp)
        },
        clearWs(fn) { //清除
            this.ws && this.ws.close();
            this.ws = null;
            this.term && this.term.dispose();
            this.term = null;
            clearInterval(this.timeClose)
            this.timeClose = null;
            this.againFlag = false; //不在重新链接
            this.errorNum = 0;
            fn && fn()
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
        inits(loginIp) { //弹框刚打开的时候调用
            this.sshByLogin = false;
            this.againFlag = true; //
            this.init(loginIp) 
        },
        init(loginIp) {
            this.clearWs(()=>{
                this.initFire(loginIp)
            })    
        },
        get_row(loginIp){
            let row = {};
            if(this.sshByLogin){
                let initRow = this.initRow;
                initRow.LoginPort = initRow.LoginPort * 1;
                row = {...initRow}
            }else{
                row.LoginIp = loginIp;
            }
            return row
        },
        initFire(loginIp){
            this.fireLoginIp = loginIp;
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
            this.term && this.term.write('连接中...\r\n')
            setTimeout(() => {
                this.fitAddon.fit()
                this.initSocket(loginIp)
                this.onTerminalResize(loginIp)
                this.onTerminalKeyPress(loginIp)
            }, 300); // 必须延时处理
        },
        handleKeyDown(event){
            // if (event.ctrlKey && event.key == 'k') {
            //    if(this.term){
            //     const row = this.get_row(this.fireLoginIp);
            //     this.term && this.ws && this.ws.send(JSON.stringify({
            //         ...row,
            //         Command: '\f'
            //     }));
            //    }
                
            //     event.preventDefault(); // 阻止默认行为
            // }
        },
        onTerminalKeyPress(loginIp) {
            const row = this.get_row(loginIp);
            // console.log('onTerminalKeyPress',row)
            this.term && this.term.onKey(e => {
                if (e.domEvent.metaKey && e.domEvent.key === 'k') {
                //   metaKey 
                     this.ws.send(JSON.stringify({
                        ...row,
                        Command: '\f'
                    }));
                }
            });
           
            this.term.onData(data => {
                
                this.ws.send(JSON.stringify({
                    ...row,
                    Command: data
                }));
            })
        },
        // resize 相关
        resizeRemoteTerminal(loginIp) {
            // const { cols, rows } = this.term
            // 调整后端终端大小 使后端与前端终端大小一致
            // this.isWsOpen() && this.ws.send(JSON.stringify({ LoginIp: loginIp, Command: '' , ...this.initRow,}))
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
        initSocket(loginIp) {
            let host = 'wss://' + window.location.host+ '/';
            if (process.env.NODE_ENV === 'development') {
               host = import.meta.env.VITE_API_WS
            }
            let w_s = host;
            const token = `${Session.get('token')}`; //terminal/ws/sshByLogin
            const ssh = this.sshByLogin?'sshByLogin':'ssh'
            const socketUrl = `${w_s}terminal/ws/${ssh}?Authorization=${token}`;
            this.ws = new WebSocket(socketUrl, ['webssh'])
            this.term && this.term.clear()
            this.term && this.term.write('连接中...\r\n')
            this.onOpenSocket(loginIp)
            this.onCloseSocket(loginIp)
            this.onErrorSocket(loginIp)
            this.onMessageSocket(loginIp)

        },

        // 打开连接
        onOpenSocket(loginIp) {
            this.againFlag = true;//断开后重新链接
            this.ws.onopen = () => {
                this.term && this.term.clear()
                const row = this.get_row(loginIp);
                const str = JSON.stringify({
                    ...row,
                     Command: '' });
                this.ws.send(str);
                // this.term && this.term.reset()
                // setTimeout(() => {
                //   this.resizeRemoteTerminal()
                // }, 500)
            }
        },

        // 关闭连接
        onCloseSocket(loginIp) {
            this.ws.onclose = () => {
                clearInterval(this.timeClose)
                if(this.againFlag){
                    this.errorNum ++
                    this.term && this.term.write(`第${this.errorNum}次重新链接...\r\n`)
                    this.timeClose = setTimeout(() => {
                        if (this.term) { this.initSocket(loginIp); }
                    }, 2000);
                }
                
            }
        },
        // 连接错误
        onErrorSocket(loginIp) {
            this.ws.onerror = () => {
                this.term && this.term.write('连接失败，请刷新！')
            }
        },
        // 接收消息
        onMessageSocket(loginIp) {
            this.ws.onmessage = res => {
                const msg = JSON.parse(res.data)
                const term = this.term;
                if ((typeof msg) == 'string') {
                    this.errorNum = 0;
                    if (this.first) {
                        this.first = false
                        term.reset()
                        term.element && term.focus()
                        // this.resizeRemoteTerminal()
                    }
                    const bMsg = this.b64_to_utf8(msg);
                    term.write(bMsg)
                }else{
                    
                    if(msg.StatusMessage){
                        term.clear()
                        term.write(msg.StatusMessage)
                    }
                    
                }

            }
        }
    }
}
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
</style>