---
title: Golang中的reflect
date: 2020-11-21
---


# 什么是reflect?  
很多人来问我说:离老师发生甚么事辣?  
我一看，哦，原来是reflect是甚么?  
reflect翻译过来就是反射，实现了运行时的反射能力，能够让程序操作不同类型的对象。  
reflect.TypeOf能获取类型信息，reflect.ValueOf能够获取数据的运行时表示  

反射包中的所有方法基本都是围绕着Type和Value这两个类型设计的。我们通过reflect.TypeOf、reflect.ValueOf可以将一个普通的变量转换成反射包提供的Type和Value  
