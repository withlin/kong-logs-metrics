
/** 
 * @Description:
 *
 * @Author: JZT.WithLin
 * @Date: 2018-03-05
*/
let serverAddress='http://localhost:9001';

//登录接口
const Login={
    //检查检查用户是否登录
    CheckLogin:serverAddress+'/api/dataplatform/UserLogin'
}



const APIS={
    ...Login
}

export default APIS