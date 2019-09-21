# rtsp流媒体服务器

1. 基于EasyDarwin二次开发
2. 支持rtsps


## 安装部署

### 编译命令

- 编译 Linux release版本 (在 bash 环境下执行)

        make 
        
- 编译 Linux debug版本 (在 bash 环境下执行)

        make debug
        
- 直接运行(Linux)

		cd go-rtspServer
		./easydarwin
		# Ctrl + C

- 以服务启动(Linux)

		cd EasyDarwin
		./start.sh
		# ./stop.sh

- 查看界面
	
	打开浏览器输入 [http://localhost:10008](http://localhost:10008), 进入控制页面,默认用户名密码是admin/admin

- 测试推流

    使用tls，推流至rtsps
    
    ffmpeg -re -i test.mp4 -rtsp_transport tcp -c copy -f rtsp rtsps://127.0.0.1:8443/test.mp4

    不使用tls，推流至rtsp
    
    ffmpeg -re -i test.mp4 -rtsp_transport tcp -c copy -f rtsp rtsp://127.0.0.1:8554/test.mp4
    
	
- 测试播放

    使用tls
    ffplay rtsps://localhost:8443/test.mp4
    
    不使用tls
    ffplay rtsp://localhost:8554/test.mp4
	// ffplay -rtsp_transport tcp rtsp://localhost/test
	// ffplay rtsp://localhost/test 




