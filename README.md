### live-service

- 基于gin的web基础架构

#### 部署

- 项目采用docker部署是最佳方式，以下基于docker-compose部署
```
git clone https://github.com/DysoDeng/live-service.git
cd live-service/deploy && copy .env.example .env
```
- 配置nginx：live-service.conf
- 启动项目
```
docker-compose up -d
```

### 基础功能计划

- [x] 基础架构搭建
- [x] JWT Token验证
- [x] 基础中间件(Token鉴权，跨域)
- [x] 集成阿里OSS
- [x] 文件上传组件
- [x] 缓存组件
- [x] 短信组件
- [ ] 微信组件
- [ ] 集成阿里云服务组件
