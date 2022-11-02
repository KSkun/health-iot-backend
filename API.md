# Health IoT Platform 后端 API 文档

## 通用规则

### 数据交换格式

后端均使用常规的 RESTful API 与前端/客户端通信，一般而言，API 请求的定义满足以下规则

- 使用的 HTTP 请求方法（Method）包括：GET、POST、PUT、DELETE
- GET、DELETE 类型的请求没有请求体，参数传递均使用 URL 参数
- POST、PUT 类型的请求有请求体，参数均通过请求体传递，需要指定 Content-Type 头为 application/json
- 需要鉴权的接口会以 \* 符号标出，鉴权方法参见后文

API 响应的定义满足以下规则

- 响应由状态码和响应体两部分内容构成
- 状态码的含义如下，前端可根据状态码判断错误类型
  - 2xx 状态码为正常情况，一般使用 200 OK
  - 4xx 状态码为请求错误
  - 5xx 状态码为服务器错误
- 响应体的格式为 JSON，示例如下，前端解析时可以直接解析响应中的 data 字段，错误的响应均可通过判断状态码过滤

```json
{
  "data": { ... }, // 如有需返回的数据，data 字段会包含数据
  "success": true, // 说明请求是否成功
  "message": "",   // 如果请求不成功，此处会包含错误提示信息
  "error": ""      // （仅调试）如果请求不成功，此处会包含错误信息，生产环境始终为空
}
```

### 鉴权方式

所有需要鉴权的接口均使用 JWT 作为鉴权方式，请求这类接口时，需要指定 Authorization 头，以 Bearer token 格式放置 token，示例如下

```
Authorization: Bearer blablablablablablablablablabla
```

## 接口定义

定义中的响应示例均为成功的响应示例中 data 字段的内容

### 用户相关

#### POST /api/v1/user

用户注册

请求：

```json
{
  "name": "test",      // 用户名
  "password": "111111" // 密码，明文
}
```

响应：

```json
{
  "id": "6362151385681f7e619fdbed" // 新注册的用户ID
}
```

#### GET /api/v1/user/token

用户登录

请求：

URL 参数 `?name=test&password=111111`

响应：

```json
{
  "expire_time": 1667458793951, // token过期时间，精确到毫秒的Unix时间戳
  "token": "JWT token"          // token文本
  }
```

#### *GET /api/v1/user

查询当前用户信息

请求：无

响应：

```json
{
  "id": "635362dea5ac2cd658b3966f", // 用户ID
  "name": "test"                    // 用户名
}
```

### 设备和数据相关

#### *POST /api/v1/device

添加设备

请求：

```json
{
  "name": "test-device1", // 设备名
  "serial": "123457"      // 设备序列号
}
```

响应：

```json
{
  "id": "6362185685681f7e619fdbee" // 新添加的设备ID
}
```

#### *GET /api/v1/device/list

查询当前用户的设备列表

请求：无

响应：

```json
{
  "devices": [
    {
      "id": "63536310a5ac2cd658b39670", // 设备ID
      "name": "test-device",            // 设备名
      "serial": "123456",               // 设备序列号
      "online": false,                  // 当前是否在线
      "battery": 10,                    // 剩余电路百分比*100
      "warning": true                   // 是否有警告
    }
  ]
}
```

#### *GET /api/v1/device/:id

查询某个设备的信息，URL 中的 :id 为设备 ID

请求：URL 中的 :id 换为设备 ID

响应：

```json
{
  "device": {
    "id": "63536310a5ac2cd658b39670",       // 设备ID
    "name": "test-device",                  // 设备名
    "serial": "123456",                     // 设备序列号
    "owner_id": "635362dea5ac2cd658b3966f", // 注册设备的用户ID
    "last_report_time": 1666882009189,      // 最后一次上报时间
    "status": {
      "battery": 10,                        // 剩余电量百分比
      "locating": true,                     // 定位是否正常
      "wearing": true                       // 穿戴是否正常
    },
    "sensor": {
      "heart_rate": 170,                    // 心率每分钟
      "blood_oxygen": 0,                    // 血氧饱和度百分比
      "longitude": 116.470098,              // 经度
      "latitude": 39.992838,                // 纬度
      "sos_warning": false,                 // 是否有求救警报
      "fall_warning": false                 // 是否有摔倒警报
    }
  },
  "online": false,                          // 是否在线
  "warnings": [                             // 警报详细内容
    {
      "field": "battery",                   // 警报数据项
      "type": 2,                            // 警报类型，1为无类型，2为值太低，3为值太高
      "message": "low battery"              // 警报消息
    }
  ]
}
```

#### *GET /api/v1/device/data

查询数据上报记录

请求：可选 URL 参数 `?device=63536310a5ac2cd658b39670`

响应：

```json
{
  "reports": [
    {
      "report": { // 数据上报信息，含义同上
        "id": "635a98751727c739df54f08e",
        "device_id": "63536310a5ac2cd658b39670",
        "time": 1666881653759,
        "status": {
          "battery": 90,
          "locating": true,
          "wearing": true
        },
        "sensor": {
          "heart_rate": 70,
          "latitude": 39.99283981323242,
          "longitude": 116.47010040283203
        }
      },
      "warnings": [] // 该条数据的警告信息，含义同上
    }
  ]
}
```

#### POST /api/v1/device/data

（仅供设备使用）上报数据

请求：

```json
{
  "serial": "123456",        // 设备序列号
  "status": {                // 以下项目同上一个API的含义
    "battery": 10,           // 必填
    "locating": true,        // 必填
    "wearing": true          // 必填
  },
  "sensor": {                // 以下字段选填，填了哪些就更新哪些
    "heart_rate": 170,
    "longitude": 114.420207,
    "latitude": 30.515331
  }
}
```