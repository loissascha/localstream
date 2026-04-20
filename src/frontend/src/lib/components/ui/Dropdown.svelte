<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		items: Snippet;
	}
	let { children, items }: Props = $props();

	let open = $state(true);
</script>

<div class="relative">
	<button
		onclick={() => {
			open = !open;
		}}
		class="cursor-pointer"
	>
		{@render children()}
	</button>
	{#if open}
		<div
			class="absolute top-full right-0 left-0 z-40 mt-2 flex w-20 flex-col gap-1 bg-neutral-600"
		>
			{@render items()}
		</div>
	{/if}
</div>
