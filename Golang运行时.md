---
title: Golang运行时
date: 2020-12-10
---

# Golang运行时  

## 什么是运行时
任何语言都需要Runtime，Runtime个人理解是把程序语言转化为机器语言的中间层  
在Java中运行时的全程叫Java Runtime Environment，里面包含了JVM(Java虚拟机)以及一些标准函数类库。  
在C中，runtime是库代码，等同于C runtime library，里面存放一系列C程序运行所需的函数  

Runtime不同语言代表的意义并不完全相同，在所有语言中，Runtime是一个通用的抽象术语，指的是计算机程序运行的时候所需要的一切代码库，框架，平台等  

## Golang中的Runtime  

在Go中，有一个runtime库，实现了