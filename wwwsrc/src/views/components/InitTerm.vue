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
            VITE_WS: import.meta.env.VITE_WS,
            timeClose: null,
            errorNum:0,//重连次数
            againFlag:true,//断开后是否重新链接
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
    beforeDestroy() {
        // console.log('beforeDestroy')
        this.clearWs()
    },
    methods: {
        clearInit(id){ //断开 - 重连
            this.againFlag = false; //不在重新链接
            this.clearWs(()=>{
                this.init(id)
            })
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
        init(id) {
            this.clearWs(()=>{
                this.initFire(id)
            })    
        },
        initFire(id){
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
                this.initSocket(id)
                this.onTerminalResize(id)
                this.onTerminalKeyPress(id)
            }, 300); // 必须延时处理
        },
        onTerminalKeyPress(id) {

            this.term && this.term.onKey(e => {
                if (e.domEvent.ctrlKey && e.domEvent.key === 'c') {
                    this.term.clear();
                }
            });

            this.term.onData(data => {
                this.ws.send(JSON.stringify({
                    Id: id,
                    Command: data
                }));
            })
        },
        // resize 相关
        resizeRemoteTerminal(id) {
            // const { cols, rows } = this.term
            // 调整后端终端大小 使后端与前端终端大小一致
            // this.isWsOpen() && this.ws.send(JSON.stringify({ Id: id, Command: '' }))
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
        initSocket(id) {
            let host = window.location.host;
            let w_s = 'wss://' + host + '/';
            if (process.env.NODE_ENV === 'development') {
                w_s = this.VITE_WS;
            }
            const token = `${Session.get('token')}`;
            const socketUrl = `${w_s}terminal/ws/ssh?Authorization=${token}`;
            this.ws = new WebSocket(socketUrl, ['webssh'])
            this.term && this.term.clear()
            this.term && this.term.write('连接中...\r\n')
            this.onOpenSocket(id)
            this.onCloseSocket(id)
            this.onErrorSocket(id)
            this.onMessageSocket(id)

        },

        // 打开连接
        onOpenSocket(id) {
            this.againFlag = true;//断开后重新链接
            this.ws.onopen = () => {
                
                this.term && this.term.clear()
                const str = JSON.stringify({ Id: id, Command: '' });
                this.ws.send(str);
                // this.term && this.term.reset()
                // setTimeout(() => {
                //   this.resizeRemoteTerminal()
                // }, 500)
            }
        },

        // 关闭连接
        onCloseSocket(id) {
            this.ws.onclose = () => {
                clearInterval(this.timeClose)
                if(this.againFlag){
                    this.errorNum ++
                    this.term.clear()
                    this.term && this.term.write(`第${this.errorNum}次重新链接...\r\n`)
                    this.timeClose = setTimeout(() => {
                        if (this.term) { this.initSocket(id); }
                    }, 2000);
                }
                
            }
        },
        // 连接错误
        onErrorSocket(id) {
            this.ws.onerror = () => {
                this.term && this.term.write('连接失败，请刷新！')
            }
        },
        // 接收消息
        onMessageSocket(id) {
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