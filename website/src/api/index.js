
/** 
 * @Description:
 *
 * @Author: JZT.WithLin
 * @Date: 2018-03-05
*/
let serverAddress='http://localhost:7777/v1/api';

//登录接口
const Login={
    //检查检查用户是否登录
    CheckLogin:serverAddress+'/checklogin',
    MixedLineAndBar:serverAddress+'/test',
    QueryUrlName:serverAddress+'/test/queryurlname',
    PieChart:serverAddress +'/test/piechart',
    ShowLog:serverAddress+'/showlogs',
    ShowLogsDetail:serverAddress+"/findlogdetailbyid",
    FindLogsByApiName:serverAddress+"/findlogsbyapiname"
    

}

const APIS={
    ...Login
}

export default APIS