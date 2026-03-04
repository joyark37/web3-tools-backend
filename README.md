# Web3 Tools Backend

Go/Gin API 服务器，用于 Web3 招聘平台。

## 技术栈

- **Go 1.21+**
- **Gin** - Web 框架
- **PostgreSQL** - 数据库

## 项目结构

```
├── cmd/
│   ├── server/           # 应用入口
│   └── migrate/          # 数据库迁移
├── internal/
│   ├── handler/          # HTTP 处理层
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑层
├── migrations/          # SQL 迁移文件
├── .env                 # 环境配置
└── .env.example        # 环境配置示例
```

## 快速开始

### 1. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入数据库配置
```

### 2. 创建数据库

```bash
createdb web3_recruitment
```

### 3. 运行迁移

```bash
go run ./cmd/migrate
```

### 4. 启动服务

```bash
# 开发模式
go run ./cmd/server

# 或使用编译后的二进制
./cmd/server/server
```

## API 端点

### 职位 (Jobs)

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/api/jobs` | 获取职位列表 |
| GET | `/api/jobs/:id` | 获取职位详情 |
| POST | `/api/jobs` | 发布新职位 |

### 申请 (Applications)

| 方法 | 端点 | 描述 |
|------|------|------|
| POST | `/api/applications` | 投递申请 |
| GET | `/api/applications/job/:jobId` | 获取职位的申请列表 |

## 示例

### 获取职位列表

```bash
curl http://localhost:8080/api/jobs
```

### 发布职位

```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Solidity Engineer",
    "company": "DeFi Protocol",
    "email": "hr@defiprotocol.com",
    "location": "Remote",
    "type": "full-time",
    "salaryMin": 150,
    "salaryMax": 250,
    "category": "engineering",
    "description": "Looking for an experienced Solidity developer...",
    "requirements": "5+ years of experience...",
    "benefits": "Competitive salary, Remote-first...",
    "tags": ["Solidity", "DeFi", "Smart Contracts"]
  }'
```

### 投递申请

```bash
curl -X POST http://localhost:8080/api/applications \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "resumeText": "简历内容...",
    "resumeFilename": "resume.pdf",
    "coverLetter": "求职信内容..."
  }'
```

## License

MIT
