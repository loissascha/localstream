<script lang="ts">
	import { resolve } from '$app/paths';
	import {
		type LibraryListItem,
		type LibraryListResponse,
		type VideoListItem,
		type VideoListResponse
	} from '$lib/types/export_types';
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';

	let videos = $state<VideoListItem[]>([]);
	let libraries = $state<LibraryListItem[]>([]);
	let selectedLibraryID = $state<string | null>(null);
	let loading = $state(true);
	let errorMessage = $state('');

	const toHumanSize = (bytes: number): string => {
		if (bytes < 1024) {
			return `${bytes} B`;
		}

		const units = ['KB', 'MB', 'GB', 'TB'];
		let size = bytes / 1024;
		let unitIndex = 0;

		while (size >= 1024 && unitIndex < units.length - 1) {
			size /= 1024;
			unitIndex += 1;
		}

		return `${size.toFixed(1)} ${units[unitIndex]}`;
	};

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

			if (selectedLibraryID == null) {
				if (libraries.length > 0) {
					selectLibraryID(libraries[0].id);
				}
			}
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loading = false;
		}
	}

	function selectLibraryID(id: string | null) {
		selectedLibraryID = id;
	}

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		// loadVideos();
		loadLibraries();
	});
</script>

<main
	class="min-h-dvh bg-[radial-gradient(circle_at_15%_10%,#d7e8ff_0%,transparent_40%),radial-gradient(circle_at_85%_0%,#c8f5e9_0%,transparent_32%),linear-gradient(180deg,#eef2f7_0%,#dce6f2_100%)] px-5 py-5 text-slate-900"
>
	<header class="mb-5 flex items-start justify-between gap-4">
		<div>
			<h1 class="m-0 text-[clamp(1.6rem,2.8vw,2.4rem)] font-semibold">Localstream</h1>
			<p class="mt-1.5 text-slate-700">Your local MP4 library</p>
		</div>
		<a
			href={resolve('/logout')}
			type="submit"
			class="cursor-pointer rounded-md border border-slate-900/20 bg-white/70 px-3 py-1.5 text-sm text-slate-800"
		>
			Log out
		</a>
	</header>

	{#if loading}
		<p>Loading video library...</p>
	{:else if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{:else if libraries.length === 0}
		<p>No libraries found.</p>
	{:else}
		<section class="flex gap-3">
			{#each libraries as library (library.id)}
				<button
					onclick={() => selectLibraryID(library.id)}
					class="cursor-pointer rounded border border-neutral-400 bg-neutral-300 p-3 shadow"
				>
					{#if selectedLibraryID == library.id}
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
</main>
