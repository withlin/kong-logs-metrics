<template>
 <Card>
    <div>
        <div style="height:20px;"></div>
        <Row>
        <Col span="4">
        <Button @click="handleSubmit" type="primary" >查询</Button>
        </Col>
        <Col span="3">
        <Select v-model="model1" @on-change="getOptionValue" style="width:200px" filterable clearable>
        <Option v-for="item in apiList" :value="item.value" :key="item.value" >{{ item.label }}</Option>
        </Select>
        </Col>
          

        <Col span="4" offset="4">
         <DatePicker type="date" :value="dateValue" @on-change="getDataValue" placeholder="选择日期" style="width: 200px"></DatePicker>
        </DatePicker>
        </Col>
        
      </Row>
      <Row>
           <div style="height:25px;"></div>
           <RadioGroup v-model="animal">
            <Col span="1">
           <span @click="handleSubmit"> <Radio label="按每天24小时聚合"></Radio> </span>
            </Col>
            <Col span="3" offset="15">
           <span @click="showPieChart"> <Radio label="按照范围聚合"></Radio> </span>
           </Col>
           <!-- <Col span="1" offset="9">
           <Radio label="Heap Map聚合"></Radio>
           </Col> -->
           </RadioGroup>
           <Col span="5" offset="">
           <label>总共聚合到{{totalCount}}条记录和{{shareCount}}个分片</label>
           </Col>
           
      </Row>
         <div style="height:50px;"></div>
         <Card>
        
        <div  v-show="agg" style="width:auto;height:700px;"  id="visite_volume_con" ></div>
        <div v-show="pie" style="width:1100px;height:700px;"  id="range_pie_chart"  ></div> 
        
        <!-- style="width:1100px;height:700px;" -->
         </Card>
        
        <Table  border stripe :columns="columns" :data="data"></Table>
    </div>

    
    </Card>
    
    
</template>

<script>
import echarts from 'echarts';
import Axios from 'axios';
import Api  from '@/api';
import moment from 'moment';


const option = {
                tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'none',
            crossStyle: {
                color: '#999'
            }
        }
    },
    toolbox: {
        feature: {
            dataView: {show: false, readOnly: false},
            magicType: {show: false, type: ['line', 'bar']},
            restore: {show: false},
            saveAsImage: {show: false}
        }
    },
    legend: {
        data:['最小耗时(ms)','平均耗时(ms)','最大耗时(ms)','请求次数']
    },
    xAxis: [
        {
            type: 'category',
            data: ['0时','1时','2时','3时','4时','5时','6时','7时','8时','9时','10时','11时','12时','13时','14时','15时','16时','17时','18时','19时','20时','21时','22时','23时'],
            axisPointer: {
                type: 'none'
            }
        }
    ],
    yAxis: [
        {
            type: 'value',
            name: '最大耗时(ms)',
            min: 0,
            max: 60000,
            interval: 10000,
            axisLabel: {
                formatter: '{value}'
            }
        },
        {
            type: 'value',
            name: '最小耗时(ms)',
            min: 0,
            max: 60000,
            interval: 10000,
            axisLabel: {
                formatter: '{value}'
            }
        }
    ],
    series: [
        
        {
            name:'最大耗时(ms)',
            type:'bar',
            barWidth:40,
            yAxisIndex:0,
            itemStyle:{normal:{color:'#ff9966'}},
            data:[]
        },
        {
            name:'最小耗时(ms)',
            type:'line',
            // barWidth:40,
            yAxisIndex:0,
            data:[]
        }
        ,
        {
            name:'请求次数',
            type:'line',
            yAxisIndex:0,
            data:[]
        }
        ,
        {
            name:'平均耗时(ms)',
            type:'line',
            yAxisIndex:0,
            data:[]
        }
    ]
 };

