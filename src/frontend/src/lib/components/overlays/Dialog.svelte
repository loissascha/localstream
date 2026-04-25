<script lang="ts">
	import { Dialog } from 'melt/builders';
	import type { Snippet } from 'svelte';

	const dialog = new Dialog();

	interface Props {
		children: Snippet;
		content: Snippet<[() => void]>;
	}
	let { children, content }: Props = $props();

	function close() {
		dialog.open = false;
	}
</script>

<button {...dialog.trigger}>{@render children()}</button>

<div {...dialog.overlay}></div>

<dialog
	{...dialog.content}
	class="m-auto max-h-[90%] max-w-[90%] w-[90%] lg:w-3xl rounded-xl border border-neutral-500 bg-neutral-800 p-4 text-white lg:max-w-4xl"
>
	{@render content(close)}
</dialog>
