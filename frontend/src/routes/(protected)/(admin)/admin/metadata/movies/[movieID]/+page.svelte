<script lang="ts">
	import { page } from '$app/state';
	import { loadMovieMetadata } from '$lib/api/movie_metadata';
	import { auth } from '$lib/auth.svelte';
	import type { MovieMetadataInfo } from '$lib/types/export_types';

	const movieID = $derived(page.params.movieID ?? '');

	let metadata = $state<MovieMetadataInfo[]>([]);
	let loadingMetadata = $state(true);

	async function loadMetadata() {
		try {
			if (!auth.token) return;
			metadata = await loadMovieMetadata(auth.token, movieID);
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		} finally {
			loadingMetadata = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.token) return;
		if (!auth.isAdmin) return;
		movieID;
		loadMetadata();
	});
</script>

Details Movie {movieID}

{#if loadingMetadata}
	Loading...
{:else if metadata.length == 0}
	<section id="metadata-search" class="my-8">
		<h2 class="text-2xl font-bold">Search Metadata</h2>
	</section>
{/if}
