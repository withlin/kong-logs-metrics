# KongApi Gateway Log Metrics

KongApi Gateway Log Metrics 后台是基于 go 写的，web框架采用的[gin](https://github.com/gin-gonic/gin),前台用的是Vue,UI组件
采用的是[iView](https://github.com/iview/iview)

## 预览（OverView）

效果图，如下所示：

![折柱混合（报表）](https://github.com/WithLin0510/kong-logs-metrics/tree/master/pictrue/echart.png)

![饼状图 （报表）](https://github.com/WithLin0510/kong-logs-metrics/tree/master/pictrue/piechart.png)

![表格 ](https://github.com/WithLin0510/kong-logs-metrics/tree/master/pictrue/viewtable.png)

![日志可视化](https://github.com/WithLin0510/kong-logs-metrics/tree/master/pictrue/logview.png)



## Getting Started

```

git clone https://github.com/WithLin0510/kong-logs-metrics.git

cd  kong-logs-metrics

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

//还原依赖
npm i 或者 cnpm i

//启动前端项目
npm run dev


//如果后台部署到其他的服务器上 需要修改前端地址

vim  kong-logs-metrics/website/src/api/index.js

let serverAddress='http://localhost:7777/v1/api'; //修改此处服务端地址即可

//发布项目

npm run build  //然后把项目静态文件部署到nginx上。

```
