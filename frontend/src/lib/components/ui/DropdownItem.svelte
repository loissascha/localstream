<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		href?: string;
		onclick?: () => void;
	}
	let { children, href = '', onclick }: Props = $props();
</script>

<a
	{href}
	{onclick}
	role="menuitem"
	class="block cursor-pointer px-4 py-2 text-left whitespace-nowrap transition-colors hover:bg-neutral-700 focus:bg-neutral-700 focus:outline-none"
>
	{@render children()}
</a>
