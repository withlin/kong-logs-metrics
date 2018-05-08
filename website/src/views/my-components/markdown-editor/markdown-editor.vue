<style lang="less">
    .vertical-center-modal{
        display: flex;
        align-items: center;
        justify-content: center;

        .ivu-modal{
            top: 0;
        }
    }

    .textarea-show{  
    border:-10;  
    -ms-overflow-style:none;
    overflow:-moz-scrollbars-none;
    // background-color:transparent;  
    // scrollbar-arrow-color:yellow;  
    // scrollbar-base-color:lightsalmon;  
    overflow: hidden;
    background-color: #e9e9e9;  
    height: 1000px;  
    width: 500px;
}  
</style>
<template>
    <div>
         <Table border stripe :loading="loading" :columns="columns" :data="data"></Table>
         <Page :total="50" size="small" show-elevator show-sizer></Page>
          <Modal
           v-model="logdetail"
           footer-hide="true"
           class-name="vertical-center-modal">
          <!-- <span></span> -->
           <textarea class="textarea-show" readonly="readonly">{{showlogsdetail}}</textarea>
          </Modal>
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
            logdetail:false,
            loading:true,
            showlogsdetail:'',
            columns:
            [
                    
                    {
                        title: '详情',
                        key: 'action',
                        width: 150,
                        align: 'center',
                        
                        render: (h, params) => {
                            return h('div', [
                                h('Button', {
                                    props: {
                                        type: 'primary',
                                        size: 'small'
                                    },
                                    style: {
                                        marginRight: '5px'
                                    },
                                    on: {
                                        click: () => {
                                            
                                           this.loadLogsDetail(params.row.id);
                                           this.logdetail=true;
                                        }
                                    }
                                }, '详情')
                            ]);
                        }
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
                        title: '总耗时',
                        key: 'usetime'
                    },
                    {
                        title: '开始时间',
                        key: 'starttime'
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
            this.loading=true;
            this.handleMethod();
            this.loading=false;
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
                                       "consumer":'',
                                       "usetime":`${element._source.latencies.request} ms`,
                                       "starttime":element._source.started_at
                               });

                               
                               });
                               this.data=tableData;
                           }

                            
                            
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        },
        loadLogsDetail(id){
                    
                    let server=Api.ShowLogsDetail;
                    let tableData=[];
                    let data={'ID':`${id}`};

                    console.log("===========id为============");
                    console.log(id);

                    Axios.post(server,data).then((res)=>{
                         
                           if(res.data.message=="ok"){
                            //    console.log(res.data.data);
                            //    console.log(res.data.data[0]);
                        //   let test= JSON.stringify(res.data.data, null, 2);
                        //   console.log(test);
                             this.showlogsdetail=res.data.data;
                           }

                            
                            
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        }
    }
};
</script>
