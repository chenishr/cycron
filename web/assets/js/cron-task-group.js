var app = new Vue({
    el: '#sub-page-wrapper',
    data:{
        taskGroupLists: [],
        taskGroup:{},
        page: 1,
    },
    mounted(){
        let page = getUrlKey('page',window.location.href)

        if(null != page){
            this.page = page
        }
        this.task_group_list()
    },
    methods:{
        open_model(index){
            if (-1 == index){
                this.taskGroup = this.init_task_group()
            }else{
                this.taskGroup = this.taskGroupLists[index]
            }

            console.log(this.taskGroup)

            $('#myModal').modal('show')
        },
        save_task_group(){
            if ("" == this.taskGroup.GroupName){
                alert("分组名称不能为空")
                return false
            }

            this.taskGroup.UserId = parseInt(this.taskGroup.UserId)
            this.taskGroup.Id = parseInt(this.taskGroup.Id)

            console.log(this.taskGroup)

            let _this = this
            const params = new URLSearchParams();
            params.append('taskGroup', JSON.stringify(this.taskGroup));
            axios.post('/group/save', params)
                .then(function (response) {
                    // handle success
                    console.log(response)
                    //_this.task_group_list()
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
        init_task_group(){
            return {
                Id:0,
                GroupName: '',
                Description: '',
                UserId: 0,
            }
        },
        task_group_list(){
            let _this = this

            const params = new URLSearchParams();
            params.append('page', this.page);
            params.append('page_size', 20);
            axios.post('/group/list',params)
                .then(function (response) {
                    // handle success
                    _this.taskGroupLists = response.data.data.list
                    console.log(_this.taskGroupLists);

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
    }
})
