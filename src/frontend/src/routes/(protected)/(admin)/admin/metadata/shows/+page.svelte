<script lang="ts">
	import AdminShowMetadataBlock from '$lib/components/admin/AdminShowMetadataBlock.svelte';
	import { shows } from '$lib/shows.svelte';

	var hideSingle = $state(true);
</script>

<div class="mb-4">
	<input id="hidesingle" bind:checked={hideSingle} class="cursor-pointer" type="checkbox" />
	<label for="hidesingle" class="cursor-pointer"> Hide items with single metadata</label>
</div>

<section class="grid grid-cols-1 gap-4">
	{#each shows.shows as show (show.id)}
		{#if !hideSingle || show.fetch_source == 'multiple' || show.fetch_source == 'empty'}
			<AdminShowMetadataBlock {show} />
		{/if}
	{/each}
</section>
