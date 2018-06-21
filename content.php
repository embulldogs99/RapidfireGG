<div <?php post_class(); ?>>
	<?php do_action( 'post_before' ); ?>
	<article>
		<div class="featured-title">
			<?php ct_chosen_featured_image(); ?>
			<div class='post-header'>
				<div class="categories">
					<span><?php
						if ( get_theme_mod( 'categories' ) != 'hide' ) {
							the_category(' / ');
						} ?>
					</span>
				</div>
				<h1 class='post-title'><?php the_title(); ?></h1>
				<div class="after-post-title">
					<?php get_template_part( 'content/comments-link' ); ?>
					<?php
						if ( get_theme_mod( 'comments_link' ) != 'hide' && get_theme_mod( 'author_name' ) != 'hide' ) {
							echo ' | ';
						}
						if ( get_theme_mod( 'author_name' ) != 'hide' ) {
							echo '<span class="author">'. get_the_author() . '</span>';
						}
						?>
				</div>
			</div>
		</div>
		<div class="post-content">
			<div class="date">
				<?php
				// translator note: Date the post was PUBLISHED on
				echo esc_html__('Published', 'chosen-gamer') . ' ' . date_i18n( get_option( 'date_format' ), strtotime( get_the_date( 'r' ) ) );
				?>
			</div>
			<?php the_content(); ?>
			<?php wp_link_pages( array(
				'before' => '<p class="singular-pagination">' . __( 'Pages:', 'chosen-gamer' ),
				'after'  => '</p>',
			) ); ?>
			<?php do_action( 'post_after' ); ?>
		</div>
		<div class="post-meta">
			<div class="tags">
				<?php the_tags('', ''); ?>
			</div>
			<?php get_template_part( 'content/post-nav' ); ?>
		</div>
	</article>
	<?php comments_template(); ?>
</div>
