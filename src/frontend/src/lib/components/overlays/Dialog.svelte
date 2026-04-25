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
	class="m-auto max-h-[90%] w-[90%] max-w-[90%] rounded-xl border border-neutral-500 bg-neutral-800 p-4 text-white lg:w-3xl lg:max-w-4xl"
>
	<button
		class="absolute top-5 right-5 cursor-pointer text-xl"
		onclick={() => {
			close();
		}}>&times;</button
	>
	{@render content(close)}
</dialog>
