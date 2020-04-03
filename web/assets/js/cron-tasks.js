var app = new Vue({
    el: '#sub-page-wrapper',
    data:{
        taskLists: [],
        taskGroupLists: [],
        task: {},
        page: 1,
    },
    mounted(){
        let page = getUrlKey('page',window.location.href)

        if(null != page){
            this.page = page
        }
        this.task_list()
        this.task_group_list()
        this.task = this.init_task()

    },
    methods:{
        open_model(index){
            if (-1 == index){
                this.task = this.init_task()
            }else{
                this.task = this.taskLists[index]
            }

            console.log(this.task)

            $('#myModal').modal('show')
        },
        del_task(taskId){
            if ("" === taskId){
                alert("任务ID不能为空")
                return false
            }

            if( !confirm("确定要删除任务吗")){
                return false
            }

            const params = new URLSearchParams();
            params.append('taskId', taskId);
            axios.post('/task/del', params)
                .then(function (response) {
                    // handle success
                    console.log(response)
                    if(0 == response.data.errno){
                        alert("任务删除成功")
                        window.location.reload()
                    }else {
                        alert("任务删除失败")
                    }
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        },
        update_task_status(taskId,taskStatus){
            if ("" === taskId){
                alert("任务ID不能为空")
                return false
            }
            if ("" === taskStatus){
                alert("任务状态不能为空")
                return false
            }

            const params = new URLSearchParams();
            params.append('taskId', taskId);
            params.append('taskStatus', taskStatus);
            axios.post('/task/update_status', params)
                .then(function (response) {
                    // handle success
                    console.log(response)
                    if(0 == response.data.errno){
                        alert("任务状态更新成功")
                        window.location.reload()
                    }else {
                        alert("任务状态更新失败")
                    }
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        },
        run_task(taskId){
            if ("" === taskId){
                alert("任务ID不能为空")
                return false
            }

            const params = new URLSearchParams();
            params.append('taskId', taskId);
            axios.post('/task/run', params)
                .then(function (response) {
                    // handle success
                    console.log(response)
                    if(0 == response.data.errno){
                        alert("任务执行成功")
                    }else {
                        alert("任务执行失败")
                    }
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        },
        save_task(){
            if ("" == this.task.TaskName){
                alert("任务名称不能为空")
                return false
            }
            if ("" == this.task.CronSpec){
                alert("cron 表达式不能为空")
                return false
            }
            if ("" == this.task.Command){
                alert("命令不能为空")
                return false
            }

            if ("" == this.task.Concurrent ){
                this.task.Concurrent = 1
            }

            if ("" == this.task.GroupId ){
                this.task.GroupId = 0
            }

            if ("" == this.task.Id ){
                this.task.Id = 0
            }

            if ("" == this.task.Timeout ){
                this.task.Timeout = 0
            }

            delete this.task.NextTime
            delete this.task.PrevTime

            this.task.Concurrent = parseInt(this.task.Concurrent)
            this.task.GroupId = parseInt(this.task.GroupId)
            this.task.Id = parseInt(this.task.Id)
            this.task.Timeout = parseInt(this.task.Timeout)
            this.task.Notify = parseInt(this.task.Notify)

            console.log(this.task)

            let _this = this
            const params = new URLSearchParams();
            params.append('task', JSON.stringify(this.task));
            axios.post('/task/save', params)
                .then(function (response) {
                    // handle success
                    console.log(response)
                    //_this.task_list()
                    if(0 == response.data.errno){
                        window.location.reload()
                    }else {
                        alert("任务保存失败")
                    }
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        },
        init_task(){
            return {
                Id:0,
                TaskName: '',
                Description: '',
                GroupId: 0,
                Concurrent: '',
                CronSpec: '',
                Command: '',
                Timeout: '',
                Notify: 0,
                NotifyEmail: '',
                NextTime: '-',
                PrevTime: '-'
            }
        },
        task_list(){
            let _this = this

            const params = new URLSearchParams();
            params.append('page', this.page);
            params.append('page_size', 20);
            axios.post('/task/list',params)
                .then(function (response) {
                    // handle success
                    _this.taskLists = response.data.data.list
                    console.log(_this.taskLists);

                    $('#pageLimit').bootstrapPaginator({
                        currentPage: _this.page,
                        totalPages: response.data.data.total_page,
                        pageUrl: function (type, page, current) {
                            return "?page=" + page;
                        },
                        size:"normal",
                        bootstrapMajorVersion: 3,
                        alignment:"right",
                        numberOfPages:10,
                        itemTexts: function (type, page, current) {
                            switch (type) {
                                case "first": return "首页";
                                case "prev": return "上一页";
                                case "next": return "下一页";
                                case "last": return "末页";
                                case "page": return page;
                            }
                        }
                    });
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });
        },
        task_group_list(){
            let _this = this
            const params = new URLSearchParams();
            params.append('page', 1);
            params.append('page_size', 1000);
            axios.post('/group/list',params)
                .then(function (response) {
                    // handle success
                    _this.taskGroupLists = response.data.data.list
                    console.log(_this.taskGroupLists);
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

function getUrlKey(name,url){
    return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)').exec(url) || [, ""])[1].replace(/\+/g, '%20')) || null
}
