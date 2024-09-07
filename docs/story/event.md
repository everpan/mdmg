# CQRS/ES

> 命令查询职责分离 (Command Query Responsibility Segregation)
> 事件溯源 (Event Sourcing)

** 原则 **

系统并不直接写入、更新数据本身，而是不停地追加事件。
基于事件驱动避免了轮询和广播数据本身，系统内的通信仅仅是时时刻刻发生的事件

## 1. 注册

任何发送的event都必需提前注册；目的是为了规范信息来源、控制版本以及规范schema

| 字段         | 类型     | 描述                                 |
|------------|--------|------------------------------------|
| reg_id     | int    | 注册号                                |
| source     | string | 来源                                 |
| event_type | string | 类型，区别同一个source的不同消息                |
| version    | int    | 从1开始，每次更新+1，同时将历史版本信息存放于历史记录中，便于追溯 |
| describe   | string | 描述，用于管理                            |
| schema     | text   | 消息结构定义，与版本对应，常用于验证消息准确性，可选使用       |

## 2. 生产

业务系统产生消息

| 字段         | 类型   | 描述                                                  |
|------------|------|-----------------------------------------------------|
| ev_id      | int  | 消息序号 ，自增；可以作为offset                                 |
| reg_id     | int  | 注册号，非注册不能记录                                         |
| version    | int  | context与对应的version相匹配，按照对应schema来产生消息，非强制。          |
| context    | text | 消息内容，消息为静态，一旦产生、不允许更新。<br/>通过追加来明确更新，即：最后一条信息表示最终状态 |
| created_at | time | 产生时段                                                |

## 3. 观察者 Observer

当关注的信息到达时刻，将主动调用或者分发出去

| 字段         | 类型      | 描述                                           |
|------------|---------|----------------------------------------------|
| ob_id      | int     | 观察者id                                        |
| source     | string  | 观察来源                                         |
| event_type | string  | 观察来源信息的类型                                    |
| trigger    | enum    | callback、dispatch<br/>通常callback比较耗时，优化策略适用。 |
| context    | string  | 逻辑，v8 code                                   |
| is_active  | tinyint | 是否启用 1-启用，0-停用                               |

## 4. 消费（主动轮询）以api形式


