# ktx

ktx 是一个针对多集群上下文管理的命令行工具，简单易用。

## 安装

```bash
go install -u github.com/poneding/ktx@latest
```

## 使用

1. 添加集群上下文

```bash
ktx add -f .kube/kind-cluster-01
ktx add -f .kube/kind-cluster-02
```

2. 列出集群上下文

```bash
ktx list
```

3. 切换集群上下文

```bash
# 切换到指定集群上下文
ktx use kind-cluster-01

# 交互式切换
ktx use
```

4. 重命名集群上下文

```bash
# 重命名指定集群上下文
ktx rename kind-cluster-01

# 交互式重命名
ktx rename
```

5. 删除集群上下文

```bash
# 删除指定集群上下文
ktx rm kind-cluster-01

# 交互式删除
ktx rm
```

6. 导出集群上下文

```bash
ktx export kind-cluster-01 -f .kube/export-01
```
