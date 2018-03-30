/**
 * Created by Administrator on 2017/7/3.
 */

var webSocketAddress = "ws://localhost:9090/responseInfo";
var process_type = [];
var xSer = [];
var gData = {};

/**
 * 传递数据格式为  {
 *                          "login":{
 *                              average:  object,
 *                              Size: object,
 *                          };
 *                  }
 */
function exist(key) {
    for(var index in process_type) {
        if(process_type[index] == key ) {
            return true
        }
    }
    return false
}
function zeroSlice() {
    var data = [];
    for(var k in gData) {
        if(gData[k] != "undefined") {
            for (var i = 0; i < gData[k].length; i++) {
                data.push(0)
            }
            return data
        }
    }
    return data;
}

function parseData(rData) {
    for(var key in rData) {
        if(!exist(key)) {
            if(!gData[key]) {
                gData[key] = {
                    average:0,
                    size:0
                };
            }
            process_type.push(key);
            gData[key].average = [];
            gData[key].average = zeroSlice();
            gData[key].size = [];
            gData[key].size = zeroSlice();
        }
    }

    for(var key in gData) {
        if (!(key in rData)) {
            gData[key].average.push(0);
            gData[key].size.push(0);
        }else {
            gData[key].average.push(rData[key].average.toFixed(2));
            gData[key].size.push(rData[key].size);
        }
    }

}
var defaultTyle = "line"
function getSeries(typ) {
    var data = [];
    for(var key in process_type) {
        data.push({
            name:process_type[key],
            type: defaultTyle,
            stack: 'ms',
            data: gData[process_type[key]][typ],
            markLine:{
                data:[{
                    type:'average',
                    name:'平均值'
                }]
            },
            markPoint:{
                data:[
                    {type:'max',name:'da'},
                    {type:"min",name:'small'}
                ]
            }
        })
    }
    return data;
}
function Name(dtype) {
    switch( dtype){
        case "size":
            return "个数"
        case "average":
            return '响应时间';
    }
    return "---"
}
function draw(dtype) {
    option = {
        title: {      //标题组件
            text: Name(dtype)
        },
        tooltip: {    //提示框组件
            trigger: 'axis'
        },
        legend: {     //图例组件
            data: process_type
        },
        toolbox: {     //工具栏
            show : true,
            feature : {
                mark : {show: true},
                dataView : {show: true, readOnly: false},
                magicType : {show: true, type: ['line', 'bar']},
                restore : {show: true},
                saveAsImage : {show: true}
            }
        },
        calculable : true,
        xAxis: [
            {       //直角坐标系 grid 中的 x 轴
                type: 'category',
                boundaryGap: false,
                data: xSer
            }
        ],
        yAxis: [
            {       //直角坐标系 grid 中的 y 轴
                type: 'value',
                axisLabel : {
                    formatter: '{value}'
                }
            }
        ],
        series: getSeries(dtype)
    }
    return option
}

var myChart = echarts.init(document.getElementById('responseDiv'));
var option;

var myChart2 = echarts.init(document.getElementById('sizeDiv'));
var socket = new WebSocket(webSocketAddress);
socket.onopen = function (event) {
    socket.onmessage = function (event) {
        var  data =  JSON.parse(event.data);
        xSer.push(data.time);
        parseData(data.data);
        myChart.setOption(draw("average"));   //参数设置方法
        myChart2.setOption(draw("size"));   //参数设置方法
    }
}