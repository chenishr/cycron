var app = new Vue({
    el: '#sub-page-wrapper',
    data:{
        logList: [],
    },
    mounted(){
        this.log_stat()
    },
    methods:{
        log_stat() {
            let _this = this
            axios.get('/log/stat')
                .then(function (response) {
                    // handle success
                    let stat = response.data.data
                    console.log(stat);
                    let data = {}
                    for (let i = 0; i < stat.length; i++){
                        if(!data[stat[i].day]){
                            data[stat[i].day] = {y:"",a:0,b:0,c:0,d:0,e:0}
                        }

                        data[stat[i].day].y = stat[i].day

                        switch (stat[i].status) {
                            case 0:
                                data[stat[i].day].a = stat[i].count
                                break
                            case 1:
                                data[stat[i].day].b = stat[i].count
                                break
                            case 2:
                                data[stat[i].day].c = stat[i].count
                                break
                            case 3:
                                data[stat[i].day].d = stat[i].count
                                break
                            case 4:
                                data[stat[i].day].e = stat[i].count
                                break
                        }
                    }
                    console.log(data)

                    let chartData = Object.values(data)
                    console.log(chartData)

                    init_chart(chartData)
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

/*
TASK_SUCCESS = 0 // 任务执行成功
TASK_ERROR   = 1 // 任务执行出错
TASK_TIMEOUT = 2 // 任务执行超时
TASK_CANCEL  = 3 // 任务被取消
TASK_IGNORE  = 4 // 任务被忽略,超出并行调度的限制数量
 */
function init_chart(data){
    Morris.Line({
        element: 'morris-line-chart',
        data: data,



        xkey: 'y',
        ykeys: ['a', 'b','c','d','e'],
        labels: ['成功','报错','超时','取消','忽略'],
        fillOpacity: 0.6,
        hideHover: 'auto',
        behaveLikeLine: true,
        resize: true,
        pointFillColors:['#ffffff'],
        pointStrokeColors: ['black'],
        lineColors:['green','red','red','red','red']

    });
}
