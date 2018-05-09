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
        <Card style="height:100px">
        <Row>
        <Col span="4">
        <Button @click="handleSubmit" type="primary" >查询</Button>
        </Col>
        <Col span="3">
        <Select v-model="model1" style="width:200px">
        <Option v-for="item in cityList" :value="item.value" :key="item.value">{{ item.label }}</Option>
        </Select>
        </Col>

        <Col span="4" offset="4">
        <DatePicker :value="value2" format="yyyy/MM/dd" type="daterange" placement="bottom-end" placeholder="选择起始日期" style="width: 200px">
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
export default {
    name: 'showlog',
    data () {
        return {
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
            this.handleMethod();
           
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
                                       "starttime":this.convertUTCTimeToLocalTime(element._source.started_at)
                               });

                               
                               });
                               this.data=tableData;
                               this.total=res.data.data.total;
                           }

                            
                        this.loading=false;    
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
             this.loading=false;
            // console.log(this.current);
        },
        convertUTCTimeToLocalTime(UTCDateString){
            if(!UTCDateString){
          return '-';
        }
        function formatFunc(str) {    //格式化显示
          return str > 9 ? str : '0' + str
        }
        let date2 = new Date(UTCDateString);     //这步是关键
        let year = date2.getFullYear();
        let mon = formatFunc(date2.getMonth() + 1);
        let day = formatFunc(date2.getDate());
        let hour = date2.getHours();
        let noon = hour >= 12 ? 'PM' : 'AM';
        hour = hour>=12?hour-12:hour;
        hour = formatFunc(hour);
        let min = formatFunc(date2.getMinutes());
        let dateStr = year+'-'+mon+'-'+day+' '+noon +' '+hour+':'+min;
        return dateStr;
        }
    }
};
</script>
