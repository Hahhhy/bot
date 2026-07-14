# Day 0 技术详解：消息如何从 QQ 到 Go 程序

- **日期**: 2026-07-14
- **主题**: 逐一解释 Day 0 涉及的每个技术概念

---

## 完整架构全景图

```
你的手机 QQ / 电脑 QQ
       │
       │ (腾讯 NT 私有协议)
       ▼
  ┌─────────────────┐
  │  NapCat         │  ← 无界面的 QQ 客户端
  │  (QQ NT 内核)   │     用自己的 QQ 号登录腾讯服务器
  │                 │     收消息、发消息
  └──────┬──────────┘
         │
         │ WebSocket (ws://0.0.0.0:3001)
         │ 协议: OneBot v11 (JSON 格式)
         ▼
  ┌─────────────────┐
  │  你的 Go 程序    │  ← main.go (qqbot)
  │                 │     读取 JSON → 判断类型 → 构造回复 → 写回
  └─────────────────┘
```

---

## 1. NapCat 是什么？

**NapCat = 没有图形界面的 QQ 客户端。**

腾讯 QQ 客户端和腾讯服务器之间的通信使用的是 NT 协议。正常情况下，你用手机 QQ 或电脑 QQ 登录，在客户端里看消息、发消息。

NapCat 做的事情一样——它内含 QQ NT 核心模块（从 `linuxqq_3.2.29_amd64.deb` 提取），但没有窗口、没有界面，在后台静默运行。你扫码登录后，它就以你的 QQ 小号身份连上腾讯服务器。

**为什么不用官方 QQ 客户端？**

官方的 QQ 客户端不对外暴露编程接口（没有 API）。NapCat 在 QQ NT 核心外面包了一层，把收到的消息**转换成 JSON 格式**，通过 WebSocket 发给你的程序。你的程序想发消息，也只需要发 JSON 给 NapCat，NapCat 帮你转发到腾讯服务器。

**类比**：NapCat 像一个翻译官+传话筒——站在 QQ 服务器和你自己的 Go 程序之间，负责把 QQ 消息翻译成 JSON 给你，把你写的 JSON 翻译成 QQ 消息发给群。

---

## 2. WebSocket 是什么？

**WebSocket 本质上是一个 "一直在线的 TCP 连接"。**

| 对比 | 普通 HTTP | WebSocket |
|------|-----------|-----------|
| 连接方式 | 客户端问 → 服务器答 → 断开 | 建立后一直保持 |
| 实时性 | 客户端不断轮询（问"有新消息吗？"） | 服务器主动推送 |
| 类比 | 发短信 | 打电话 |

Go 代码里建立 WebSocket 连接的位置：

```go
// main.go 第 35 行
conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3001", nil)
```

- `ws://` 协议前缀表示使用 WebSocket
- `localhost` 表示本机
- `3001` 是端口号

连接成功后，NapCat 和你的 Go 程序之间就一直保持这条通道。NapCat 收到 QQ 消息后，立刻通过这个通道推送给你的程序，不需要你的程序反复询问。

---

## 3. 端口 3001 是什么？为什么需要一个端口？

**端口就是操作系统的 "门牌号"。**

你的电脑同时运行着很多网络程序——浏览器、SSH、NapCat、你的 Go 程序等等。操作系统通过 **IP + 端口号** 来区分网络数据包应该送给哪个程序。

- `127.0.0.1:3001` 的含义：
  - `127.0.0.1`（localhost）：本机
  - `3001`：端口号
- 数据包到达本机后，操作系统看端口号是 3001，就知道要交给 NapCat

**类比**：快递只知道你的家庭地址（127.0.0.1），但不知道应该交给你还是你室友。端口号就是 "交给谁" 的明确指令。

**为什么这个端口默认不开启？**

NapCat 安装后默认不开启 WebSocket 服务。需要你明确配置它"在 3001 端口开启 WebSocket 等待连接"。这就是 `onebot11_2381484102.json` 的作用。

