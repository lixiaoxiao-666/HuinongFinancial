# 惠农金融项目 Docker & Kubernetes 部署指南

## 项目概述

惠农金融项目是一个前后端分离的金融服务平台，包含：
- **后端服务**: Go + Gin框架，提供API服务
- **前端管理端**: Vue 3 + Element Plus，管理员界面
- **前端用户端**: Vue 3 + Element Plus，用户界面

## 🚀 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+
- Kubernetes 1.20+ (用于K8s部署)
- kubectl (用于K8s部署)

### 一键部署

```bash
# Kubernetes部署
./quick-deploy.sh

# 本地Docker Compose部署
./quick-deploy.sh local

# 仅构建镜像
./quick-deploy.sh build-only
```

### 本地开发部署 (Docker Compose)

1. **克隆项目**
```bash
git clone <repository-url>
cd HuinongFinancial
```

2. **使用Docker Compose启动**
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

3. **访问服务**
- 前端用户端: http://localhost:3000
- 前端管理端: http://localhost:3001
- 后端API: http://localhost:8080

### 生产环境部署 (Kubernetes)

#### 方式1：完整部署（推荐）

```bash
# 构建镜像并部署
./quick-deploy.sh

# 或者分步执行
./build-images.sh
./deploy-k8s.sh
```

#### 方式2：分步部署

```bash
# 1. 构建Docker镜像
./build-images.sh

# 2. 分步部署到Kubernetes
./deploy-k8s.sh deploy-step

# 3. 部署HPA自动扩缩容
./deploy-k8s.sh deploy-hpa
```

#### 方式3：使用完整配置文件

```bash
# 直接应用完整配置
kubectl apply -f kubernetes/complete-deployment.yaml

# 应用HPA配置
kubectl apply -f kubernetes/hpa.yaml
```

## 📁 项目结构

```
HuinongFinancial/
├── backend/                    # Go后端服务
│   ├── Dockerfile
│   ├── main.go
│   ├── go.mod
│   └── ...
├── frontend/
│   ├── admin/                  # Vue管理端
│   │   ├── Dockerfile
│   │   ├── nginx.conf
│   │   ├── package.json
│   │   ├── .dockerignore
│   │   └── ...
│   └── users/                  # Vue用户端
│       ├── Dockerfile
│       ├── nginx.conf
│       ├── package.json
│       ├── .dockerignore
│       └── ...
├── kubernetes/                 # K8s配置文件
│   ├── complete-deployment.yaml    # 完整部署配置
│   ├── hpa.yaml                    # 自动扩缩容配置
│   ├── namespace.yaml
│   ├── backend-deployment.yaml
│   ├── frontend-admin-deployment.yaml
│   ├── frontend-users-deployment.yaml
│   └── ingress.yaml
├── build-images.sh            # 镜像构建脚本
├── deploy-k8s.sh             # K8s部署脚本
├── quick-deploy.sh           # 一键部署脚本
├── docker-compose.yml        # 本地开发配置
└── README.md                 # 本文档
```

## 🔧 配置说明

### 端口配置

- **后端服务**: 8080
- **前端用户端**: 3000
- **前端管理端**: 3001

### 资源配置

#### 后端服务
- **请求资源**: 256Mi内存, 200m CPU
- **限制资源**: 512Mi内存, 500m CPU
- **副本数**: 2-10 (自动扩缩容)

#### 前端服务
- **请求资源**: 64Mi内存, 50m CPU
- **限制资源**: 128Mi内存, 100m CPU
- **副本数**: 2-5/8 (自动扩缩容)

### 环境变量

- `GIN_MODE`: 运行模式 (debug/release)
- `TZ`: 时区设置 (Asia/Shanghai)
- `CONFIG_PATH`: 配置文件路径

## 🛠️ 常用命令

### Docker Compose

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重新构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f [service_name]

# 进入容器
docker-compose exec [service_name] sh
```

### Kubernetes

```bash
# 查看所有资源
kubectl get all -n huinong-financial

