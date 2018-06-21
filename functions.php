<?php

//----------------------------------------------------------------------------------
// Include files
//----------------------------------------------------------------------------------
require_once( trailingslashit( get_stylesheet_directory() ) . 'inc/customizer.php' );

//----------------------------------------------------------------------------------
// Setup Chosen Gamer
//----------------------------------------------------------------------------------
function chosen_gamer_setup_theme() {
  load_child_theme_textdomain( 'chosen-gamer', get_stylesheet_directory() . '/languages' );
}
add_action( 'after_setup_theme', 'chosen_gamer_setup_theme' );

//----------------------------------------------------------------------------------
// Load scripts and stylesheets for the front-end
//----------------------------------------------------------------------------------
function chosen_gamer_enqueue_styles() {
  // load Chosen stylesheet
  $parent_style = 'ct-chosen-style';
  wp_enqueue_style( $parent_style, get_template_directory_uri() . '/style.css' );
  // load stylesheet after Chosen's stylesheet (uses it as a dependency)
  wp_enqueue_style( 'ct-chosen-gamer-style',
    get_stylesheet_directory_uri() . '/style.css',
    array( $parent_style )
  );
  // load Google Fonts (Montserrate & PT Serif)
  $font_args = array(
		'family' => str_replace( '%2B', '+', urlencode( 'Montserrat:400:700|PT+Serif:400,400i' ))
	);
	$fonts_url = add_query_arg( $font_args, '//fonts.googleapis.com/css' );
  wp_enqueue_style( 'ct-chosen-gamer-google-fonts', $fonts_url );
}
add_action( 'wp_enqueue_scripts', 'chosen_gamer_enqueue_styles' );

//----------------------------------------------------------------------------------
// Load scripts and stylesheets on the back-end
//----------------------------------------------------------------------------------
function chosen_gamer_admin_enqueue_scripts() {
  wp_enqueue_script( 'chosen-gamer-admin-js', trailingslashit(get_stylesheet_directory_uri()) . 'js/admin.js', array( 'jquery' ), '', true );
}
add_action( 'admin_enqueue_scripts', 'chosen_gamer_admin_enqueue_scripts' );

//----------------------------------------------------------------------------------
// Dequeue Google Fonts loaded by Chosen
//----------------------------------------------------------------------------------
function chosen_gamer_dequeue_parent_fonts() {
  wp_dequeue_style( 'ct-chosen-google-fonts' );
}
add_action( 'wp_enqueue_scripts', 'chosen_gamer_dequeue_parent_fonts', 999 );

//----------------------------------------------------------------------------------
// Add classes to body element
//----------------------------------------------------------------------------------
function chosen_gamer_add_body_classes($classes) {

  $search_bar = get_theme_mod( 'search_bar' );
  $extra_wide_post = get_theme_mod( 'full_width_post' );

  if ( $search_bar == 'show' ) {
    $classes[] = 'search-bar';
  }
  if ( $extra_wide_post != 'no' ) {
    $classes[] = 'extra-wide-post';
  }

  return $classes;
}
add_action( 'body_class', 'chosen_gamer_add_body_classes' );

//----------------------------------------------------------------------------------
// Filter the footer text
//----------------------------------------------------------------------------------
function chosen_gamer_filter_footer_text($footer_text) {
  $footer_text = sprintf( __( '<a href="%s">Chosen Gamer Theme</a> by Compete Themes.', 'chosen-gamer' ), 'https://www.competethemes.com/chosen-gamer/' );
  return $footer_text;
}
add_filter( 'ct_chosen_footer_text' , 'chosen_gamer_filter_footer_text' );

//----------------------------------------------------------------------------------
// Remove subtitles on archive pages
//----------------------------------------------------------------------------------
function chosen_gamer_no_subtitles_archives() {
  global $wp_query;
  if ( is_singular() || (is_home() && $wp_query->current_post == 0 && get_theme_mod('full_width_post') != 'no') ) {
    return true;
  } else {
    return false;
  }
}
add_filter( 'subtitle_exists', 'chosen_gamer_no_subtitles_archives' );

//----------------------------------------------------------------------------------
// Recommend the Subtitles plugin in an admin notice
//----------------------------------------------------------------------------------
function chosen_gamer_recommend_subtitle_plugin() {
  if ( !get_option('chosen-gamer-notice-dismissed') ) {
    $plugins = get_plugins();
    if ( !array_key_exists( 'subtitles/subtitles.php', $plugins ) ) {
      $plugin_search = add_query_arg( array(
        'tab' => 'search',
        's'   => 'subtitles'
      ), admin_url( 'plugin-install.php' ) );
      ?>
      <div class="notice notice-info is-dismissible chosen-gamer-subtitles">
          <p><?php printf( __( 'Please install the <a href="%s">Subtitles</a> plugin to display subtitles on posts and pages.', 'chosen-gamer' ), $plugin_search); ?></p>
      </div>
      <?php
    }
  }
}
add_action( 'admin_notices', 'chosen_gamer_recommend_subtitle_plugin' );

//----------------------------------------------------------------------------------
// Ajax handler for permanently dismissing admin notice
//----------------------------------------------------------------------------------
function chosen_gamer_ajax_notice_handler() {
  if ( $_POST['chosen_gamer_dismissed'] ) {
    update_option( 'chosen-gamer-notice-dismissed', true );
  }
}
add_action( 'wp_ajax_dismissed_notice_handler', 'chosen_gamer_ajax_notice_handler' );

//----------------------------------------------------------------------------------
// Add Chosen Gamer's new Customizer settings to Chosen's reset option
//----------------------------------------------------------------------------------
function chosen_gamer_reset_customizer_settings($mods_array) {
  array_push($mods_array, 'categories', 'comments_link', 'author_name');
  return $mods_array;
}
add_filter( 'ct_chosen_mods_to_remove', 'chosen_gamer_reset_customizer_settings' );
