const PAGE_LOGIN = "cron-login.html"
const PAGE_TASK = "cron-tasks.html"
const PAGE_DASH = "cron-dash.html"

var app = new Vue({
    el: '#wrapper',
    data:{
        actIndex: 0,
        user: {},
        uptUser: {},
        modTitle: "更新用户信息",
        target: PAGE_TASK
    },
    mounted(){
        this.get_user()
    },
    methods:{
        logout(){
            localStorage.setItem('cycron-token','')
            window.location.reload()
        },
        open_model(index){
            if (-1 == index){
                this.uptUser = this.init_user()
                this.modTitle = "添加用户"
            }else{
                this.uptUser = this.user
                this.uptUser['Passwd'] = ""
                this.modTitle = "更新用户信息"
            }

            console.log(this.task)

            $('#myModal').modal('show')
        },
        init_user(){
            return {
                UserName:"",
                Email: "",
                Passwd: "",
            }
        },
        select_menu(i){
            this.actIndex = i
        },
        get_user(){
            let _this = this
            axios.get('/user/info',)
                .then(function (response) {
                    // handle success
                    if(response.data.errno === 1000){
                        _this.target = PAGE_LOGIN
                        return;
                    }

                    _this.user = response.data.data
                    console.log(_this.user)
                    _this.target = PAGE_DASH
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        },
        save_user(){
            if ("" == this.uptUser.UserName){
                alert("用户名称不能为空")
                return false
            }
            if ("" == this.uptUser.Email){
                alert("邮箱不能为空")
                return false
            }

            let _this = this
            const params = new URLSearchParams();
            params.append('user', JSON.stringify(this.uptUser));
            params.append('newPassword', this.uptUser.Passwd);
            axios.post('/user/save', params)
                .then(function (response) {
                    // handle success
                    console.log(response)
                    //_this.task_list()
                    window.location.reload()
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        },
    }
})
