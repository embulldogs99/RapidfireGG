<?php

/* Add customizer panels, sections, settings, and controls */
function chosen_gamer_add_customizer_settings( $wp_customize ) {

  // Remove option to hide author name in byline
  $wp_customize->remove_control('author_byline');

  // categories - setting
  $wp_customize->add_setting( 'categories', array(
    'default'           => 'show',
    'sanitize_callback' => 'ct_chosen_sanitize_all_show_hide_settings'
  ) );
  // categories - control
  $wp_customize->add_control( 'categories', array(
    'label'    => __( 'Show the categories above the post title?', 'chosen-gamer' ),
    'section'  => 'chosen_additional',
    'settings' => 'categories',
    'type'     => 'radio',
    'choices'  => array(
      'show' => __( 'Show', 'chosen-gamer' ),
      'hide'  => __( 'Hide', 'chosen-gamer' )
    )
  ) );

  // comments link - setting
	$wp_customize->add_setting( 'comments_link', array(
		'default'           => 'show',
		'sanitize_callback' => 'ct_chosen_sanitize_all_show_hide_settings'
	) );
	// comments link - control
	$wp_customize->add_control( 'comments_link', array(
		'label'    => __( 'Show the comments link after the post title?', 'chosen-gamer' ),
		'section'  => 'chosen_additional',
		'settings' => 'comments_link',
		'type'     => 'radio',
		'choices'  => array(
			'show' => __( 'Show', 'chosen-gamer' ),
			'hide'  => __( 'Hide', 'chosen-gamer' )
		)
  ) );

  // author name - setting
	$wp_customize->add_setting( 'author_name', array(
		'default'           => 'show',
		'sanitize_callback' => 'ct_chosen_sanitize_all_show_hide_settings'
	) );
	// author name - control
	$wp_customize->add_control( 'author_name', array(
		'label'    => __( 'Show the author name after the post title?', 'chosen-gamer' ),
		'section'  => 'chosen_additional',
		'settings' => 'author_name',
		'type'     => 'radio',
		'choices'  => array(
			'show' => __( 'Show', 'chosen-gamer' ),
			'hide'  => __( 'Hide', 'chosen-gamer' )
		)
	) );
}
add_action( 'customize_register', 'chosen_gamer_add_customizer_settings', 20 );
