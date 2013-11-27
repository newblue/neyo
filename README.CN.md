* [中文介绍](#chinese-introduction)
    * [安装](#installation-安装)
    * [快速入门](#quick-start-快速入门)

#中文介绍

## yo -- Golang编写的静态博客引擎（派生自`gor`）

`yo`是使用 [Go](http://golang.org/) 实现的类 Ruhoh 静态博客引擎（Ruhoh like），基本兼容 ruhoh 1.x 规范。
相当于与 ruhoh 的官方实现（ ruby 实现），派生自[gor](http://github.com/wendal/gor.git)有以下优点：

1. 速度完胜 -- 编译 wendal.net 近200篇博客,仅需要1秒
2. 安装简单 -- 得益于 golang 的特性，编译后仅一个可运行程序，无依赖

## Installation 安装
====================

> go get github.com/newblue/neyo/yo

**在 Mac下使用 brew 的用户**

如果是通过 [brew](https://github.com/mxcl/homebrew) 来安装`go`，并且没有设置`$GOROOT`跟`$GOPATH`的话，请使用如下命令（路径请更改为自己对应的 golang 的版本信息）

> ln -s /usr/local/Cellar/go/1.0.3/bin/yo /usr/local/bin

## Quick Start 快速入门
======================

新建站点
--------

> yo new example.com

执行完毕后, 会生成example.com文件夹，包含基本素材及演示文章

新建单篇博客
-----------

> cd example.com
> yo post "goodday" [dir/to/img/files]

即可生成 post/goodday.md文件，打开你的markdown编辑器即可编写

如果输入可选参数 `[dir/to/img/files]`，gor 会从该目录拷贝图片文件到配置的目录，同时在 `goodday.md` 中自动插入对应的 `<img>` 标签。

基本配置
--------

打开站点根目录下的`site.yml`文件

1. 填入 title，作者等信息
2. 填入邮箱等信息

打开站点根目录下的 config.yml 文件

1. `production_url`：为你的网站地址，例如 [http://www.makechan.com/](http://www.makechan.com/)最后面不需要加入`/`，生成`rss.xml`等文件时会用到
2. `summary_lines`：首页的文章摘要的长度,按你喜欢的呗
3. `latest`：首页显示多少文章
4. `imgs`：自动插入`<img>`的相关配置
   * `imgtag`：要插入的 <img> 标签的基本格式，`%s` 部分会被自动替换为 `urlperfix/post_name/img_file_name` 的格式
   * `urlperfix`：图片地址前缀
   * `localdir`：图片文件在博客内的本地存放目录

打开`widgets`目录, 可以看到基本的挂件，里面有`config.yml`配置文件

1. `analytics`：暂时只支持`google analytics`，填入`tracking_id`
2. `comments`：暂时只支持`disqus`，请填入`short_name`
3. `google_prettify`：代码高亮,一般不修改


编译生成静态网页
--------------

> yo compile

瞬间完成，生成 compiled 文件夹，包含站点所有资源

本地预览
-------

> yo http

打开你的浏览器，访问 http://127.0.0.1:8080

部署
-----

你可以使用[github pages](http://pages.github.com/)等服务，或者放到你的自己的`vps`下，因为是纯静态文件,不需要`php/mysql/java`等环境的支持

