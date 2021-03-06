
####  如何解决卖超问题

 - 在sql加上判断防止数据边为负数 
 - 数据库加唯一索引防止用户重复购买
 - redis预减库存减少数据库访问　内存标记减少redis访问　请求先入队列缓冲，异步下单，增强用户体验
 - 在项目中利用哈希环来在把用户访问秒杀接口的记录保存在各个节点。防止重复下单秒杀。

####  对象级缓存redis
     redis永久缓存对象减少压力
     redis预减库存减少数据库访
     内存标记方法减少redis访问
#### 订单处理队列rabbitmq
     请求先入队缓冲，异步下单，增强用户体验
     请求出队，生成订单，减少库存
     客户端定时轮询检查是否秒杀成功 
#### 分布式身份验证验证解决
    利用在cookie中保存uid和sign，节点只需要验证token和sign有效性即可。
#### 秒杀安全 -- 安全性设计
     秒杀接口隐藏
     数字公式验证码
     接口防刷限流(通用 注解，拦截器方式)
#### redis的库存如何与数据库的库存保持一致
    redis的数量不是库存,他的作用仅仅只是为了阻挡多余的请求透穿到DB，起到一个保护的作用
    因为秒杀的商品有限，比如10个，让1万个请求区访问DB是没有意义的，因为最多也就只能10个
    请求下单成功，所有这个是一个伪命题，我们是不需要保持一致的
    同时这边也利用了redis的lua脚本功能，来保证扣除库存的事务性。
#### redis 预减成功，DB扣减库存失败怎么办

    -其实我们可以不用太在意，对用户而言，秒杀不中是正常现象，秒杀中才是意外，单个用户秒杀中
    -1.本来就是小概率事件，出现这种情况对于用户而言没有任何影响
    -2.对于商户而言，本来就是为了活动拉流量人气的，卖不完还可以省一部分费用，但是活动还参与了，也就没有了任何影响
    -3.对网站而言，最重要的是体验，只要网站不崩溃，对用户而言没有任何影响
#### redis怎么控制库存数量不为负数
    利用lua脚本
#### 为什么要单独维护一个秒杀结束标志
     -1.前提所有的秒杀相关的接口都要加上活动是否结束的标志，如果结束就直接返回，包括轮寻的接口防止一直轮寻
     -2.管理后台也可以手动的更改这个标志，防止出现活动开始以后就没办法结束这种意外的事件

#### rabbitmq如何做到消息不重复不丢失即使服务器重启
     -1.exchange持久化
     -2.queue持久化
     -3.发送消息设置MessageDeliveryMode.persisent这个也是默认的行为
     -4.手动确认
     （在Rabbitmq文件下有详细说明）
#### redis 分布式锁实现方法
    推荐redis作者提出的redislock算法。在golang实现的redis连接池中，也有很好的实现了Redlock算法。：https://github.com/bsm/redislock。
#### 秒杀类似场景sql的写法注意事项
    1.在秒杀一类的场景里面，因为数据量亿万级所有即使有的有缓存有的时候也是扛不住的，不可避免的透穿到DB
     所有在写一些sql的时候就要注意：
     1.一定要避免全表扫描，如果扫一张大表的数据就会造成慢查询，导致数据的连接池直接塞满,导致事故
     首先考虑在where和order by 设计的列上建立索引
     例如： 1. where 子句中对字段进行 null 值判断 . 
           2. 应尽量避免在 where 子句中使用!=或<>操作符 
           3. 应尽量避免在 where 子句中使用 or 来连接条件
           4. in 和 not in 也要慎用，否则会导致全表扫描( 如果索引 会优先走索引 不会导致全表扫描 
            字段上建了索引后，使用in不会全表扫描，而用not in 会全表扫描 低版本的mysql是两种情况都会全表扫描。
            5.5版本后以修。而且在优化大表连接查询的时候，有一个方法就是将join操作拆分为in查询)
           5. select id from t where name like '%abc%' 或者
           6.select id from t where name like '%abc' 或者
           7. 若要提高效率，可以考虑全文检索。 
           8.而select id from t where name like 'abc%' 才用到索引 慢查询一般在测试环境不容易复现
           9.应尽量避免在 where 子句中对字段进行表达式操作 where num/2  num=100*2
     2.合理的使用索引  索引并不是越多越好，使用不当会造成性能开销
     3.尽量避免大事务操作，提高系统并发能力
     4.尽量避免象客户端返回大量数据，如果返回则要考虑是否需求合理，实在不得已则需要在设计一波了！！！！！