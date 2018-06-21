<?php if ( get_theme_mod( 'comments_link' ) != 'hide' ) : ?>
  <span class="comments-link">
    <i class="fa fa-comment" title="<?php esc_attr_e( 'comment icon', 'chosen-gamer' ); ?>" aria-hidden="true"></i>
    <?php
    if ( ! comments_open() && get_comments_number() < 1 ) :
      comments_number( __( 'Comments closed', 'chosen-gamer' ), __( '1 Comment', 'chosen-gamer' ), _x( '% Comments', 'noun: 5 comments', 'chosen-gamer' ) );
    else :
      echo '<a href="' . esc_url( get_comments_link() ) . '">';
      comments_number( __( 'Leave a Comment', 'chosen-gamer' ), __( '1 Comment', 'chosen-gamer' ), _x( '% Comments', 'noun: 5 comments', 'chosen-gamer' ) );
      echo '</a>';
    endif;
    ?>
  </span>
<?php endif;
