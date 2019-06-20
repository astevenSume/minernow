
# 运行
```shell
errorgenerater --config ./otc.json --out-dir ../../otc/controllers/error/

errorgenerater --config ./eusd.json --out-dir ../../eusd/eusdplus/

```



# 错误码描述

200 成功

## 服务端内部错误（客户端不使用，直接提示“未知错误”）
（0, 199）


## 公共错误

（201, 400)


## EUSD 错误
(2000, 3000)