---

## 4. JSON 配置文件的作用

`~/.config/napcat-qq-plugin/config/onebot11_2381484102.json`：

```json
"websocketServers": [
  {
    "name": "ws-main",
    "enable": true,
    "host": "0.0.0.0",    ← 监听所有网络接口
    "port": 3001           ← 监听 3001 端口
  }
]
```

这是 **NapCat 的 OneBot 协议配置文件**。它告诉 NapCat：

- "在 3001 端口开启一个 WebSocket 服务器"
- "监听所有 IP 地址（0.0.0.0）"
- "使用 OneBot v11 协议格式"

**何时被读取？** NapCat 启动时读取，**不能热加载**。所以配置必须在启动 NapCat 之前写好——这是之前反复调试失败的关键原因：先启动了 NapCat（读的是旧配置），后改了文件，NapCat 不会重新读取。

---

## 5. OneBot v11 是什么？

**OneBot v11 是一套标准化的 JSON 消息格式。**

不同的人写 QQ 机器人用了不同的桥接程序——有人用 go-cqhttp、有人用 NapCat、有人用 LLOneBot。但他们都遵循同一套消息格式，这就是 OneBot 标准。

这样你的 Go 程序只需要关心 OneBot 格式，不需要关心底层是哪个桥接程序。

### OneBot 下行事件格式（NapCat → 你的程序）

群里有新消息时，NapCat 发送给你的 JSON：

```json
{
  "post_type": "message",
  "message_type": "group",
  "group_id": 1063339847,
  "sender": {
    "user_id": 3558509716,
    "nickname": "D O"
  },
  "message": [
    {
      "type": "text",
      "data": {"text": "哈哈哈哈"}
    }
  ]
}
```

对应 Go 代码中的结构体（main.go 第 13-22 行）：

```go
type Event struct {
    PostType    string `json:"post_type"`     // "message"
    MessageType string `json:"message_type"`  // "group"（群聊）
    GroupID     int64  `json:"group_id"`      // 群号
    Sender      struct {
        UserID   int64  `json:"user_id"`      // 发送者 QQ 号
        Nickname string `json:"nickname"`     // 昵称
    } `json:"sender"`
    Message json.RawMessage `json:"message"`   // 消息内容（可以是文本、图片等）
}
```

### OneBot 上行 API 格式（你的程序 → NapCat）

回复消息时，你的程序构造这样的 JSON 发给 NapCat：

```json
{
  "action": "send_group_msg",
  "params": {
    "group_id": 1063339847,
    "message": [
      {"type": "text", "data": {"text": "你好，我上线了！"}}
    ]
  }
}
```

对应 Go 代码（main.go 第 25-29 行）：

```go
type APIRequest struct {
    Action string      `json:"action"`   // "send_group_msg"
    Params interface{} `json:"params"`   // 群号和消息内容
}
```

### 常见的 OneBot 字段

| 字段 | 含义 | 示例值 |
|------|------|--------|
| `post_type` | 事件类型 | `"message"`（消息）, `"notice"`（通知） |
| `message_type` | 消息类型 | `"group"`（群聊）, `"private"`（私聊） |
| `group_id` | 群号 | `1063339847` |
| `user_id` | QQ 号 | `3558509716` |
| `action` | API 动作 | `"send_group_msg"`（发群消息） |

---

## 6. screen 是什么？

**screen 是 Linux 终端会话管理器。**

```bash
screen -dmS napcat napcat-qq
```

- `screen`：创建新的会话
- `-d`（detach）：启动后立即脱离（后台运行）
- `-m`：强制创建新会话
- `-S napcat`：给会话命名为 "napcat"
- 最后跟要执行的命令 `napcat-qq`

**为什么调试时 screen 不合适？**

NapCat 启动后需要**扫码登录**。在 screen 里运行时你看不到终端输出，也就没法扫码。没有登录的 NapCat 不会开启 OneBot WebSocket 服务，所以 Go 程序连不上 → 报错 `connection refused`。

