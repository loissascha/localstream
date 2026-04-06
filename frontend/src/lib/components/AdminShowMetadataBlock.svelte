<script lang="ts">
	import { loadShowMetadata } from '$lib/api/show_metadata';
	import { auth } from '$lib/auth.svelte';
	import { type ShowMetadataInfo, type ShowInfo } from '$lib/types/export_types';

	let { show }: { show: ShowInfo } = $props();
	let metadata = $state<ShowMetadataInfo[]>([]);
	let loading = $state(true);

	async function loadMetadata(show: ShowInfo) {
		try {
			if (!auth.token) return;
			metadata = await loadShowMetadata(auth.token, show.id);
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadMetadata(show);
	});
</script>

<div class="border border-neutral-500 p-4">
	<div>
		{show.name}
	</div>
	{#if loading}
		Loading metadata...
	{:else}
		<div>
			Metadata: {metadata.length}
		</div>
	{/if}
</div>
