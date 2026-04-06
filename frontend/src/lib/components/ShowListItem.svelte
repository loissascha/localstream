<script lang="ts">
	import { resolve } from '$app/paths';
	import { loadShowMetadata } from '$lib/api/show_metadata';
	import { auth } from '$lib/auth.svelte';
	import { type ShowMetadataInfo, type ShowInfo } from '$lib/types/export_types';

	let { show }: { show: ShowInfo } = $props();

	let metadata = $state<ShowMetadataInfo | null>(null);

	async function loadMetadata() {
		try {
			if (!auth.token) return;
			const mlist = await loadShowMetadata(auth.token, show.id);
			if (mlist.length == 1) {
				metadata = mlist[0];
			}
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadMetadata();
	});
</script>

<a
	href={resolve('/(protected)/(user)/shows/[showID]', { showID: show.id })}
	class="w-60 cursor-pointer rounded-lg bg-neutral-800 p-4 hover:bg-neutral-700"
>
	{#if metadata == null || metadata.medium_image_url == ''}
		{show.name}
	{:else}
		<img alt={show.name} src={metadata.medium_image_url} />
	{/if}
</a>
