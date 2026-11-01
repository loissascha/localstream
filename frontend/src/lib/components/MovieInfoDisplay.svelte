<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { type MovieInfo } from '$lib/types/export_types';
	import PercentageBar from './ui/PercentageBar.svelte';

	let { movie, nameLink = false }: { movie: MovieInfo; nameLink?: boolean } = $props();
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="transition-all duration-200 group-hover:scale-105">
	{#if movie.medium_image_url != ''}
		<img alt={movie.name} class="w-full rounded" src={movie.medium_image_url} />
	{/if}
	{#if movie.percentage > 0}
		<div>
			<PercentageBar percentage={movie.percentage} />
		</div>
	{/if}
	{#if movie.had_warnings || movie.had_errors || !movie.playable}
		<div class="absolute top-2 right-2">
			{#if movie.had_warnings}
				<div class="rounded-full bg-yellow-300 px-4 py-2 text-xl font-bold text-black">!</div>
			{:else if movie.had_errors}
				<div class="rounded-full bg-red-500 px-4 py-2 text-xl font-bold">!</div>
			{:else if !movie.playable}
				<div class="rounded-full bg-red-500 px-4 py-2 text-xl font-bold">X</div>
			{/if}
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
		class={`mt-1 w-full max-w-full cursor-pointer truncate font-bold text-neutral-200 ${nameLink ? 'hover:text-white' : ''}`}
	>
		<div class="w-full text-center">{movie.name}</div>
		{#if movie.year > 0}
			<div class="text-center text-sm text-neutral-400">{movie.year}</div>
		{/if}
	</button>
</div>
