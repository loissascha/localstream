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

<div>
	{#each libraries as library}
		<div>
			{library.name}
		</div>
	{/each}
</div>
