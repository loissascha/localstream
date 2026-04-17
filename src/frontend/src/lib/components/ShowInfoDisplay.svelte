<script lang="ts">
	import { type ShowInfo } from '$lib/types/export_types';

	let { show }: { show: ShowInfo } = $props();
</script>

<div>
	{#if show.medium_image_url != ''}
		<img alt={show.name} class="w-full rounded" src={show.medium_image_url} />
	{/if}
	<div class="my-1 text-center font-bold">{show.name}</div>
</div>
