jQuery(document).ready(function($) {
  $( document ).on( 'click', '.chosen-gamer-subtitles', function () {
    $.ajax( ajaxurl,
      {
        type: 'POST',
        data: {
          action: 'dismissed_notice_handler',
          chosen_gamer_dismissed: true
        }
      } );
  } );
});
