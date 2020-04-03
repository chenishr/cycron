var app = new Vue({
    el: '#sub-page-wrapper',
    data:{
        logList: [],
        log:{},
        taskId: null,
        taskLists: [],
        page: 1,
    },
    mounted(){
        this.taskId = getUrlKey('task_id',window.location.href)

        let page = getUrlKey('page',window.location.href)
        if(null != page){
            this.page = page
        }

        console.log("query TaskId:" + this.taskId)

        if(null != this.taskId){
            this.log_list(this.taskId)
        }

        this.task_list()
    },
    methods:{
        show_log_detail(logId){
            console.log(logId)

            let _this = this
            const params = new URLSearchParams();
            params.append('logId', logId);
            axios.post('/log/detail',params)
                .then(function (response) {
                    // handle success
                    _this.log = response.data.data
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                })
                .then(function () {
                    // always executed
                });

            $('#myModal').modal('show')
        },
        task_id_change(){
            this.log_list(this.taskId)
        },
        task_list(){
            let _this = this
            const params = new URLSearchParams();
            params.append('page', 1);
            params.append('page_size', 1000);
            axios.post('/task/list',params)
                .then(function (response) {
                    // handle success
                    _this.taskLists = response.data.data.list
                    console.log(_this.taskLists);
                    if(null == _this.taskId){
                        console.log("TaskId:" + _this.taskId)
                        _this.taskId = _this.taskLists[0]['Id']
                        _this.log_list(_this.taskId)
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
        log_list(taskId){
            console.log("get log")
            let _this = this
            const params = new URLSearchParams();
            params.append('taskId', taskId);
            params.append('page', this.page);
            params.append('page_size', 20);
            axios.post('/log/list',params)
                .then(function (response) {
                    // handle success
                    _this.logList = response.data.data.list
                    console.log(_this.logList);

                    $('#pageLimit').bootstrapPaginator({
                        currentPage: _this.page,
                        totalPages: response.data.data.total_page,
                        pageUrl: function (type, page, current) {
                            return "?task_id=" + _this.taskId + "&page=" + page;
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
        }
    }
})
