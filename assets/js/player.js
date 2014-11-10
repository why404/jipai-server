var start = function(){
  var width = $('body').width();
  var height = $('body').height()

  var canvas = document.getElementById('cover');

  canvas.setAttribute('width', width);
  canvas.setAttribute('height', height);

  var ctx = canvas.getContext('2d');

  var radius = 10;

  ctx.beginPath();
  ctx.fillStyle = "#FFE1DB";

  ctx.arc(canvas.width/2, canvas.height/2, radius, 0, Math.PI*2, true);
  ctx.fill();

  var interval = setInterval(function(){
    radius += 20
    ctx.arc(canvas.width/2, canvas.height/2, radius, 0, Math.PI*2, true);
    ctx.fill();

    if(radius > Math.max(width, height)){
      clearInterval(interval);
    }
  }, 16);
}

$(window).load(function(){
});

$(function(){
  setTimeout(start, 500);
  setTimeout(function(){
    $('body').append('<img class="player-start" src="images/player-start-2.gif" />');
    setTimeout(function(){
      $('.player-halo').show();
    }, 1000 * 1.5);
  }, 750);
  // start();

  // $(window).bind('resize', function(){
  //   start();
  // });

});