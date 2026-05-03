<script lang="ts">
	import type { MovieInfo } from '$lib/types/export_types';
	import type { Snippet } from 'svelte';
	import MovieMetadataSearchOverlay from './overlays/MovieMetadataSearchOverlay.svelte';
	import ContextMenu from './ui/ContextMenu.svelte';

	interface Props {
		children: Snippet;
		movie: MovieInfo;
	}
	let { children, movie }: Props = $props();

	let movieMetadataOverlayOpen = $state(false);
</script>

<ContextMenu closeOnItemClick={true}>
	{@render children()}
	{#snippet items(closeMenu)}
		<button
			role="menuitem"
			onclick={(e) => {
				e.preventDefault();
				e.stopPropagation();
				closeMenu();
				movieMetadataOverlayOpen = true;
			}}
			class="cursor-pointer px-4 py-2 hover:bg-neutral-700">Update Metadata</button
		>
	{/snippet}
</ContextMenu>

{#if movieMetadataOverlayOpen}
	<MovieMetadataSearchOverlay
		{movie}
		close={() => {
			movieMetadataOverlayOpen = false;
		}}
	/>
{/if}
