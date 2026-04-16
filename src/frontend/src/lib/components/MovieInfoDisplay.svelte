<script lang="ts">
	import { loadMovieMetadata } from '$lib/api/movie_metadata';
	import { auth } from '$lib/auth.svelte';
	import { type MovieMetadataInfo, type MovieInfo } from '$lib/types/export_types';

	let { movie }: { movie: MovieInfo } = $props();

	let metadata = $state<MovieMetadataInfo | null>(null);

	async function loadMetadata() {
		try {
			if (!auth.token) return;
			const mlist = await loadMovieMetadata(auth.token, movie.id);
			if (mlist.length == 1) {
				metadata = mlist[0];
			}
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadMetadata();
	});
</script>

{#if metadata == null || metadata.medium_image_url == ''}
	<span class="font-bold">{movie.name}</span>
{:else}
	<div>
		<img alt={movie.name} class="rounded w-full" src={metadata.medium_image_url} />
		<div class="my-1 text-center font-bold">{movie.name}</div>
	</div>
{/if}
