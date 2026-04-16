<script lang="ts">
	import { resolve } from '$app/paths';
	import { loadLibraries } from '$lib/api/libraries';
	import { auth } from '$lib/auth.svelte';
	import type { LibraryListItem } from '$lib/types/export_types';

	let libraries = $state<LibraryListItem[]>([]);

	async function loadLibs() {
		if (!auth.token) {
			throw new Error('no auth token');
		}
		const data = await loadLibraries(auth.token);
		libraries = data.libraries;
	}

	$effect(() => {
		if (!auth.token) return;
		loadLibs();
	});
</script>

<section id="stats_and_actions" class="mb-4 flex justify-center items-center gap-4">
	<span>{libraries.length} Libraries</span>
	<a
		href={resolve('/(protected)/(admin)/admin/libraries/new')}
		class="rounded-full bg-neutral-800 hover:bg-neutral-700 px-4 py-2">+ New Library</a
	>
</section>

<section id="libraries" class="flex justify-center gap-2">
	{#each libraries as library (library.id)}
		<div
			class="flex h-40 w-80 flex-col items-end justify-between rounded border border-neutral-500 p-4"
		>
			<div class="grow">Settings</div>
			<div class="flex w-full flex-col items-center">
				<span class="font-bold">
					{library.name}
				</span>
				<span>
					{library.library_type}
				</span>
				<span>
					{library.path}
				</span>
			</div>
		</div>
	{/each}
</section>