# 查看Pod日志
kubectl logs -f deployment/huinong-backend -n huinong-financial

# 进入Pod
kubectl exec -it deployment/huinong-backend -n huinong-financial -- sh

# 端口转发 (本地测试)
kubectl port-forward service/huinong-backend 8080:8080 -n huinong-financial

# 扩缩容服务
./deploy-k8s.sh scale huinong-backend 5

# 重启服务
./deploy-k8s.sh restart huinong-backend

# 查看HPA状态
kubectl get hpa -n huinong-financial

# 删除所有资源
kubectl delete namespace huinong-financial
```

## 📊 监控和维护

### 健康检查

所有服务都配置了完整的健康检查：
- **Liveness Probe**: 存活性检查
- **Readiness Probe**: 就绪性检查  
- **Startup Probe**: 启动检查（后端服务）

### 自动扩缩容 (HPA)

- **后端服务**: CPU 70%, 内存 80%, 2-10副本
- **前端服务**: CPU 70%, 内存 80%, 2-5/8副本
- **扩容策略**: 快速扩容，缓慢缩容

### 资源监控

```bash
# 查看资源使用情况
kubectl top pods -n huinong-financial
kubectl top nodes

# 查看HPA状态
kubectl get hpa -n huinong-financial

# 查看事件
kubectl get events -n huinong-financial
```

## 🌐 访问方式

### NodePort访问

- 前端用户端: http://\<NODE_IP\>:30080
- 前端管理端: http://\<NODE_IP\>:30081
- 后端API: http://\<NODE_IP\>:30082

### Ingress访问

配置hosts文件后访问：
- 前端用户端: http://huinong-users.local
- 前端管理端: http://huinong-admin.local
- 后端API: http://huinong-api.local

### 本地开发访问

- 前端用户端: http://localhost:3000
- 前端管理端: http://localhost:3001
- 后端API: http://localhost:8080

## 🔍 故障排除

### 常见问题

1. **镜像构建失败**
   - 检查Docker是否正常运行
   - 确保网络连接正常，能够下载依赖
   - 检查.dockerignore文件是否正确排除node_modules

2. **Pod启动失败**
   - 查看Pod日志: `kubectl logs <pod-name> -n huinong-financial`
   - 检查资源限制是否合理
   - 确保镜像已正确构建

3. **服务无法访问**
   - 检查Service配置
   - 确认端口映射正确
   - 检查防火墙设置

4. **前端无法连接后端**
   - 检查nginx代理配置
   - 确认后端服务正常运行
   - 验证网络连通性

5. **nginx启动失败**
   - 检查nginx配置语法
   - 确认后端服务域名解析
   - 查看错误页面配置

### 日志查看

```bash
# 查看所有服务日志
./deploy-k8s.sh logs huinong-backend
./deploy-k8s.sh logs huinong-admin
./deploy-k8s.sh logs huinong-users

# 查看实时日志
kubectl logs -f deployment/huinong-backend -n huinong-financial
```

### 性能调优

```bash
# 手动扩容
./deploy-k8s.sh scale huinong-backend 5

# 查看资源使用
kubectl top pods -n huinong-financial

# 调整HPA参数
kubectl edit hpa huinong-backend-hpa -n huinong-financial
```

## 🔐 安全建议

1. **密码管理**
   - 使用Kubernetes Secrets管理敏感信息
   - 定期更换访问密钥

2. **网络安全**
   - 配置网络策略限制Pod间通信
   - 使用TLS加密传输

3. **镜像安全**
   - 定期更新基础镜像
   - 扫描镜像漏洞

4. **资源限制**
   - 设置合理的资源请求和限制
   - 配置安全上下文

## 📞 支持

如有问题，请：
1. 查看日志文件
2. 检查配置文件
3. 参考故障排除章节
4. 联系开发团队

---

**版本**: 2.0.0  
**更新日期**: 2024年1月  
**维护者**: AI Assistant 