// WebSocketService.js
import { Session } from '/@/utils/storage';
class WebSocketService {
  constructor(url,initialMessage,error) {
    this.url = url;
    this.socket = null; // WebSocket 实例
    this.reconnectInterval = null; // 重连定时器
    this.shouldReconnect = true; // 控制是否重连
    this.messageHandlers = []; // 消息处理器
    this.onOpenCallbacks = []; // 连接打开的回调
    this.onCloseCallbacks = []; // 连接关闭的回调
    this.initialMessage = initialMessage; // 初始发送的信息
    this.errorCinnect = error;//创建失败的回调
    this.connect(); // 初始化连接
  }

  connect() {
    const token = `${Session.get('token')}`; 
    if(!token){
      return
    }
    let host = window.location.host;
    if (process.env.NODE_ENV === 'development') {
       host = import.meta.env.VITE_API_WS
       
    }
    // let host = 'fops.fsgit.cc'
    const w_s = 'wss://' + host;
    const ws = `${w_s}/${this.url}?Authorization=${token}`
    this.socket = new WebSocket(ws);
    if(this.socket){
      this.socket.onopen = () => {
        this.clearReconnect(); // 清除重连定时器
        this.shouldReconnect = true; // 允许重连
        // console.log(this.socket.readyState)
        if(this.socket.readyState == 1){ //不判断链接状态 可能会报错
          this.socket.send(this.initialMessage);
        }
        
      };
    }
    
    this.socket.onmessage = (event) => {
      if(event.data){
        const d = JSON.parse(event.data)
        // console.log(d) 
        if(d.StatusCode == 401 || d.StatusCode == 403){
          const msg = d.StatusMessage;
          this.shouldReconnect = false; // 禁止重连
          this.socket.close();
          if(d.StatusCode == 401){
           
          }
          if(d.StatusCode == 403){
            if(msg){
              this.error && this.error(msg)
            }
          }
          return
        }
        this.messageHandlers.forEach(handler => handler(d));
      }
      
    };
    this.socket.onerror = (error) => {
      console.error('WebSocket 错误:', error);
    };
    this.socket.onclose = () => {
      this.onCloseCallbacks.forEach(callback => callback()); // 执行关闭回调
      if (this.shouldReconnect) {
        this.reconnect(); // 开始重连
      }
    };
  }

  reconnect() {
    if (!this.reconnectInterval) {
      this.reconnectInterval = setInterval(() => {
        // console.log('尝试重连...');
        this.connect(); // 重新连接
      }, 3000); // 每 5 秒重连一次
    }
  }

  clearReconnect() {
    if (this.reconnectInterval) {
      clearInterval(this.reconnectInterval);
      this.reconnectInterval = null;
    }
  }

  sendMessage(message) {
    // console.log('发送消息:', message); // 监听发送的消息
    this.initialMessage = message;
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket && this.socket.send(message);
    } else {
      console.error('WebSocket 连接未打开，无法发送消息');
    }
  }
  onMessage(handler) {
    this.messageHandlers.push(handler);
  }

  onOpen(callback) {
    this.onOpenCallbacks.push(callback);
  }

  onClose(callback) {
    this.onCloseCallbacks.push(callback);
  }

  close() {
    this.shouldReconnect = false; // 禁止重连
    if (this.socket) {
      this.socket.close();
    }
    this.clearReconnect(); // 清除重连定时器
  }
}

export default WebSocketService;