<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth.svelte';
	import { type MovieInfo } from '$lib/types/export_types';
	import DropdownItem from './ui/DropdownItem.svelte';
	import PercentageBar from './ui/PercentageBar.svelte';
	import { Popover } from 'melt/builders';

	let { movie, nameLink = false }: { movie: MovieInfo; nameLink?: boolean } = $props();

	const popover = new Popover();
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	oncontextmenu={(e) => {
		e.preventDefault();
		popover.open = true;
	}}
>
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
		class={`my-1 w-full cursor-pointer font-bold text-neutral-200 ${nameLink ? 'hover:text-white' : ''}`}
	>
		<div class="w-full text-center">{movie.name}</div>
		{#if movie.year > 0}
			<div class="text-center text-sm text-neutral-400">{movie.year}</div>
		{/if}
	</button>
	<div {...popover.trigger}></div>
	<div
		{...popover.content}
		class="rounded-md border border-neutral-500 bg-neutral-800 text-white shadow-lg"
	>
		{#if auth.isAdmin}
			<DropdownItem
				onclick={() => {
					alert('Item 1');
				}}>Metadata</DropdownItem
			>
		{/if}
	</div>
</div>
