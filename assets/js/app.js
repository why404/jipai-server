$(window).load(function() {

  if($('.home-page').size()){
    // The slider being synced must be initialized first
    $('#carousel').flexslider({
      animation: "fade",
      controlNav: false,
      directionNav: false,
      animationLoop: true,
      asNavFor: '#slider'
    });

    $('#slider').flexslider({
      animation: "fade",
      controlNav: false,
      directionNav: false,
      animationLoop: true,
      sync: "#carousel"
    });
  }
});


$(function(){

  // player page
  if($('.player-page').size()){

    // videojs.options.flash.swf = "http://pili-static.qiniudn.com/video.js/4.8.4/video-js.swf"
    window.player = videojs('video');

    var initializeVideoPlayerSize = function(){
      var $playerWrap = $('#player-wrap');
          $playerWrap.height($playerWrap.width() * 9 / 16)
    }

    initializeVideoPlayerSize();

    $(window).bind('resize', initializeVideoPlayerSize);

    $('body').on('click', '#player-button', function(){
      $('#content').show();
      window.player.play();
      window.player.requestFullscreen();
    });

    var drawWave = function(){
      if (window.SW) {
        window.SW.stop();
        window.SW._clear();
        $('#wave canvas').remove();
      }
      window.SW = new SiriWave({
        container: document.getElementById('wave'),
        width: window.innerWidth,
        height: 120,
        color: '#FF6C48',
        frequency: 345,
        speed: 0.02,
        amplitude: 0.1,
        autostart: true,
      });
    }

    setTimeout(function(){$('#cover, #player-button').show();}, 1000);
    setTimeout(function(){
      $('#soundwave').show();
      $('#wave').fadeIn();
    }, 1000 * 2);

    drawWave();

    $(window).bind('resize', function(){
      drawWave();
    });
  }


});