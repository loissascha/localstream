<script lang="ts">
	import type { CollectionInfo } from '$lib/types/export_types';

	let { collection }: { collection: CollectionInfo } = $props();

	let images = $state<string[]>([]);
	let gridCols = $derived(images.length <= 1 ? 'grid-cols-1' : 'grid-cols-2');

	$effect(() => {
		if (!collection) return;

		let allImages = [];
		for (const movie of collection.movies) {
			if (movie.medium_image_url) {
				allImages.push(movie.medium_image_url);
			}
		}
		for (const show of collection.shows) {
			if (show.medium_image_url) {
				allImages.push(show.medium_image_url);
			}
		}
		images = allImages.slice(0, 4);
	});
</script>

<div class={`grid ${gridCols}`}>
	{#each images as image (image)}
		<img src={image} alt={image} class="" />
	{/each}
	{#if images.length == 0}
		<div class="aspect-video w-full bg-neutral-800"></div>
	{/if}
</div>
