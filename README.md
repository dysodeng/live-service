### 直播

- 基于gin和阿里云视频直播服务的直播平台

#### 部署

- 项目采用docker部署是最佳方式，以下基于docker-compose部署
- get code
``` git clone https://github.com/DysoDeng/live-service.git ```
- 部署配置
```
cd live-service/deploy && copy .env.example .env
```
- 配置nginx：live-service.conf
- 运行docker-compose
```
docker-compose up -d
```