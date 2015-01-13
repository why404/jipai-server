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

    var initializeVideoPlayerSize = function(){
      var $playerWrap = $('#player-wrap');
          $playerWrap.height($playerWrap.width() * 9 / 16)
    }

    initializeVideoPlayerSize();

    $(window).bind('resize', initializeVideoPlayerSize);


    $('body').on('click', '#player-button', function(){
      $('#content').show();
      initializeVideoPlayerSize();

      var $videoSource = $('#video source');

      if( $videoSource.attr('type') == "rtmp/mp4" ){
        SewisePlayer.setup({
          server: "live",
          type: "rtmp",
          buffer: 1,
          streamurl: $videoSource.attr('src'),
          autostart: "true",
          skin: "liveWhite",
          claritybutton: "disable",
          timedisplay: "disable",
          topBarDisplay: "disable",
          playerName: "Jipai Player",
          copyright: "(C) All right reserved the Jipai.in",
          logo: "/assets/images/headline.png"
        }, "player-wrap");

      }else{

        SewisePlayer.setup({
          server: "vod",
          type: "m3u8",
          videourl: $videoSource.attr('src'),
          autostart: "true",
          skin: "vodWhite",
          claritybutton: "disable",
          timedisplay: "disable",
          topBarDisplay: "disable",
          playerName: "Jipai Player",
          copyright: "(C) All right reserved the Jipai.in",
          lang: "zh_CN",
          logo: "/assets/images/headline.png"
        }, "player-wrap");

      }

      $('#soundwave').hide();
      $('#player-button').hide();
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