#使用yo发布到Github Pages

首先你需要有一个github的帐号，新建一个仓库（Repo）。

##安装Go语言

如果你不懂请到网上自行找教程，安装Go语言对我这类只会在GNU/Debian
下面用包管理工具笨蛋来说，有些麻烦，如果你要我一步一步较你如何安装
难度太高了。

##安装yo

>go get github.com/newblue/neyo

安装好了之后，请在你的`GOROOT/bin`或者`GOPATH/bin`找一个叫做yo的可执行
文件，如果你找到了，说明你已经安装好这个工具。

##创建一个新的博客项目

>yo new blog

一个静态博客已经基本生成好了。然后你就可以开始写你自己的博客了。

进入`blog`这个目录

> yo post 小明上广州

用你喜欢的编辑器编辑 `posts/小明上广州.md`

##编译

进入blog这个目录

> yo compile

默认情况下，最后结果会在`blog/public`生成一个静态的网站，你可以把这个
静态网站同步的Github Pages。

##同步到Pages

> cd `blog/public`
> git init
> git checkout --orphan gh-pages
> git add *
> git commit

> git remote add origin https://github.com/<username>/<repo>.git
> git push -u origin master

大概就是这样。
