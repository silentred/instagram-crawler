$(function(){

var userId = $("#userId").text();
var doneCnt = 0;
var cntDom = $("#doneCnt");

var ws;

$("#wsConn").click(function(){
    ws = new WebSocket("ws://50.117.7.122:3000/ws?userId="+userId);

    ws.addEventListener("message", function(e){
        //console.log("recieve message...");
        console.log(e);

        var obj = $.parseJSON(e.data);
        var cnt = parseInt(obj.data);
        doneCnt += cnt;
        cntDom.text(doneCnt);
    })
    
    ws.addEventListener("open", function(e){
        console.log(e)
    })
});

$("#wsSend").click(function(){
    var text = $("#userIds").text();
    var aIds = text.trim().split(/[., -]/);
    var numRegex = /d+/;
    aIds = aIds.map(function(item){
        if (numRegex.test(item)) {return item};
    });
    data = aIds.join(',');

    var msg = {action: "start", data: data};

    console.log("sending the message");
    ws.send(JSON.stringify(msg));
})

$("#wsQuit").click(function(){
    var msg = {action: "stop"};
    console.log("stoping the message");
    ws.send(JSON.stringify(msg));
});


});

