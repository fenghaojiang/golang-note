---
title: HTTP网络协议相关知识笔记
date: 2020-11-2
---


# HTTP网络协议相关

**URI**  
URI包含URL跟URN  
URI全称Uniform Resource Identifier 统一资源标识符  
URL全称Uniform Resource Locator 统一资源定位符  
URN全称Uniform Resource Name 统一资源名称  
看这中文名跟英文名就知道URI包含URL跟URN，懂的都懂，真的是给他懂完了。  


**请求报文和响应报文**  
请求报文的第一行是请求行，包含了方法字段。  
+ 请求报文 

![request](request.png) 

+ 响应报文

![response](response.png)

**HTTP方法**  
请求方法字段  
+ GET  
请求指定的页面信息，并返回实体主体
+ HEAD  
与GET请求类似，但不返回报文实体主体部分。  
主要用于确认URL的有效性以及资源的更新的日期时间。
+ POST
+ PUT
+ DELETE
+ PATCH
+ OPTIONS
+ CONNECT
+ TRACE


