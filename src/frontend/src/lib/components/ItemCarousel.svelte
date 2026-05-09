<script lang="ts">
	import ChevronLeftIcon from '$lib/icons/ChevronLeftIcon.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let scroller: HTMLDivElement;

	let canScrollLeft = $state(false);
	let canScrollRight = $state(false);
	let canScroll = $state(false);
	let isHovered = $state(false);

	function updateScrollButtons() {
		if (!scroller) return;

		const maxScrollLeft = scroller.scrollWidth - scroller.clientWidth;

		canScroll = maxScrollLeft >= 1;
		canScrollLeft = scroller.scrollLeft > 1;
		canScrollRight = scroller.scrollLeft < maxScrollLeft;
	}

	function scroll(direction: 'left' | 'right') {
		const amount = scroller.clientWidth * 0.8;

		scroller.scrollBy({
			left: direction === 'right' ? amount : -amount,
			behavior: 'smooth'
		});
	}

	$effect(() => {
		if (!scroller) return;

		updateScrollButtons();

		const resizeObserver = new ResizeObserver(() => {
			updateScrollButtons();
		});

		resizeObserver.observe(scroller);

		return () => {
			resizeObserver.disconnect();
		};
	});
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<section
	class="relative space-y-3"
	onpointerenter={() => (isHovered = true)}
	onpointerleave={() => (isHovered = false)}
>
	<div
		bind:this={scroller}
		onscroll={updateScrollButtons}
		class="flex gap-4 overflow-x-auto overflow-y-hidden scroll-smooth pb-2 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
	>
		{@render children()}
	</div>

	{#if canScroll && canScrollLeft}
		<button
			onclick={() => scroll('left')}
			class="absolute top-0 bottom-0 left-0 cursor-pointer bg-linear-to-l via-black/55 to-black/70 px-4 text-neutral-300 transition-all duration-300 hover:via-black/80 hover:to-black/90 hover:text-neutral-50"
			class:opacity-0={!isHovered}
			class:pointer-events-none={!isHovered}
		>
			<ChevronLeftIcon />
		</button>
	{/if}

	{#if canScroll && canScrollRight}
		<button
			onclick={() => scroll('right')}
			class="absolute top-0 right-0 bottom-0 cursor-pointer bg-linear-to-r via-black/55 to-black/70 px-4 text-neutral-300 transition-all duration-300 hover:via-black/80 hover:to-black/90 hover:text-neutral-50"
			class:opacity-0={!isHovered}
			class:pointer-events-none={!isHovered}
		>
			<ChevronRightIcon />
		</button>
	{/if}
</section>
