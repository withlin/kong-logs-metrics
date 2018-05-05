<template>
    <div>
        <div style="height:10px;"></div>
        <Row>
        <Col span="8">
        <Button @click="handleSubmit" type="primary" >显示图表</Button>
        <Select v-model="model1" style="width:200px">
        <Option v-for="item in cityList" :value="item.value" :key="item.value">{{ item.label }}</Option>
        </Select>
        
        </Col>

        <Col span="1" offset="1">
        <DatePicker :value="value2" format="yyyy/MM/dd" type="daterange" placement="bottom-end" placeholder="选择起始日期" style="width: 200px">
        </DatePicker>
        </Col>
      </Row>
         <div style="height:50px;"></div>
        <div style="width:1100px;height:700px;" id="visite_volume_con"></div>
        <Table stripe :columns="columns1" :data="data1" style="width:1100px;"></Table>
    </div>
    
    
</template>

<script>
import echarts from 'echarts';
import Axios from 'axios';
import Api  from '@/api';
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
        data:['最小耗时','平均耗时','最大耗时']
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
            name: '最大耗时',
            min: 0,
            max: 50000,
            interval: 5000,
            axisLabel: {
                formatter: '{value}'
            }
        },
        {
            type: 'value',
            name: '最小耗时',
            min: 0,
            max: 2000,
            interval: 200,
            axisLabel: {
                formatter: '{value}'
            }
        }
    ],
    series: [
        {
            name:'最大耗时(ms)',
            type:'bar',
            barWidth:35,
            itemStyle:{normal:{color:'#ff9966'}},
            data:[]
        },
        {
            name:'最小耗时(ms)',
            type:'bar',
            barWidth:40,
            data:[]
        },
        {
            name:'平均耗时(ms)',
            type:'line',
            yAxisIndex: 1,
            data:[]
        }
    ]
            };
export default {
    name: 'visiteVolume',
    data () {
        return {
            //
            // result:this.handleSubmit()
        };
    },
    mounted () {
        this.$nextTick(() => {
            
        });
    },
    methods: {
        handleSubmit () {
              let visiteVolume = echarts.init(document.getElementById('visite_volume_con'));

            

               visiteVolume.setOption(option);

            //    window.addEventListener('resize', function () {
            //      visiteVolume.resize();
            //    });
                    let server=Api.MixedLineAndBar;
                    
                    Axios.get(server).then((res)=>{
                        console.log(res.data);
                        visiteVolume.hideLoading();
                        visiteVolume.showLoading();

                        
                        if(res.data.message=="ok"){
                            setTimeout(()=>{  //未来让加载动画效果明显,这里加入了setTimeout,实现2s延时
                           visiteVolume.hideLoading(); //隐藏加载动画
                             
                           visiteVolume.setOption({
                                series: [{
                                data: res.data.data.max
                               },
                               {
                                data: res.data.data.min
                               },
                               {
                                data: res.data.data.avg
                               }
                            ]
                           });
                             }, 1000 )

                        }
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
                    
        }
    }
};
</script>
