# **********************************************************
# * Author           : ricky97gr
# * Email            : forgocode@163.com
# * Github           : https://github.com/ricky97gr
# * Create Time      : 2024-03-14 15:50
# * FileName         : prepare.sh
# * Description      : 
# **********************************************************

# 设置错误时立即退出
set -e

# 提前获取sudo权限，避免多次输入密码
echo "正在获取sudo权限..."
sudo echo "sudo权限获取成功!"

echo "正在检查docker是否安装..."
if [[ -z $(which docker) ]]; then
	echo "错误: docker 未安装，请先安装docker!"
	exit 1
else
	echo "docker 已安装"
fi

check_image(){
	echo "正在检查镜像..."
	if [[ -z $(sudo docker image list | grep mysql) ]];then
		echo "镜像 mysql 不存在，开始拉取..."
		pull_mysql
	else
		echo "镜像 mysql 已存在"
	fi
	if [[ -z $(sudo docker image list | grep redis) ]];then
		echo "镜像 redis 不存在，开始拉取..."
		pull_redis
	else
		echo "镜像 redis 已存在"
	fi
	if [[ -z $(sudo docker image list | grep nginx) ]];then
  		echo "镜像 nginx 不存在，开始拉取..."
  		pull_nginx
  	else
  		echo "镜像 nginx 已存在"
  	fi

  	if [[ -z $(sudo docker image list | grep mongo) ]];then
    		echo "镜像 mongo 不存在，开始拉取..."
    		pull_mongo
    	else
    		echo "镜像 mongo 已存在"
    	fi
}

pull_mysql(){
	echo "开始拉取 bitnami/mysql:latest..."
	sudo docker pull bitnami/mysql
	echo "拉取 bitnami/mysql:latest 成功!"
}

pull_redis(){
	echo "开始拉取 bitnami/redis:latest..."
	sudo docker pull bitnami/redis
	echo "拉取 bitnami/redis:latest 成功!"
}

pull_mongo(){
	echo "开始拉取 bitnami/mongodb:latest..."
	sudo docker pull bitnami/mongodb
	echo "拉取 bitnami/mongodb:latest 成功!"
}

pull_nginx(){
	echo "开始拉取 bitnami/nginx:latest..."
	sudo docker pull bitnami/nginx
	echo "拉取 bitnami/nginx:latest 成功!"
}

run_mysql(){
	echo "正在启动 mysql 容器..."
	# 先检查是否存在mysql容器，存在则停止并删除
	if [[ ! -z $(sudo docker ps -a | grep mysql) ]]; then
		echo "mysql 容器已存在，正在停止并删除..."
		sudo docker stop mysql || true
		sudo docker rm mysql || true
	fi
	id=$(sudo docker image list | grep mysql | awk -F ' ' '{print $3}')
	if [[ -z $id ]]; then
		echo "错误: 未找到 mysql 镜像!"
		exit 1
	fi
	sudo docker run --restart always --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -e TZ=Asia/Shanghai -e MYSQL_DATABASE=test -d $id
	echo "mysql 容器启动成功!"
}

run_redis(){
	echo "正在启动 redis 容器..."
	# 先检查是否存在redis容器，存在则停止并删除
	if [[ ! -z $(sudo docker ps -a | grep redis) ]]; then
		echo "redis 容器已存在，正在停止并删除..."
		sudo docker stop redis || true
		sudo docker rm redis || true
	fi
	id=$(sudo docker image list | grep redis | awk -F ' ' '{print $3}')
	if [[ -z $id ]]; then
		echo "错误: 未找到 redis 镜像!"
		exit 1
	fi
	sudo docker run --restart always --name redis -p 6379:6379 -e ALLOW_EMPTY_PASSWORD=yes -d $id 
	echo "redis 容器启动成功!"
}

# 执行主流程
echo "开始准备数据库环境..."
check_image
run_mysql
run_redis

echo
echo "=================================="
echo "测试环境准备完成!"
echo "=================================="
echo "mysql:"
echo -e "database: \e[32mtest\e[0m"
echo -e "user:     \e[32mroot\e[0m"
echo -e "passwd:   \e[32m123456\e[0m"
echo 
echo "redis:"
echo -e "redis 无密码"
echo "=================================="
