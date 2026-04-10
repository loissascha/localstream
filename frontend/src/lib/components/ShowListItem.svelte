<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ShowInfo } from '$lib/types/export_types';
	import ListItemA from './ListItemA.svelte';
	import ShowInfoDisplay from './ShowInfoDisplay.svelte';

	let { show }: { show: ShowInfo } = $props();
</script>

<ListItemA href={resolve('/(protected)/(user)/shows/[showID]', { showID: show.id })}>
	<ShowInfoDisplay {show} />
</ListItemA>
