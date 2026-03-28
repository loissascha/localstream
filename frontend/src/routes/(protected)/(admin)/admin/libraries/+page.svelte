<script lang="ts">
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

<section id="stats_and_actions" class="mb-4 flex justify-center gap-2">
	<span>{libraries.length} Libraries</span>
	<button>+ New Library</button>
</section>

<section id="libraries" class="flex justify-center gap-2">
	{#each libraries as library}
		<div
			class="flex h-40 w-80 flex-col items-end justify-between rounded border border-neutral-500 p-4"
		>
			<div class="grow">Settings</div>
			<div class="flex flex-col items-center w-full">
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
