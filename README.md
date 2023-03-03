# hhttp一个为了方便而生的Go语言http client
关于为什么是h开头？

- 可以是我的姓氏（侯）的第一个字母；
- 可以是我们学变成的第一个程序hello ,world的第一个字母；
- 也可以是hello、hi；
- ……

# 1 介绍

该库对Go原生http库做了封装，添加了错误处理，超时处理及重试功能。为了实现更快捷、更方便的使用Go的http client发送请求。

推荐使用Go1.18及以上版本。

# 2 主要功能

- [x] GET、POST、DELETE、PUT、PATCH支持；
- [x] 错误处理、错误回调；
- [x] 超时处理；
- [x] 重试功能；
- [x] 自定义header和cookie；
- [ ] 友好的支持文件上传和下载；
- [ ] ……

# 3 基本使用

## 3.1 全局参数设置

```go
hhttp.Load(
    hhttp.WithTimeout(time.Second), // 设置全局超时时间
    hhttp.WithRetryCount(1), // 设置全局重试次数
    hhttp.WithRetryWait(time.Second),// 设置全局重试等待时间
    // 设置全局重试错误时的回调方法
    hhttp.WithRetryError(func(ctx context.Context, r hhttp.Request) error {
        return nil
    }),
)
```



## 3.2 GET示例

```go
// 正常发送一个Get请求
res, err := hhttp.New().Get(context.Background(), "https://www.houzhenkai.com")
if err != nil {
    fmt.Println(err)
}

// 添加header和cookie
hhttp.New().SetHeader("token", "1234567890").SetCookies(&http.Cookie{
            Name:  "token",
            Value: "abcdefg",
    }).Get(context.Background(), "https://www.houzhenkai.com")


// 正常发送一个Get请求,追加get参数，以key-value格式
res, err := hhttp.New().Get(context.Background(), "https://www.houzhenkai.com",hhttp.NewKVParam("name", "jankin"))
if err != nil {
    fmt.Println(err)
}

// 正常发送一个Get请求,追加get参数，以map格式
res, err := hhttp.New().Get(context.Background(), "https://www.houzhenkai.com",hhttp.NewMapParams(map[string]interface{}{
        "name": "jankin",
    }))
if err != nil {
    fmt.Println(err)
}

// 正常发送一个Get请求,追加get参数，以字符串格式
res, err := hhttp.New().Get(context.Background(), "https://www.houzhenkai.com",hhttp.NewQueryParam("name=jankin"))
if err != nil {
    fmt.Println(err)
}
```



## 3.3 POST 示例

```go
// -- application/x-www-form-urlencoded -- // 
// 以map的形式添加post参数
hhttp.New().Post(context.Background(), "https://www.houzhenkai.com/test/login",hhttp.NewWWWFormPayload(map[string]interface{}{
        "username": "jankin",
    }))


// -- application/json -- //
hhttp.New().SetHeader(hhttp.SerializationType,hhttp.SerializationTypeJSON).Post(context.Background(), "https://www.houzhenkai.com/test/login", hhttp.NewJSONPayload("username=jankin"))

```

## 3.4 超时处理
hhttp共有两种方式可以实现超时处理，具体说明和使用方式如下所示
1. hhttp.WithTimeout(time.Second)
```Go 
res, err := New(WithTimeout(time.Second)).Get(ctx, urlStr)
```
2. 在Get、Post等方法里传入一个带有timeout的context
```Go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
res, err := New(WithTimeout(6*time.Second)).Get(ctx, urlStr)
```
​    注意：第二种方式的优先级将优于第一种，也就是说如果两种方式同时在使用，则以传入ctx的时间为实际超时时间。

## 3.5请求重试
hhttp集成了请求重试，如果在请求失败(含请求超时)后，可以进行请求重试，即在延时一段时间后(可以是0秒)，重新发起请求。
```Go
urlStr := "https://www.houzhenkai.com"
// 请求失败后会再次进行重试请求
ctx := context.Background()
res1, err := New(WithRetryCount(2)).Get(ctx, urlStr)
if err != nil {
    t.Error(err)
}
t.Log(string(res1))
```



# 4 支持

在使用过程中有任何问题，欢迎在[Issues](https://github.com/JankinHou/hhttp/issues)里留言互动， 也欢迎您能指出项目存在的问题。
