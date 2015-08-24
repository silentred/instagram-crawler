$(function(){

var userId = $("#userId").text();


var ws;

$("#wsConn").click(function(){
    ws = new WebSocket("ws://50.117.7.122:3000/ws?userId="+userId);

    ws.addEventListener("message", function(e){
        console.log("recieve message...");
        console.log(e.data);
    })
    
    ws.addEventListener("open", function(e){
        console.log(e)
    })
});

$("#wsSend").click(function(){
    var msg = {action: "start"};
    console.log("sending the message");
    ws.send(JSON.stringify(msg));
})


});

