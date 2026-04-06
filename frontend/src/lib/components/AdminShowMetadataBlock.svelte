<script lang="ts">
	import type { ShowInfo } from '$lib/types/export_types';

	let { show }: { show: ShowInfo } = $props();
</script>

<div>
	{show.name}
</div>
