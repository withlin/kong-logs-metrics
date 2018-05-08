<template>
    <div>
         <Table border stripe :columns="columns" :data="data"></Table>
         <Page :total="50" size="small" show-elevator show-sizer></Page>
    </div>
</template>

<script>
import Axios from 'axios';
import Api  from '@/api';
export default {
    name: 'showlog',
    data () {
        return {
            data:[],
            columns:
            [
                    {
                        title: '详情',
                        key: 'detail'
                    },
                    {
                        title: 'Uris',
                        key: 'uris'
                    },
                    {
                        title: 'Methods',
                        key: 'methods'
                    },
                    {
                        title: 'UpstreamUrl',
                        key: 'upstreamurl'
                    },
                    {
                        title: 'Name',
                        key: 'name'
                    },
                    {
                        title: '消费者',
                        key: 'consumer'
                    },
                    {
                        title: 'ID',
                        key: 'id'
                    }
            ]
        }
    },
    mounted () {
        this.$nextTick(()=>{
            this.handleMethod();
        });
    },
    methods: {
        handleMethod(){

                    let server=Api.ShowLog;
                    //  let server=Api.MixedLineAndBar;
                    let tableData=[];
                    Axios.get(server).then((res)=>{
                           if(res.data.message=="ok"){
                               console.log("=======================================展示日志的数据");
                               console.log(res.data.data);
                               res.data.data.forEach(element => {
                                   tableData.push({
                                       "detail":'',
                                       "uris":element._source.request.uri,
                                       "methods":element._source.request.method,
                                       "upstreamurl":element._source.api.upstream_url,
                                       "name":element._source.api.name,
                                       "id":element._id,
                                       "consumer":''
                               });

                               console.log(tableData);
                               this.data=tableData;
                               });
                               
                           }
                            
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        }
    }
};
</script>
