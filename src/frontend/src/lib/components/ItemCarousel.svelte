<script lang="ts">
	import ChevronLeftIcon from '$lib/icons/ChevronLeftIcon.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let scroller: HTMLDivElement;

	function scroll(direction: 'left' | 'right') {
		const amount = scroller.clientWidth * 0.8;

		scroller.scrollBy({
			left: direction === 'right' ? amount : -amount,
			behavior: 'smooth'
		});
	}
</script>

<section class="relative space-y-3">
	<div
		bind:this={scroller}
		class="flex gap-4 overflow-x-auto scroll-smooth pb-2 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
	>
		{@render children()}
	</div>

	<button
		onclick={() => scroll('left')}
		class="absolute top-0 bottom-0 left-0 cursor-pointer bg-linear-to-l via-black/55 to-black/70 px-4 text-neutral-300 transition-all duration-300 hover:via-black/80 hover:to-black/90 hover:text-neutral-50"
	>
		<ChevronLeftIcon />
	</button>

	<button
		onclick={() => scroll('right')}
		class="absolute top-0 right-0 bottom-0 cursor-pointer bg-linear-to-r via-black/55 to-black/70 px-4 text-neutral-300 transition-all duration-300 hover:via-black/80 hover:to-black/90 hover:text-neutral-50"
	>
		<ChevronRightIcon />
	</button>
</section>
