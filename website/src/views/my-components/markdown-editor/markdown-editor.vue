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
         <Table border stripe :loading="loading" :columns="columns" :data="data" :size="tableSize"></Table>
         <div style="margin: 10px;overflow: hidden">
           <div style="float: right;">
            <Page :total=total :page-size=200 @on-change="changePage"  show-elevator  show-total></Page>
           </div>
         </div>
          <Modal
           v-model="logdetail"
           footer-hide=true
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
            total:200,
            tableSize: 'default',
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
        handleMethod(data){

                    let server=Api.ShowLog;
                    //  let server=Api.MixedLineAndBar;
                    let tableData=[];
                    console.log("=============页码跳转=========");
                    console.log(data);
                    if (data===undefined) {
                        data={"pagesize":200,"pagenumber":1}
                    }
                    Axios.post(server,data).then((res)=>{
                           
                           if(res.data.message=="ok"){
                               console.log("=======================================展示日志的数据=======================");
                               console.log(res.data.data);
                               res.data.data.hits.forEach(element => {
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
                               this.total=res.data.data.total;
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
                             this.showlogsdetail=res.data.data;
                           }

                            
                            
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        },
        changePage(index){
            console.log(index);
            // 200 400 From(page.PageNumber).Size(page.PageSize)
            let data={"pagesize":200,"pagenumber":(index-1)*200}
            this.handleMethod(data);
            // console.log(this.current);
        }
    }
};
</script>
