<script lang="ts">
	import { resolve } from '$app/paths';
	import {
		type ShowInfo,
		type LibraryListItem,
		type LibraryListResponse,
		type ShowListResponse
	} from '$lib/types/export_types';
	import { auth } from '$lib/auth.svelte';
	import LastWatched from '$lib/components/LastWatched.svelte';

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

	async function loadLibraries() {
		try {
			const res = await fetch('/api/libraries', {
				headers: {
					Authorization: 'Bearer ' + auth.token
				}
			});
			if (!res.ok) {
				throw new Error(`Failed to load videos: ${res.status}`);
			}

			const data = (await res.json()) as LibraryListResponse;
			libraries = data.libraries;
			console.log('data', data);

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
		loadLibraries();
		loadShows();
	});
</script>

<main class="min-h-dvh px-5 py-5">
	{#if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{/if}

	<LastWatched />

	{#if loadingShows}
		<p>Loading shows...</p>
	{:else}
		<section class="my-4">
			<h2>Shows</h2>
			<div class="flex gap-3">
				{#each shows as show (show.id)}
					<a
						href={resolve('/(protected)/(user)/shows/[showID]', { showID: show.id })}
						class="w-60 cursor-pointer rounded-lg border border-blue-500 bg-blue-800 p-4 shadow-lg shadow-blue-600/50"
					>
						{show.name}
					</a>
				{/each}
			</div>
		</section>
	{/if}
</main>