**正确流程：**

```bash
# 1) 调试时：前台运行，能看二维码
napcat-qq

# 2) 扫码登录，等待看到 "Server Started :::3001"

# 3) 另开终端，启动 Go 程序
cd /home/asus/Dev/QQbot && ./qqbot

# 4) 正式部署时：用 screen 后台运行
screen -dmS napcat napcat-qq
```

---

## 7. 完整消息流转（逐行对应 main.go）

### 建立连接（第 32-40 行）

```go
wsURL := "ws://localhost:3001"
conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 操作系统找到监听 3001 端口的 NapCat，建立 WebSocket 长连接
```

### 循环等待消息（第 52-57 行）

```go
for {
    _, raw, err := conn.ReadMessage()
    // 阻塞等待。没有消息时程序在这里等待。
    // NapCat 把 QQ 消息翻译成 JSON → 通过 conn 推送过来
```

### 解析 JSON（第 59-63 行）

```go
    var event Event
    if err := json.Unmarshal(raw, &event); err != nil {
        // 把 raw 字节流解析到 Event 结构体
        continue  // 解析失败就跳过这条
    }
```

### 判断消息类型（第 66 行）

```go
    if event.PostType == "message" && event.MessageType == "group" {
        // 只处理群聊消息，忽略私聊和系统通知
```

### 构造回复（第 70-78 行）

```go
        reply := APIRequest{
            Action: "send_group_msg",
            Params: map[string]interface{}{
                "group_id": event.GroupID,
                "message": []map[string]interface{}{
                    {"type": "text", "data": map[string]string{"text": "你好，我上线了！"}},
                },
            },
        }
```

### 发送回复（第 79-81 行）

```go
        replyBytes, _ := json.Marshal(reply)
        conn.WriteMessage(websocket.TextMessage, replyBytes)
        // 把 Go 结构体变回 JSON → 通过 WebSocket 发给 NapCat
        // NapCat 调用 QQ API 把消息发到群里
```

---

## 8. 完整链路总结

| 步骤 | 谁 | 做什么 |
|------|-----|--------|
| 1 | 群友 | 在 QQ 群发了一条消息 |
| 2 | 腾讯服务器 | 把消息推送给 NapCat（它也是已登录的 QQ 客户端） |
| 3 | NapCat | 把 QQ 消息翻译成 OneBot v11 JSON |
| 4 | NapCat | 通过 WebSocket 把 JSON 推给你的 Go 程序 |
| 5 | main.go | `conn.ReadMessage()` 收到 JSON 字节 |
| 6 | main.go | `json.Unmarshal` 解析到 Event 结构体 |
| 7 | main.go | 判断 post_type=message 且 group → 群消息 |
| 8 | main.go | 构造 APIRequest（send_group_msg） |
| 9 | main.go | `conn.WriteMessage` 把 JSON 写回 NapCat |
| 10 | NapCat | 执行 QQ API 调用，把消息发到群里 |
| 11 | 群友 | 看到机器人的回复 |

---

## 核心组件速查表

| 组件 | 是什么 | 为什么需要它 |
|------|--------|-------------|
| **QQ 服务器** | 所有消息的中转站 | 消息的终极来源和目的地 |
| **NapCat** | 无界面的 QQ 客户端 | 连接 QQ 服务器，把消息翻译成 JSON |
| **WebSocket (3001端口)** | NapCat 和 Go 程序之间的长连接 | 实时的双向消息通道 |
| **OneBot v11** | 标准化的 JSON 消息格式 | 让 NapCat 和你的程序用同一种语言对话 |
| **onebot11_*.json** | NapCat 的配置文件 | 控制是否开启 WebSocket、在哪个端口 |
| **main.go** | 你的 QQ 机器人程序 | 接收消息 → 业务逻辑 → 回复消息 |