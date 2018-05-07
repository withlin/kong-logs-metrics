<template>
<div style="width:1100px;height:700px;"  id="range_chart"></div> 
</template>

<script>
import echarts from 'echarts';
import Axios from 'axios';
import Api  from '@/api';
const option = {
    title : {
        text: '同名数量统计',
        subtext: '纯属虚构',
        x:'center'
    },
    tooltip : {
        trigger: 'item',
        formatter: "{a} <br/>{b} : {c} ({d}%)"
    },
    legend: {
        type: 'scroll',
        orient: 'vertical',
        right: 10,
        top: 20,
        bottom: 20,
        data: data.legendData,

        selected: data.selected
    },
    series : [
        {
            name: '姓名',
            type: 'pie',
            radius : '55%',
            center: ['40%', '50%'],
            data: data.seriesData,
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
     name:'rangechart',
     data(){
         return{
            
         }
     },
     mounted () {
         //do something
     },
     methods: {
         // do something
         
     }
}
</script>


