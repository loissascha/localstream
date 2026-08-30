<script lang="ts">
	import { debug } from '$lib/debug.svelte';

	interface Props {
		percentage: number;
	}
	let { percentage }: Props = $props();
</script>

<div class="bg-neutral-600">
	<div style={`width: ${percentage}%;`} class={`h-2 bg-brand text-sm`}></div>
	{#if debug.showDebug}
		<span class="text-shadow text-white text-shadow-black text-shadow-lg"
			>{percentage.toFixed(2)}%</span
		>
	{/if}
</div>
