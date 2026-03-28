<script lang="ts">
	import { resolve } from '$app/paths';
	import {
		type ShowInfo,
		type LibraryListItem,
		type ShowListResponse
	} from '$lib/types/export_types';
	import { auth } from '$lib/auth.svelte';
	import LastWatched from '$lib/components/LastWatched.svelte';
	import { loadLibraries } from '$lib/api/libraries';

	let libraries = $state<LibraryListItem[]>([]);
	let shows = $state<ShowInfo[]>([]);
	let selectedLibrary = $state<LibraryListItem | null>(null);
	let loading = $state(true);
	let loadingShows = $state(true);
	let errorMessage = $state('');

	async function loadShows() {
		try {
			const res = await fetch('/api/shows', {
				headers: {
					Authorization: 'Bearer ' + auth.token
				}
			});
			if (!res.ok) {
				throw new Error(`Failed to load shows: ${res.status}`);
			}

			const data = (await res.json()) as ShowListResponse;
			shows = data.shows;
			console.log('shows', data);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loadingShows = false;
		}
	}

	async function loadLibrariesData() {
		try {
			if (!auth.token) {
				throw new Error('Auth token not loaded.');
			}
			const data = await loadLibraries(auth.token);
			libraries = data.libraries;
			console.log('libraries', data);

			if (selectedLibrary == null) {
				if (libraries.length > 0) {
					selectLibrary(libraries[0]);
				}
			}
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loading = false;
		}
	}

	function selectLibrary(lib: LibraryListItem | null) {
		selectedLibrary = lib;
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadLibrariesData();
		loadShows();
	});
</script>

<main>
	{#if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{/if}

	<LastWatched />

	{#if loadingShows}
		<p>Loading shows...</p>
	{:else}
		<section class="my-4">
			<h2 class="text-xl tracking-wider">Shows</h2>
			<div class="flex gap-3">
				{#each shows as show (show.id)}
					<a
						href={resolve('/(protected)/(user)/shows/[showID]', { showID: show.id })}
						class="w-60 cursor-pointer rounded-lg bg-neutral-800 p-4 hover:bg-neutral-700"
					>
						{show.name}
					</a>
				{/each}
			</div>
		</section>
	{/if}
</main>
