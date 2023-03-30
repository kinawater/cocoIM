import {w3cwebsocket,IMessageEvent,ICloseEvent} from 'websocket';
import {Buffer} from 'buffer';

export const PingCode : number = 100;
export const PongCode : number = 101;
export const Ping = new Uint8Array([0,PingCode,0,0,0,0])
export const Pong = new Uint8Array([0,PongCode,0,0,0,0])

// 心跳间隔
const heartbeatIntervalTime = 10 //秒
const readDeadLineTimeOutLimit = heartbeatIntervalTime * 3 * 1000 //毫秒
const heartbeatLoopCheckGapTime = heartbeatIntervalTime * 1000 //毫秒
const raadDeadLineLoopCheckGapTime = 1000 //毫秒


// 状态枚举
export enum State {
    INIT,
    CONNECTING,
    CONNECTED,
    RECONNECTING,
    CLOSEING,
    CLOSED,
}
// ACK 枚举
export enum Ack {
    Success = "Success",
    Timeout = "Timeout",
    Loginfailed = "LoginFailed",
    Logined = "Logined",
}

// sleep函数
export let sleep = async (second:number):Promise<void> => {
    return new Promise((resolve,_) => {
        setTimeout(() => {
            resolve()
        },second * 1000)
    })
}

export let doLogin = async (url:string):Promise<{status:string,conn:w3cwebsocket}> => {
    const LoginTimeout = 5 //秒
    return new Promise((resolve,_)=>{
        let conn = new w3cwebsocket(url)
        conn.binaryType = "arraybuffer"

        // 超时登录
        let loginTimeout = setTimeout(()=>{
            resolve({status:Ack.Timeout,conn:conn});
        },LoginTimeout * 1000);

        // 成功
        conn.onopen = () => {
            console.info("websocket link start - open -Status:",conn.readyState)
            if (conn.readyState == w3cwebsocket.OPEN) {
                clearTimeout(loginTimeout)
                resolve({status:Ack.Success,conn:conn})
            }
        }
        conn.onerror = (error:Error) => {
            clearTimeout(loginTimeout)
            console.debug(error)
            resolve({status:Ack.Loginfailed,conn:conn})
        }
    })
}

export class UserClient {
    wsUrl:string
    state = State.INIT
    private userConn : w3cwebsocket | null
    private lastRead : number

    constructor(url:string,user:string) {
        // 从链接接收变量
        this.wsUrl = `${url}?user=${user}`
        this.userConn = null
        this.lastRead = Date.now()
    }
    // 登录方法
    async login():Promise<{status:string}> {
        // 已经链接->已经登录
        if (this.state == State.CONNECTED) {
            return {status:Ack.Logined}
        }
        // 设置状态为正在链接
        this.state = State.CONNECTING

        // 链接
        let status,conn;
        try{
            let result = await doLogin(this.wsUrl);
            status = result.status
            conn = result.conn
            console.info("login access - ",status)
        } catch (error) {
            console.error("login error, please check , err is " , error)
            return {status:Ack.Loginfailed}
        }

        if (status !== Ack.Success) {
            this.state = State.INIT
            return {status}
        }

        conn.onmessage = (msgEvent:IMessageEvent) => {
            try{
                this.lastRead = Date.now()
                let buf = Buffer.from(<ArrayBuffer>msgEvent.data);
                let command = buf.readInt16BE(0)
                let msgLength = buf.readInt32BE(2)
                if (command == PingCode) {
                    console.info("<<<<<<<<<<< pong")
                }
            }catch (error) {
                console.error(msgEvent.data,error)
            }
        }
        conn.onerror = (error) => {
            console.info("websocket link error :",error)
            // TODO 重连 错误处理 errorHandle
        }
        conn.onclose = (msgEvent:ICloseEvent) => {
            console.debug("connect [onclose] is fired")
            if (this.state == State.CLOSEING) {
                // TODO onclose
                return
            }
            //TODO 重连 错误处理 errorHandle
        }
        this.userConn = conn
        this.state = State.CONNECTED
        this.heartbeatLoop()
        // TODO 读取超时 readDeadLineLoop

        return {status}
    }
    // 退出
    logginOut(){
        // 正在关闭，弹出，避免重复关闭
        if (this.state === State.CLOSEING) {
            return
        }
        this.state = State.CLOSEING
        // 空连接
        if (!this.userConn) {
            return
        }
        console.info("user connect closing......")
        this.userConn.close()
    }
    // 心跳
    private heartbeatLoop() {
        console.debug("user connect heartbeat loop check start~")
         let loop = () => {
            if (this.state != State.CONNECTED) {
                console.debug("now connect is not connected,heartbeatLoop exited")
                return
            }
            // TODO send
             console.log("send ping")
             setTimeout(loop,heartbeatLoopCheckGapTime)
         }

        setTimeout(loop,heartbeatLoopCheckGapTime)
    }
    // 读超时
    private readDeadLineLoop() {
        console.debug("deadLine timeOut check Loop start")
        let loop = () => {
            if (this.state != State.CONNECTED) {
                console.debug("now connect is not connected,heartbeatLoop exited")
                return
            }
            if ((Date.now() - this.lastRead) > readDeadLineTimeOutLimit) {
                //TODO 重连 错误处理 errorHandle
            }
            setTimeout(loop,raadDeadLineLoopCheckGapTime)
        }
        setTimeout(loop,raadDeadLineLoopCheckGapTime)
    }
    // 连接关闭
}