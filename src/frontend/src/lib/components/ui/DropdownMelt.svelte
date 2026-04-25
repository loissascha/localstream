<script lang="ts">
	import { Popover } from 'melt/builders';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		items: Snippet;
		className?: string;
	}
	let { children, items, className = '' }: Props = $props();

	const popover = new Popover();
</script>

<button {...popover.trigger} class={`cursor-pointer ${className}`}>{@render children()}</button>
<div
	{...popover.content}
	class="rounded-md border border-neutral-500 bg-neutral-800 text-white shadow-lg"
>
	<div {...popover.arrow}></div>
	{@render items()}
</div>
