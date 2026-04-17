<script lang="ts">
	import { type MovieInfo } from '$lib/types/export_types';

	let { movie }: { movie: MovieInfo } = $props();
</script>

<div>
	{#if movie.medium_image_url != ''}
		<img alt={movie.name} class="w-full rounded" src={movie.medium_image_url} />
	{/if}
	<div class="my-1 text-center font-bold">{movie.name}</div>
</div>
