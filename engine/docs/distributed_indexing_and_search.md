分布式索引和搜索
===

分布式搜索的原理如下：

当文档数量较多无法在一台机器内存中索引时，可以将文档按照文本内容的hash值裂分(sharding)，不同块交由不同服务器索引。在查找时同一请求分发到所有裂分服务器上，然后将所有服务器返回的结果归并重排序作为最终搜索结果输出。

为了保证裂分的均匀性，建议使用Go语言实现的Murmur3 hash函数:

https://github.com/huichen/murmur

按照上面的原理很容易用悟空引擎实现分布式搜索（每个裂分服务器运行一个悟空引擎），但这样的分布式系统多数是高度定制的，比如任务的调度依赖于分布式环境，有时需要添加额外层的服务器以均衡负载，因此就不在这里实现了。