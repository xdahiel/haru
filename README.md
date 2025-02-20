# haru

一个使用 Go 语言开发的全文搜索引擎。你可以使用它构建出属于自己的搜索引擎。

## 特性

* 支持中文搜索。中文分词算法用Go重写[jieba](https://github.com/fxsjy/jieba)分词。
* 并发搜索。使用Goroutine更加轻量、并发度更大。
