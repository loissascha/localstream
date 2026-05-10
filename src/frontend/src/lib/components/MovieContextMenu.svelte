<script lang="ts">
	import type { MovieInfo } from '$lib/types/export_types';
	import type { Snippet } from 'svelte';
	import MovieMetadataSearchOverlay from './overlays/MovieMetadataSearchOverlay.svelte';
	import ContextMenu from './ui/ContextMenu.svelte';
	import { deleteWatchstateMovie, setWatchstateFinishedMovie } from '$lib/api/watchstate_movie';
	import { auth } from '$lib/auth.svelte';
	import { reloadSingleMovie } from '$lib/movies.svelte';

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
		{#if auth.isAdmin}
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
		{/if}
		{#if movie.finished}
			<button
				role="menuitem"
				onclick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					if (auth.token) {
						deleteWatchstateMovie(auth.token, movie.id).then(() => {
							// loadMoviesDatabase();
							reloadSingleMovie(movie.id);
						});
					}
					closeMenu();
				}}
				class="cursor-pointer px-4 py-2 hover:bg-neutral-700">Set Not Watched</button
			>
		{:else}
			<button
				role="menuitem"
				onclick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					if (auth.token) {
						setWatchstateFinishedMovie(auth.token, movie.id).then(() => {
							// loadMoviesDatabase();
							reloadSingleMovie(movie.id);
						});
					}
					closeMenu();
				}}
				class="cursor-pointer px-4 py-2 hover:bg-neutral-700">Set Watched</button
			>
		{/if}
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
