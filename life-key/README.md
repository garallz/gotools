# Life Map (有生命周期的Map)

## Use Way And Order (使用方法和顺序）

- `First:` Initialize the definition struct (初始化定义结构)
``` go
var a = lifekey.LifeData{}
```
- `Second:` Start timing cleanup (启动定时清理)
``` go
a.GcData(time.Second)
// You can adjust the cleaning cycle according to the situation.
// 你可以根据情况需求调节定时检查清理周期
```
- `Third:` Call the function method (调用函数方法)
``` go
a.Set(key, life_time)
a.SetAddData(data, key, life_time)          // If data was null, input nil
a.Get(key)
a.Check(key)
a.UpdateData(data, key)         // If data was null, input nil
a.Delete(key)
```

## Precautions (注意事项)
  1. When you use it package, please read the execution order carefully.
  (当你使用该包时，请认真阅读执行顺序）
  2. The package is use second as a lifecycle unit
  (该包用秒做为生命周期单位)
