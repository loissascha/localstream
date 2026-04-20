<script lang="ts">
	import type { Snippet } from 'svelte';
	import { scale } from 'svelte/transition';

	interface Props {
		children: Snippet;
		items: Snippet;
	}
	let { children, items }: Props = $props();

	let open = $state(false);
	let container: HTMLDivElement | undefined = $state();

	function handleDocumentClick(event: MouseEvent) {
		if (!open || !container) return;

		const target = event.target;
		if (target instanceof Node && !container.contains(target)) {
			open = false;
		}
	}

	function handleDocumentKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			open = false;
		}
	}
</script>

<svelte:document onclick={handleDocumentClick} onkeydown={handleDocumentKeydown} />

<div bind:this={container} class="relative inline-block">
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
			transition:scale={{ duration: 140, start: 0.95 }}
			class="absolute top-full left-1/2 z-40 mt-2 flex w-max min-w-full -translate-x-1/2 flex-col gap-1 overflow-hidden rounded-md bg-neutral-800 shadow-lg"
		>
			{@render items()}
		</div>
	{/if}
</div>
