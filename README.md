## Table of Contents
* [English Introduction](#english-introduction)
    * [Installation](#installation)
    * [Quick Start](#quick-start)
* [Copyright and License](#copyright-and-license)

[Chinese README.md](README.CN.md)

# English Introduction

This project is fork from http://github.com/wendal/gor.git

## yo -- A static websites and blog generator engine written in Go fork from gor

Transform your plain text into static websites and blogs.
`yo(fork gor)` is a [Ruhoh](http://ruhoh.com/) like websites and blog generator engine written in [Go](http://golang.org/). It's almost compatible to ruhoh 1.x specification. You can treat yo as a replacement of the official implementation what is written in [Ruby](http://www.ruby-lang.org/en/).

Why reinvent a wheel? `yo` has following awesome benefits:

1. Speed -- Less than 1 second when compiling all my near 200 blogs on wendal.net
2. Simple -- Only one single executable file generated after compiling, no other dependence

## Installation
====================
To install:

> go get github.com/newblue/neyo/yo

If you use [brew](https://github.com/mxcl/homebrew) on Mac, and you didn't set `$GOROOT` and `$GOPATH` environment variable
Please using this command:

> ln -s /usr/local/Cellar/go/1.0.3/bin/yo /usr/local/bin

## Quick Start
======================

Create a new website
--------------------

> yo new example.com
> cd example.com

After execution, a folder named example.com will be generated, including a scaffold & some sample posts.

Create a new post
----------

> yo post "goodday" [dir/to/img/files]

Generate a new post file: post/goodday.md, open it with your markdown editor to write.

`[dir/to/img/files]` is optionl. If it's provided, all files in that dir will be copy into blog dir(configurable dir), and insert `<img>` tag into post file.

Configuration
-------------

Open the `site.yml` file in root folder

1. Input title, author etc.
2. Input email etc.

Open the config.yml file in root folder

1. `production_url` is your website address, such as `http://wendal.net`, don't add `'/'` at last, it will be used to generate `rss.xml` etc.
2. `summary_lines` is the length of abstract on homepage, any number as you like.
3. `latest` is how many posts will be shown on homepage
4. `imgs` parts is auto img config
   * `imgtag`：basic format for <img> tag to be insert. the `%s` part will to replaced by `urlperfix/post_name/img_file_name`
   * `urlperfix`：img file url perfix
   * `localdir`：location inside blog repo for img file storage

Open `widgets` folder, you can see some widgets here, there is a `config.yml` file of each widget for configuration.

1. `analytics` only support `google analytics` by now, please input `tracking_id` here
2. `comments` only support `disqus` by now, please input your `short_name` of disqus here
3. `google_prettify` for code highlighting, normally it's not necessary to change

Compile to generate static web page
--------------

> yo compile

Finished instantly. A new folder named public will be generated, all website is in it.

Local preview
-------
gor also comes with a built-in development server that will allow you to preview what the generated site will look like in your browser locally.

> yo http

Open your favorite web browser and visit: http://127.0.0.1:8080

Deployment
-----

You can deploy it to [GitHub Pages](http://pages.github.com/), or put it to your own `VPS`, because there are only static files(HTML, CSS, js etc.), no need of `php/mysql/java` etc.

Copyright and License
----------------------

This project is licensed under the BSD license.

Copyright (C) 2013, by NewBlue newblue@gmail.com

If you are also using gor, please don't hesitate to tell me by email or open an issue.
