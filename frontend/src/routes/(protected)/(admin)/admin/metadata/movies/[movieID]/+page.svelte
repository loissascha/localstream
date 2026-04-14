<script lang="ts">
	import { page } from '$app/state';
	import { loadMovieMetadata } from '$lib/api/movie_metadata';
	import { getMovie } from '$lib/api/movies';
	import { auth } from '$lib/auth.svelte';
	import SearchIcon from '$lib/icons/SearchIcon.svelte';
	import { type MovieInfo, type MovieMetadataInfo } from '$lib/types/export_types';

	const movieID = $derived(page.params.movieID ?? '');

	let metadata = $state<MovieMetadataInfo[]>([]);
	let movie = $state<MovieInfo | null>(null);
	let loadingMovie = $state(true);
	let loadingMetadata = $state(true);

	async function loadMovieData() {
		try {
			if (!auth.token) return;
			movie = await getMovie(auth.token, movieID);
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		} finally {
			loadingMovie = false;
		}
	}

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
		loadMovieData();
		loadMetadata();
	});

	async function submitMetadataSearchForm() {}
</script>

<h1 class="text-2xl font-bold">{movie?.name}</h1>

{#if loadingMetadata}
	Loading...
{:else if metadata.length == 0}
	<section id="metadata-search" class="my-8">
		<h2 class="text-xl font-bold">Search Metadata</h2>
		<form
			onsubmit={(e) => {
				e.preventDefault();
				submitMetadataSearchForm();
			}}
			class="flex items-center gap-2"
		>
			<input
				type="text"
				class="my-4 w-full rounded bg-neutral-800 px-4 py-2"
				placeholder="Search by Name"
			/>
			<button
				class="flex cursor-pointer gap-2 rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600"
			>
				<SearchIcon /> Search
			</button>
		</form>
	</section>
{/if}
