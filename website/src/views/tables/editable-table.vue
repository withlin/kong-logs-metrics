<template>
    <div>
         <Button @click="handleSubmit" type="primary" >显示图表</Button>
        <div style="width:1000px;height:600px;" id="visite_volume_con"></div>
    </div>
    
    
</template>

<script>
import echarts from 'echarts';
import Axios from 'axios';
import Api  from '@/api';
export default {
    name: 'visiteVolume',
    data () {
        return {
            //
            result:this.handleSubmit()
        };
    },
    mounted () {
        this.$nextTick(() => {
            let visiteVolume = echarts.init(document.getElementById('visite_volume_con'));

            const option = {
                tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
            crossStyle: {
                color: '#999'
            }
        }
    },
    toolbox: {
        feature: {
            dataView: {show: true, readOnly: false},
            magicType: {show: true, type: ['line', 'bar']},
            restore: {show: true},
            saveAsImage: {show: true}
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
                type: 'shadow'
            }
        }
    ],
    yAxis: [
        {
            type: 'value',
            name: '最大耗时',
            min: 0,
            max: 250,
            interval: 50,
            axisLabel: {
                formatter: '{value}'
            }
        },
        {
            type: 'value',
            name: '最小耗时',
            min: 0,
            max: 25,
            interval: 5,
            axisLabel: {
                formatter: '{value}'
            }
        }
    ],
    series: [
        {
            name:'最大耗时',
            type:'bar',
            data:this.result.max
        },
        {
            name:'最小耗时',
            type:'bar',
            data:this.result.min
        },
        {
            name:'平均耗时',
            type:'line',
            yAxisIndex: 1,
            data:this.result.avg
        }
    ]
            };

            visiteVolume.setOption(option);

            window.addEventListener('resize', function () {
                visiteVolume.resize();
            });
        });
    },
    methods: {
        handleSubmit () {

                    let server=Api.MixedLineAndBar;

                    Axios.get(server).then((res)=>{
                        console.log(res.data);

                        if(res.data.message=="ok"){
                           return res.data.data;
                        }
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
                    
        }
    }
};
</script>
