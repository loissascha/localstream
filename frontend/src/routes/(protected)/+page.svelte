<script lang="ts">
	import { resolve } from '$app/paths';
	import {
		type ShowInfo,
		type LibraryListItem,
		type LibraryListResponse,
		type ShowListResponse
	} from '$lib/types/export_types';
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';

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
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		loadLibraries();
		loadShows();
	});
</script>

<main class="min-h-dvh px-5 py-5">
	<header class="mb-5 flex items-start justify-between gap-4">
		<div>
			<h1 class="m-0 text-[clamp(1.6rem,2.8vw,2.4rem)] font-semibold">Localstream</h1>
			<p class="mt-1.5">Your local MP4 library</p>
		</div>
		<a
			href={resolve('/logout')}
			type="submit"
			class="cursor-pointer rounded-md border border-neutral-500 bg-neutral-600 px-3 py-1.5 text-sm hover:bg-neutral-500"
		>
			Log out
		</a>
	</header>

	{#if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{/if}

	{#if loading}
		<p>Loading video library...</p>
	{:else if libraries.length === 0}
		<p>No libraries found.</p>
	{:else}
		<section class="flex gap-3">
			{#each libraries as library (library.id)}
				<button
					onclick={() => selectLibrary(library)}
					class="cursor-pointer rounded border border-neutral-600 bg-neutral-800 p-3"
				>
					{#if selectedLibrary?.id == library.id}
						<strong>{library.name}</strong>
					{:else}
						{library.name}
					{/if}
					<br />
					{library.library_type}
				</button>
			{/each}
		</section>
	{/if}

	{#if loadingShows}
		<p>Loading shows...</p>
	{:else}
		<section class="my-4 flex gap-3">
			{#each shows as show (show.id)}
				<a
					href={resolve('/(protected)/shows/[showID]', { showID: show.id })}
					class="w-60 cursor-pointer rounded-lg border border-blue-500 bg-blue-800 p-4 shadow-lg shadow-blue-600/50"
				>
					{show.name}
				</a>
			{/each}
		</section>
	{/if}
</main>
