# Kong Statistical Analysis Of Performance

KongApi Gateway Log Metrics 后台是基于 go 写的，web框架采用的[gin](https://github.com/gin-gonic/gin),前台用的是Vue,UI组件
采用的是[iView](https://github.com/iview/iview)

## 预览（OverView）

效果图，如下所示：

![折柱混合（报表）](https://user-images.githubusercontent.com/22409551/40306164-33347c54-5d30-11e8-8b6f-068580a426b6.jpg)

![饼状图 （报表）](https://user-images.githubusercontent.com/22409551/40306210-67edc0ae-5d30-11e8-997e-571ade958d7b.jpg)

![表格 ](https://user-images.githubusercontent.com/22409551/40306211-684a5cb0-5d30-11e8-8223-4909924d0acf.jpg)

![日志可视化](https://user-images.githubusercontent.com/22409551/40306209-6781a626-5d30-11e8-9aca-745fc2955e61.jpg)



## Getting Started

```
//Logstash 配置参考
input {
  tcp {
    type => "tcp"
    port => 6666
    mode => "server"
  }
}
filter  {
      json {
            source => "message"
        }
      date {
        match => ["started_at", "yyyy-MM-dd HH:mm:ss,SSS", "UNIX_MS"]
        target=>"started_at"
        locale => "cn"
        timezone => "Asia/Shanghai"
    }
}
output {
  elasticsearch { hosts => ["localhost:9200"] }
  stdout { codec => rubydebug }
}

使用kong gateway的admina pi添加插件tcp log,配置的地址就写logstash的地址，端口为上面配置设置的端口

git clone https://github.com/DevWithLin/kong-logs-metrics.git

cd  kong-logs-metrics

//vendor文件夹的东西拷贝到 GOPATH目录的SRC下面

//设置ElasticSearch服务端地址
vim config.json

配置说明
{
    "elasticsearch":{
        "Host":"http://192.168.199.17:9200", //配置ES地址
        "SetSniff":false                     //SetSniff 是否为false
    },
    "go": {
        "LogDir":"",                    
        "APIPrefix":"/v1/api",              //API前缀
        "Port":7777,                        //go启动Server监控的端口
        "TokenMaxAge":7200,
        "Env":"production"                  // /*模式(开发(development)，测试(test)，产品(production)*/
	}
}

//启动后台服务
go run main.go


cd  kong-logs-metrics/website


//还原依赖  注意：Node版本使用v8.10.0 最新版本不兼容（webpack）
npm i 或者 cnpm i

//启动前端项目
npm run dev


//如果后台部署到其他的服务器上 需要修改前端地址

vim  kong-logs-metrics/website/src/api/index.js

let serverAddress='http://localhost:7777/v1/api'; //修改此处服务端地址即可

//发布项目

npm run build  //然后把项目静态文件部署到nginx上。

```
