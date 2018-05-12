<style lang="less">
@import './login.less';
</style>

<template>
    <div class="login" @keydown.enter="handleSubmit">
        <div class="login-con">
            <Card :bordered="false">
                <p slot="title">
                    <Icon type="log-in"></Icon>
                    网关性能统计
                </p>
                <div class="form-con">
                    <Form ref="loginForm" :model="form" :rules="rules">
                        <FormItem prop="userName">
                            <Input v-model="form.userName" placeholder="请输入用户名">
                            <span slot="prepend">
                                <Icon :size="16" type="person"></Icon>
                            </span>
                            </Input>
                        </FormItem>
                        <FormItem prop="password">
                            <Input type="password" v-model="form.password" placeholder="请输入密码">
                            <span slot="prepend">
                                <Icon :size="14" type="locked"></Icon>
                            </span>
                            </Input>
                        </FormItem>
                        <FormItem>
                            <Button @click="handleSubmit" type="primary" long>登录</Button>
                        </FormItem>
                    </Form>
                    <p class="login-tip">请输入用户名和密码</p>
                </div>
            </Card>
        </div>
    </div>
</template>

<script>
import Cookies from 'js-cookie';
import Axios from 'axios';
import Api  from '@/api';
export default {
    data () {
        return {
            form: {
                userName: 'admin',
                password: ''
            },
            rules: {
                userName: [
                    { required: true, message: '账号不能为空', trigger: 'blur' }
                ],
                password: [
                    { required: true, message: '密码不能为空', trigger: 'blur' }
                ]
            }
        };
    },
    methods: {
        handleSubmit () {
            this.$refs.loginForm.validate((valid) => {
                if (valid) {
                    
                    this.$store.commit('setAvator', 'https://avatars1.githubusercontent.com/u/22409551?s=400&u=bafa72dbbfd895c17aa4bbbfeff2d2de164db146&v=4');
                    
                    let data={'username':this.form.userName,'password':this.form.password};

                    
                    let server=Api.CheckLogin;

                     this.$router.push('/home');
                    Axios.post(server,data).then((res)=>{
                        console.log(res.data);
                        
                        if(res.data.message=="ok"){
                            Cookies.set('user', this.form.userName);
                            Cookies.set('password', this.form.password);
                            this.$Message.success('登录成功');
                            this.$router.push('/home');
                        }else{
                            this.$Message.error('登录失败，账户或者密码错误!');
                            this.$router.push('/login');
                        }
                    }).catch((err)=>{
                        this.$Message.error(err.message);
                        console.log(err);
                         this.$router.push('/home');
                    });
                    if (food === 'pizza')
{ 
    alert('Pizza ;-)');  
}
else
{  
    alert('Not pizza ;-(');
}

function foo(items) {
  return items
    .filter(item => item.checked)
    .map(item => item.value)
  ;
}

const food = [
  'pizza',
  'burger',
  'pasta',
]

        //  if (this.form.userName === 'iview_admin') {
        //                 Cookies.set('access', 0);
        //             } else {
        //                 Cookies.set('access', 1);
        //             }
        //             this.$router.push({
        //                 name: 'home_index'
        //             });
        //         }
        //     });
        }
    }
};
</script>

<style>

</style>
