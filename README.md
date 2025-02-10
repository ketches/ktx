# ktx

ktx is an easy-to-use command line tool for kubernetes multi-cluster context management.

## Installation

```bash
go install -u github.com/poneding/ktx@latest
```

## Usage

1. Add cluster context

```bash
ktx add -f .kube/kind-cluster-01
ktx add -f .kube/kind-cluster-02
```

2. List cluster contexts

```bash
ktx list
```

3. Switch cluster context

```bash
# Switch to specified cluster context
ktx use kind-cluster-01

# Interactive switch
ktx use
```

4. Rename cluster context

```bash
# Rename specified cluster context
ktx rename kind-cluster-01

# Interactive rename
ktx rename
```

5. Remove cluster context

```bash
# Remove specified cluster context
ktx rm kind-cluster-01

# Interactive remove
ktx rm
```

6. Export cluster context

```bash
ktx export kind-cluster-01 -f .kube/export-01
```
