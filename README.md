# push系统

这套系统能解决什么问题：
1. 给符合指定条件的用户发送业务推广链接
2. 之前是通过php脚本实现，没有界面操作，也不提供接口，不方便重复操作，每次变更（推送文案变更、推送链接变更、获取用户条件变更等），都需要直接修改脚本，不方便也不安全
3. 之前获取的用户是存储在redis里，合并过滤用户都相对麻烦，也不安全，同时也不方便管理
4. 性能问题，之前推送很慢，如果想并发推，还得手动启多个脚本同时运行，操作极不方便，也不安全

这套系统的最初设想：
1. 有图形化界面可以操作，配置信息都落地到db里，管理简单，安全
2. 性能要足够强劲，获取指定条件用户必须可以并发执行，推送时也能并发执行
3. 要安全，不能出现重复推送的情况，这种互斥锁机制实现起来要足够简单
4. 支持更多的状态控制，可以随意启停，可以实时监控等

这套系统现在长什么样？进展到什么程度？
1. 这套系统大体分为两部分，一个是获取及推送用户服务，一个是操作界面（操作界面目前借用的通用op系统: https://github.com/simplejia/op）
2. 已经在线上长期使用，功能稳定，性能卓越

安装使用：
1. 下载源代码
> go get github.com/simplejia/push
2. 配置数据库
> 目前的配置信息存储在mongo db，需要修改配置文件：mongo/push.json
3. 使用
> 进入push目录，启动编译好的push程序，比如：./push -env dev
> 使用通用op系统配置并管理： https://github.com/simplejia/op

注意：
1. “获取用户”任务这一部分，目前开源代码只是摘选了一种实现：见：service/constraint/worker.go，在函数：CreateGetMembersFunc内部，只实现了一种ConditionKindConcrete的类型，这种类型只推送到指定的用户列表，其他类型可以根据自己的项目情况自行扩充
2. “推送用户”任务这一部分，目前开源代码牵涉到具体推送用户这一步给省略了，见：service/project/worker.go，在函数：Push内部，请根据自己的项目情况自行扩充。
3. 这套push系统开源的最主要目的，是为了技术分享，如果大家正好也有类似需求，不妨尝试一下，这套系统完全借助db实现了并发控制，做到了安全，简单，方便等业务目的。


内部实现：（粗略）
1. 以http server方式对外提供服务
2. 内部以任务方式分为两部分，一部分是“获取用户”任务，一部分是“推送用户”任务（下面如果没有特殊说明，均以“获取用户”任务做说明，因为大部分情况下，这两部分在技术实现上是类似的）
3. “获取用户”和“推送用户”这两种任务，内部实现上，均以并发方式（多worker）执行，任务一共包括7种状态，状态流转图如下：

￼
![1](https://github.com/simplejia/nothing/raw/master/push_1.png)


* ready：任务初始状态
* started：任务启动状态，通过接口手动设置，所有其他状态均可以流转成started状态，当重新设置成started状态时，所有参数会被重置，注意，已经是started状态的不可以重新设置成started状态，这样做是防止参数重置造成异常情况发生
* finished：任务正常结束状态
* failed：任务执行失败状态，比如：a：程序异常错误 b：运行错误数超过一定阈值（db调用失败等）  
* stopped：任务被手动停止
* pause：任务被手动暂停，注意，当通过专门的恢复执行接口重新设置成started状态时，所有参数不会被重置，会从上一次任务执行结果处恢复执行
* reboot：这个状态相对复杂，是任务执行过程中触发特定条件后自动转换的状态，比如：有worker挂掉（心跳数据不更新），任务转换为reboot状态，然后等到同一任务下所有机器上的所有worker均保持一段较长时间（1分钟）没有更新心跳数据，任务自动转换为started状态，之所以这样设计，是为了解决某些特定情况下，比如其中一台机器服务重启，worker能够重新分配，保证每一台机器上的worker数量基本均衡


4. 任务以并发方式执行：单master+多worker方式，如下：
> master（一个）：
>> master会定期检测处于started状态的任务，根据每一个任务配置的worker数量，慢慢（不是快速，防止多机器worker分配不均）启动worker，master和worker之间建立通信通道，所有的worker通过唯一的一个通信通道给master发送消息（目前主要是状态变更消息），每一个worker均会产生一个自己的通信通道，master通过这个通道给worker发送消息（目前主要是状态变更消息）。

> worker（多个）：
>> worker之间借助db建立数据通信，可以满足现有场景下所有原子性操作，比如每一个worker只关注自己要扫描的用户表、每一个worker可以原子记数、每一个worker可以原子更新每一个用户推送状态等

## 依赖
    wsp: github.com/simplejia/wsp
    clog: github.com/simplejia/clog
    utils: github.com/simplejia/utils
    mongo: gopkg.in/mgo.v2

## 注意
> 如果在controller里修改了路由，编译前需执行go generate，实际是运行了wsp这个工具，所以需要提前go get github.com/simplejia/wsp
