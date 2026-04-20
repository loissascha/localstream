<script lang="ts">
	import { resolve } from '$app/paths';
	import { setPrimaryMetadataForMovie } from '$lib/api/admin/movie_metadata';
	import { loadMovieMetadata } from '$lib/api/movie_metadata';
	import { auth } from '$lib/auth.svelte';
	import type { MovieInfo, MovieMetadataInfo } from '$lib/types/export_types';
	import MovieMetadataSearchOverlay from '../overlays/MovieMetadataSearchOverlay.svelte';

	let { movie }: { movie: MovieInfo } = $props();
	let metadata = $state<MovieMetadataInfo[]>([]);
	let loading = $state(true);

	var showMovieMetadataSearchOverlay = $state(false);
	var movieMetadataSearchOverlayMovieID = $state('');

	async function loadMetadata() {
		try {
			if (!auth.token) return;
			metadata = await loadMovieMetadata(auth.token, movie.id);
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadMetadata();
	});
</script>

<div class="rounded bg-neutral-800 p-4">
	<div class="font-bold">
		{movie.name}
	</div>
	{#if loading}
		Loading metadata...
	{:else}
		<div class="flex items-center justify-between">
			<span>
				Metadata: {metadata.length}
			</span>
			<button
				class="cursor-pointer"
				onclick={() => {
					showMovieMetadataSearchOverlay = true;
					movieMetadataSearchOverlayMovieID = movie.id;
				}}
			>
				Details</button
			>
		</div>
		<div class="mt-4 flex flex-col gap-2">
			{#each metadata as m (m.id)}
				<div class="">
					<div class="font-bold">{m.id}</div>
					<div class="font-bold">{m.name}</div>
					<div class="grid grid-cols-2">
						<div>
							<p>{m.description}</p>
							{#if metadata.length > 1}
								<button
									onclick={() => {
										if (!auth.token) return;
										setPrimaryMetadataForMovie(auth.token, movie.id, m.id)
											.then(() => {
												loadMetadata();
											})
											.catch((e) => {
												const m = (e as Error).message;
												alert(m);
											});
									}}
									class="mt-2 cursor-pointer rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600"
									>Select as Primary</button
								>
							{/if}
						</div>
						<div>
							<img class="w-full" src={m.medium_image_url} alt={m.name} />
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

{#if showMovieMetadataSearchOverlay}
	<MovieMetadataSearchOverlay
		close={() => {
			showMovieMetadataSearchOverlay = false;
			loadMetadata();
		}}
		{movie}
	/>
{/if}
