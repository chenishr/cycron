const PAGE_DASH = "cron-dash.html"

var app = new Vue({
    el: '#sub-page-wrapper',
    data:{
        email: "",
        password: "",
    },
    mounted(){
    },
    methods:{
        login(){
            if("" == this.email){
                alert("请输入邮箱")
                return true;
            }

            let _this = this
            const params = new URLSearchParams();
            params.append('email', this.email);
            params.append('password', this.password);
            axios.post('/user/login', params)
                .then(function (response) {
                    // handle success

                    console.log("请求完成");
                    if(0 !== response.data.errno){
                        alert("用户名或密码错误")
                        return true
                    }

                    // 保存 token
                    localStorage.setItem('cycron-token',response.data.data.token)

                    // 重定向
                    parent.location.reload();
                })
                .catch(function (error) {
                    // handle error
                    console.log("请求失败");
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        }
    }
})
