## 安装：
yum install -y gcc gcc-c++
./configure --prefix=/usr/local/nginx --with-pcre=/usr/local/src/pcre-8.38 --with-http_stub_status_module --with-http_gzip_static_module --add-module=/usr/local/src/ngx_cache_purge-2.3
make
make install

http://nginx.org/en/docs/

1. 工作线程数和并发连接数
   worker_rlimit_nofile 20480; #每个进程打开的最大的文件数=worker_connections*2是安全的，受限于操作系统/etc/security/limits.conf
   vi /etc/security/limits.conf

* hard nofile 204800
* soft nofile 204800
* soft core unlimited
* soft stack 204800

worker_processes 4; #cpu，如果nginx单独在一台机器上
worker_processes auto;
events {
    worker_connections 10240;#每一个进程打开的最大连接数，包含了nginx与客户端和nginx与upstream之间的连接
    multi_accept on; #可以一次建立多个连接
    use epoll;
}

2. 操作系统优化
   配置文件/etc/sysctl.conf
   sysctl -w net.ipv4.tcp_syncookies=1#防止一个套接字在有过多试图连接到达时引起过载
   sysctl-w net.core.somaxconn=1024#默认128，连接队列
   sysctl-w net.ipv4.tcp_fin_timeout=10 # timewait的超时时间
   sysctl -w net.ipv4.tcp_tw_reuse=1 #os直接使用timewait的连接
   sysctl -w net.ipv4.tcp_tw_recycle = 0 #回收禁用

3. Keepalive长连接
   Nginx与upstream server：
   upstream server_pool{
           server localhost:8080 weight=1 max_fails=2 fail_timeout=30s;
           keepalive 300;  #300个长连接
   }
   同时要在location中设置：
   location /  {
               proxy_http_version 1.1;
   	proxy_set_header Upgrade $http_upgrade;
   	proxy_set_header Connection "upgrade";
   }
   客户端与nginx（默认是打开的）：
   keepalive_timeout  60s; #长连接的超时时间
   keepalive_requests 100; #100个请求之后就关闭连接，可以调大
   keepalive_disable msie6; #ie6禁用

4. 启用压缩
   gzip on;
   gzip_http_version 1.1;
   gzip_disable "MSIE [1-6]\.(?!.*SV1)";
   gzip_proxied any;
   gzip_types text/plain text/css application/javascript application/x-javascript application/json application/xml application/vnd.ms-fontobject application/x-font-ttf application/svg+xml application/x-icon;
   gzip_vary on; #Vary: Accept-Encoding
   gzip_static on; #如果有压缩好的 直接使用

5. 状态监控
   location = /nginx_status {
   	stub_status on;
   	access_log off;
   	allow <YOURIPADDRESS>;
   	deny all;
   }
   输出结果：
   Active connections: 1 
   server accepts handled requests
    17122 17122 34873 
   Reading: 0 Writing: 1 Waiting: 0 
   Active connections：当前实时的并发连接数
   accepts：收到的总连接数，
   handled：处理的总连接数
   requests：处理的总请求数
   Reading：当前有都少个读，读取客户端的请求
   Writing：当前有多少个写，向客户端输出
   Waiting：当前有多少个长连接（reading + writing）
   reading – nginx reads request header
   writing – nginx reads request body, processes request, or writes response to a client
   waiting – keep-alive connections, actually it is active - (reading + writing)

6. 实时请求信息统计ngxtop
   https://github.com/lebinh/ngxtop
   (1)安装python-pip
   yum install epel-release
   yum install python-pip
   (2)安装ngxtop
   pip install ngxtop
   (3)使用
   指定配置文件：           ngxtop -c ./conf/nginx.conf
   查询状态是200：        ngxtop -c ./conf/nginx.conf  --filter 'status == 200'
   查询那个ip访问最多： ngxtop -c ./conf/nginx.conf  --group-by remote_addr