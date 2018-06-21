<div <?php post_class(); ?>>
	<?php do_action( 'page_before' ); ?>
	<article>
  <div class="featured-title">
    <?php ct_chosen_featured_image(); ?>
    <div class='post-header'>
      <h1 class='post-title'><?php the_title(); ?></h1>
    </div>
  </div>
		<div class="post-content">
			<?php the_content(); ?>
			<?php wp_link_pages( array(
				'before' => '<p class="singular-pagination">' . __( 'Pages:', 'chosen-gamer' ),
				'after'  => '</p>',
			) ); ?>
			<?php do_action( 'page_after' ); ?>
		</div>
	</article>
	<?php comments_template(); ?>
</div>
