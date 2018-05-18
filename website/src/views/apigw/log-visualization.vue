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
    // overflow: hidden;
    background-color: #e9e9e9;  
    height: 1000px;  
    width: 500px;
}  
</style>
<template>
    <div>
        <Card style="height:100px">
        <Row>
        <Col span="4" offset="2">
        <Button @click="selectdata" type="primary" >查询</Button>
        </Col>
        <Col span="3">
        <Select v-model="model1" @on-change="getOptionValue" style="width:200px" filterable clearable>
        <Option v-for="item in apiList" :value="item.value" :key="item.value" >{{ item.label }}</Option>
        </Select>
        </Col>

        <Col span="4" offset="4">
        <DatePicker type="date" :value="dateValue" @on-change="getDataValue"placeholder="选择日期" style="width: 200px"></DatePicker>
        </DatePicker>
        </Col>
        </Row>
        </Card>
         <Table  :columns="columns" :data="data" :size="tableSize"  border stripe ></Table>
         <div style="margin: 10px;overflow: hidden">
           <div style="float: right;">
            <Page :total=total :page-size=200 @on-change="changePage"  show-elevator  show-total></Page>
           </div>
           
         </div>
          <Modal
           v-model="logdetail"
           :footer-hide="true"
           class-name="vertical-center-modal">

          <!-- <span></span> -->
           <textarea class="textarea-show" readonly="readonly">{{showlogsdetail}}</textarea>
          </Modal>
    </div>
    
</template>

<script>
import Axios from 'axios';
import Api  from '@/api';
import moment from 'moment';
import _  from 'lodash';
export default {
    name: 'showlog',
    data () {
        return {
            pageNumber:1,
            model1:'',
            optionValue:'',
            dateValue:'',
            apiList:[],
            data:[],
            total:200,
            tableSize: 'default',
            logdetail:false,           
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
                                            
                                           this.loadLogsDetail(params.row.id,params.row.indexName);
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
                        title: '总耗时(ms)',
                        key: 'usetime'
                    },
                    {
                        title: '开始时间',
                        key: 'starttime',
                        sortable: true
                    },
                    {
                        title: '消费者',
                        key: 'consumer'
                    },
                    {
                        title: 'ID',
                        key: 'id'
                    },
                    {
                        title: '索引名称',
                        key: 'indexName'
                    }
            ]
        }
    },
    mounted () {
        this.$nextTick(()=>{
            this.dateValue=moment().format('YYYY.MM.DD');
            this.handleMethod();
            this.queryUrlName();         
        });
    },
    methods: {
        selectdata(page){
            let server=Api.FindLogsByApiName;
            let data=null;
            console.log("这里是==========selectdata方法=========================")
            console.log(this.model1);
            console.log(page);
            if (page==1) {
                      data={"name":this.model1,"datevalue":`logstash-${this.dateValue}`,"pagesize":200,"pagenumber":1}
                        
                    }else{
                        data={"name":this.model1,"datevalue":`logstash-${this.dateValue}`,"pagesize":200,"pagenumber":( this.pageNumber-1)*200}

                    }
                    let tableData=[];

                    if(this.dateValue !="" && this.model1==""){
                         console.log("if====================");
                        let  test={"pagesize":200,"pagenumber":1,"datevalue":`logstash-${this.dateValue}`}
                        this.handleMethod(test);
                    }
                    else if(this.dateValue=="" || this.model1 == ""){
                         console.log("这里进入了else if====================");
                        let  test={"pagesize":200,"pagenumber":1,"datevalue":`logstash-${moment().format('YYYY.MM.DD')}`}
                        this.handleMethod(test);
                    }else{
                        Axios.post(server,data).then((res)=>{
                           console.log("这里进入了else====================");
                           if(res.data.message=="ok"){
                            
                               console.log(res.data.data);
                               res.data.data.hits.forEach(element => {
                                   let date=new Date(element._source.started_at);
                                   
                                   tableData.push({
                                       "detail":'',
                                       "uris":element._source.request.uri,
                                       "methods":element._source.request.method,
                                       "upstreamurl":element._source.api.upstream_url,
                                       "name":element._source.api.name,
                                       "id":element._id,
                                       "consumer":element._source.consumer.username,
                                       "usetime":`${element._source.latencies.request}`,
                                       "starttime":moment(date).format('YYYY-MM-DD HH:mm:ss'),
                                       "indexName":element._index
                               });

                               
                               });
                               tableData=_.orderBy(tableData,['starttime'],['desc'])
                               this.data=tableData;
                               this.total=res.data.data.total;
                           }

                            
                        this.loading=false;    
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
                    }
                    

        },
        handleMethod(data){

                    let server=Api.ShowLog;
                    //  let server=Api.MixedLineAndBar;
                    let tableData=[];
                    console.log("handleMethod时间啊啊===================");
                    console.log(data);
                    let dateTimeNow= moment().format('YYYY.MM.DD');
                    console.log(dateTimeNow);
                    if (data===undefined) {
                        data={"pagesize":200,"pagenumber":1,"datevalue":`logstash-${dateTimeNow}`}
                    }
                    console.log("请求参数是：=============================",data);
                    Axios.post(server,data).then((res)=>{
                           
                           if(res.data.message=="ok"){
                            
                               console.log(res.data.data);
                               res.data.data.hits.forEach(element => {
                                   let date=new Date(element._source.started_at);
                                   
                                   tableData.push({
                                       "detail":'',
                                       "uris":element._source.request.uri,
                                       "methods":element._source.request.method,
                                       "upstreamurl":element._source.api.upstream_url,
                                       "name":element._source.api.name,
                                       "id":element._id,
                                       "consumer":element._source.consumer.username,
                                       "usetime":`${element._source.latencies.request}`,
                                       "starttime":moment.utc(element._source.started_at).local().format('YYYY-MM-DD HH:mm:ss'),
                                       "indexName":element._index
                               });

                               
                               });
                               tableData=_.orderBy(tableData,['starttime'],['desc'])
                               this.data=tableData;
                               this.total=res.data.data.total;
                           }

                            
                        this.loading=false;    
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        },
        loadLogsDetail(id,indexName){
                    
                    let server=Api.ShowLogsDetail;
                    let tableData=[];
                    let data={'ID':`${id}`,'indexname':indexName};
                    console.log("===========id为============");
                    console.log(id);

                    Axios.post(server,data).then((res)=>{
                         
                           if(res.data.message=="ok"){
                             this.showlogsdetail=res.data.data[0]._source;
                             console.log(res.data.data[0]._source);
                           }

                            
                            
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                    });
        },
        changePage(index){
            console.log("进入了changePage方法============================");
            console.log(index);
            // 200 400 From(page.PageNumber).Size(page.PageSize)
            let data={"pagesize":200,"pagenumber":(index-1)*200,"datevalue":`logstash-${moment(this.dateValue).format('YYYY.MM.DD')}`}
            this.handleMethod(data);
            this.pageNumber=index;
            this.loading=false;
            // console.log(this.current);
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
        getDataValue(value){
            let test=value.replace("-",".").replace("-",".");
            this.dateValue=test;
            console.log(this.dateValue);
            this.queryUrlName();
        },
        getOptionValue(value){
            console.log(this.model1);
        }
    }
};
</script>
