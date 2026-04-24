<script lang="ts">
	import { Popover } from 'melt/builders';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		items: Snippet;
	}
	let { children, items }: Props = $props();

	const popover = new Popover();
</script>

<button {...popover.trigger} class="cursor-pointer">{@render children()}</button>
<div {...popover.content} class="rounded-md bg-neutral-800 text-white shadow-lg">
	<div {...popover.arrow}></div>
	{@render items()}
</div>
