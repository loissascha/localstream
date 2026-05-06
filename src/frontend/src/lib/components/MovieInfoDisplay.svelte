<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { type MovieInfo } from '$lib/types/export_types';
	import PercentageBar from './ui/PercentageBar.svelte';

	let { movie, nameLink = false }: { movie: MovieInfo; nameLink?: boolean } = $props();
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div>
	{#if movie.medium_image_url != ''}
		<img alt={movie.name} class="w-full rounded" src={movie.medium_image_url} />
	{/if}
	{#if movie.percentage > 0}
		<div>
			<PercentageBar percentage={movie.percentage} />
		</div>
	{/if}
	<button
		onclick={(e) => {
			if (nameLink) {
				e.preventDefault();
				e.stopPropagation();
				goto(resolve('/(protected)/(user)/movies/[movieID]', { movieID: movie.id }));
			}
		}}
		class={`my-1 w-full max-w-full cursor-pointer truncate font-bold text-neutral-200 ${nameLink ? 'hover:text-white' : ''}`}
	>
		<div class="w-full text-center">{movie.name}</div>
		{#if movie.year > 0}
			<div class="text-center text-sm text-neutral-400">{movie.year}</div>
		{/if}
	</button>
</div>
