import Cookies from 'js-cookie'

var websock = null;

var global_callback = null;

var wsuri = "ws://172.18.10.46:8929/echo"

function createWebSocket(callback) {
    if (websock == null || typeof websock !== WebSocket) {
        initWebSocket(callback);
    }
}

function initWebSocket(callback) {
    global_callback = callback;
    // 初始化websocket
    websock = new WebSocket(wsuri);
    websock.onmessage = function (e) {
        websocketonmessage(e);
    };
    websock.onclose = function (e) {
        websocketclose(e);
    };
    websock.onopen = function () {
        websocketOpen();
    };
    // 连接发生错误的回调方法
    websock.onerror = function () {
        console.log("WebSocket连接发生错误");
        websock.onclose = function (e) {
            console.log("关闭连接");
            websocketclose(e);
        };
        // console.log("重新连接");
        // createWebSocket();
    };
}

// 实际调用的方法
function sendSock(agentData) {
    if (websock.readyState === websock.OPEN) {
        // 若是ws开启状态
        websocketsend(agentData);
    } else if (websock.readyState === websock.CONNECTING) {
        // 若是正在开启状态，则等待1s后重新调用
        setTimeout(function () {
            sendSock(agentData);
        }, 1000);
    } else {
        // 若未开启 ，则等待1s后重新调用
        setTimeout(function () {
            sendSock(agentData);
        }, 1000);
    }
}

function closeSock() {
    websock.close();
}

// 数据接收
function websocketonmessage(msg) {
    let Uuid = Cookies.get('Uuid')
    console.log(msg);
    let result = JSON.parse(msg.data);
    if(String(result.Uuid).length > 0 && JSON.parse(result.Uuid).constructor == Array){
        JSON.parse(result.Uuid).forEach(e =>{
            if(Uuid == e){
                global_callback(result);
            }
        })
    }
}

// 数据发送
function websocketsend(agentData) {
    console.log("发送数据：" + JSON.stringify(agentData));
    websock.send(JSON.stringify(agentData));
}

// 关闭
function websocketclose(e) {
    console.log("连接关闭",JSON.stringify(e));
}

function websocketOpen(e) {
    console.log("连接打开");
}

export { sendSock, createWebSocket, closeSock };