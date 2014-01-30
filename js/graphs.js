
var data = [{x:0.0, y:0.0}],
    maxLen = 100;

var graph = new Rickshaw.Graph( {
        element: document.querySelector("#chart"),
        renderer: 'line',
        width: 540,
        height: 240,
        min: 'auto',
        series: [ {
                color: 'steelblue',
                data: data
        } ]
} );

graph.render();

var goSrv;
var first = true;
var onMessage = function (event) {
  var tuple = JSON.parse(event.data);

  data.push(tuple);

  for (; data.length >= maxLen; data.shift()){}

  graph.update();
};

var onOpen = function (event) {
  console.log("Lol!" + event);
};

// Very aggressively reconnect!
var onClose = function(event) {
  init();
};

var init = function() {
  goSrv = new WebSocket("ws://127.0.0.1:8080/ws");
  goSrv.onopen = onOpen
  goSrv.onmessage = onMessage;
  goSrv.onclose = onClose;
};
init();
