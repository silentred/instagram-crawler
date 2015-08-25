$(function(){

var userId = $("#userId").text();
var doneCnt = 0;
var cntDom = $("#doneCnt");

var ws;

$("#wsConn").click(function(){
    ws = new WebSocket("ws://50.117.7.122:3000/ws?userId="+userId);

    ws.addEventListener("message", function(e){
        //console.log("recieve message...");
    
        var obj = $.parseJSON(e.data);
        //console.log(obj);
        var cnt = parseInt(obj.data);

        // need to be precise. cnt is NaN
        doneCnt += 1;
        cntDom.text(doneCnt);
    })
    
    ws.addEventListener("open", function(e){
        console.log(e)
        // @TODO alert if the connection is established
    })
});

$("#wsSend").click(function(){
    var text = $("#userIds").val();
    var aIds = text.trim().split(/[\., -]/);
    var numRegex = /\d+/;
    aIds = aIds.map(function(item){
        if (numRegex.test(item)) {return item};
    }); 
    var data = aIds.join(',');

    var msg = {action: "start", data: text};

    console.log("sending the message");
    ws.send(JSON.stringify(msg));

    // empty the textArea
    $("#userIds").val("");
})

$("#wsQuit").click(function(){
    var msg = {action: "stop"};
    console.log("stoping the message");
    ws.send(JSON.stringify(msg));
});

$("#queryUserIdBtn").click(function(){
    var url = $("#homePage").val();
    $.ajax({
        url:'/oauth/ajaxGetUserId',
        data:{    
            url : url
        },    
        type:'get',    
        cache:false,    
        dataType:'json',    
        success:function(data) {    
            $("#resultUserId").val(data.userId)
        },    
        error : function(e) {    
             console.log(e);    
        }    

    });
    return false;
});


});

