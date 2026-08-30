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
		<div
			class="pointer-events-none absolute top-0 bottom-0 left-0 flex items-center justify-center px-4"
		>
			<button
				onclick={() => scroll('left')}
				class="pointer-events-auto cursor-pointer rounded-full bg-neutral-800 p-4 text-brand shadow transition-all duration-200 hover:scale-110"
				class:opacity-0={!isHovered}
				class:pointer-events-none={!isHovered}
			>
				<ChevronLeftIcon size={26} />
			</button>
		</div>
	{/if}

	{#if canScroll && canScrollRight}
		<div
			class="pointer-events-none absolute top-0 right-0 bottom-0 flex items-center justify-center px-4"
		>
			<button
				onclick={() => scroll('right')}
				class="pointer-events-auto cursor-pointer rounded-full bg-neutral-800 p-4 text-brand shadow transition-all duration-200 hover:scale-110"
				class:opacity-0={!isHovered}
				class:pointer-events-none={!isHovered}
			>
				<ChevronRightIcon size={26} />
			</button>
		</div>
	{/if}
</section>
