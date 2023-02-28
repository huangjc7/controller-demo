# 有点乱 待整理

* 代码生成器
  * pkg/apis/stable/v1beta1/register.go 拷贝源码
  * pkg/apis/stable/v1beta1/type.go 编写


* hack/目录必要
  * hack/boilerplate.go.txt
  * tools.go 引入包k8s.io/code-generator不使用
  * update-codegen.sh 生成代码 需要修改脚本容器
  * 项目/位置执行 go mod vendor 早期版本包管理 项目需要 引入kubernetes部分代码库
