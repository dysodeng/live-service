### 直播

- 基于gin的web基础架构

#### 部署

- 项目采用docker部署是最佳方式，以下基于docker-compose部署
- 获取项目代码
```
git clone https://github.com/DysoDeng/live-service.git
```
- 部署配置
```
cd live-service/deploy && copy .env.example .env
```
- 配置nginx：live-service.conf
- 运行docker-compose
```
docker-compose up -d
```
