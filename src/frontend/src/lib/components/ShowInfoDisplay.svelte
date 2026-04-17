<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { type ShowInfo } from '$lib/types/export_types';

	let { show, nameLink = false }: { show: ShowInfo; nameLink?: boolean } = $props();
</script>

<div>
	{#if show.medium_image_url != ''}
		<img alt={show.name} class="w-full rounded" src={show.medium_image_url} />
	{/if}
	<button
		onclick={(e) => {
			if (nameLink) {
				e.preventDefault();
				e.stopPropagation();
				goto(resolve('/(protected)/(user)/shows/[showID]', { showID: show.id }));
			}
		}}
		class={`my-1 cursor-pointer text-center font-bold text-neutral-200 ${nameLink ? 'hover:text-white' : ''}`}
	>
		{show.name}
	</button>
</div>
