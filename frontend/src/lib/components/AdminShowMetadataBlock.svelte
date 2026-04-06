<script lang="ts">
	import type { ShowInfo } from '$lib/types/export_types';

	let { show }: { show: ShowInfo } = $props();

	async function loadShowMetadata(show: ShowInfo) {

	}

	$effect(() => {
		loadShowMetadata(show);
	});
</script>

<div class="border border-neutral-500 p-4">
	{show.name}
</div>