const  optionPie = {
    title : {
        text: '网关访问速度',
        subtext: '网关',
        x:'center'
    },
    tooltip : {
        trigger: 'item',
        formatter: "{a} <br/>{b} : {c} ({d}%)"
    },
    legend: {
        orient: 'vertical',
        left: 'left',
        data: ['30ms','30ms-90ms','90ms-120ms','120ms-150ms','150ms-180ms','180ms-210ms','210ms-240ms','240ms-270ms','270ms-300ms','300ms-500ms']
    },
    series : [
        {
            name: '网关',
            type: 'pie',
            radius : '55%',
            center: ['50%', '60%'],
            data:[
            ],
            itemStyle: {
                emphasis: {
                    shadowBlur: 10,
                    shadowOffsetX: 0,
                    shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
            }
        }
    ]
};

export default {
    name: 'visiteVolume',
    data () {
        return {
            model1:'',
            dateValue:'',
            loading:true,
            shareCount:0,
            agg:true,
            pie:false,
            apiList: [],
            model1:'',
            animal: 'test',
            totalCount:0,
            columns:
            [
                   
                    {
                        title: '时间',
                        key: 'time'
                    },
                    {
                        title: '请求数量',
                        key: 'countRequest'
                    },
                    {
                        title: '最大耗时(ms)',
                        key: 'maxTime'
                    },
                    {
                        title: '最小耗时(ms)',
                        key: 'minTime'
                    },
                    {
                        title: '平均耗时(ms)',
                        key: 'avgTime'
                    }
            ],
            data:
            [
            ]
        };
    },
    mounted () {
        this.$nextTick(() => {
            this.dateValue=moment().format('YYYY.MM.DD');
            this.handleSubmit();
            this.queryUrlName();
            // this.showPieChart();
        });
    },
    methods: {
        handleSubmit () {
             this.agg=true;
             this.pie=false;
        
              let visiteVolume = echarts.init(document.getElementById('visite_volume_con'));

               visiteVolume.setOption(option);

                    let server=Api.MixedLineAndBar;
                    let data=null;
                    if (this.model1==""&& this.dateValue!=""){
                        data={"logstastname":`logstash-${this.dateValue}`}
                    }else{
                         data={"logstastname":`logstash-${this.dateValue}`,"name":this.model1}
                    }
                    Axios.post(server,data).then((res)=>{
                        visiteVolume.hideLoading();
                        visiteVolume.showLoading();
                        let tabledata=[];
                        if(res.data.message=="ok"){
                            setTimeout(()=>{  //未来让加载动画效果明显,这里加入了setTimeout,实现2s延时
                           visiteVolume.hideLoading(); //隐藏加载动画
                           this.totalCount=res.data.data.totalCount
                           visiteVolume.setOption({
                                series: [
                               {
                                data: res.data.data.max
                               },
                               {
                                data: res.data.data.min
                               }
                               ,
                               {
                                data: res.data.data.count
                               }
                               ,
                               {
                                   data:res.data.data.avg
                               }
                            ]
                           });
                             }, 0 );
                             for(let i=0; i<res.data.data.avg.length; i++){
                               
                                 console.log(name);
                                 tabledata.push({
                                     time:`${i}时`,
                                     maxTime:res.data.data.max[i],
                                     minTime:res.data.data.min[i],
                                     avgTime:res.data.data.avg[i],
                                     countRequest:res.data.data.count[i]

                                 });

                             }
                             this.shareCount=res.data.data.shareTotalCount;
                             this.data=tabledata;

                        }
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
                    
        },
        queryUrlName(){
            let server=Api.QueryUrlName;
            let apis=[]
            let data={"logstastname":`logstash-${this.dateValue}`}
            Axios.post(server,data).then((res)=>{
                        console.log(res.data);

                          if(res.data.message=="ok"){
                              for (let index = 0; index < res.data.data.length; index++) {
                                 apis.push({
                                     value:res.data.data[index].key,
                                     label:res.data.data[index].key
                                 })
                                  
                              }
                              this.apiList=apis;
                          }else{
                              console.log(res.data.data);
                              this.$Notice.warning({
                                         duration:6,
                                         title: '警告',
                                         desc:res.data.data
                              });
                              this.apiList=[];
                          }
                      
                           
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        },
        showPieChart(){ 
                      this.agg=false;
                      this.pie=true;
                      let rangechart = echarts.init(document.getElementById('range_pie_chart'));
                      rangechart.setOption(optionPie);
            
                     let server=Api.PieChart;
                     let data=null;
                    if (this.model1==""&& this.dateValue!=""){
                        data={"logstastname":`logstash-${this.dateValue}`}
                    }else{
                         data={"logstastname":`logstash-${this.dateValue}`,"name":this.model1}
                    }
                    Axios.post(server,data).then((res)=>{
                           if(res.data.message=="ok")
                            rangechart.setOption({
                                 series : [
                                     {
                                         data:res.data.data
                                     }
                             ]
                           });
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        },
        getDataValue(value){
            let test=value.replace("-",".").replace("-",".");
            this.dateValue=test;
            console.log(this.dateValue);
            this.queryUrlName();
        },
        getOptionValue(value){

        }
    }
};
</script>
