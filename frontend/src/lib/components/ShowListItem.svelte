<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ShowInfo } from '$lib/types/export_types';
	import ShowInfoDisplay from './ShowInfoDisplay.svelte';

	let { show }: { show: ShowInfo } = $props();
</script>

<a
	href={resolve('/(protected)/(user)/shows/[showID]', { showID: show.id })}
	class="w-60 cursor-pointer rounded-lg bg-neutral-800 p-4 hover:bg-neutral-700"
>
	<ShowInfoDisplay {show} />
</a>
