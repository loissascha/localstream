<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		close: () => void;
	}
	let { children, close }: Props = $props();

	function stopOverlayClick(event: MouseEvent) {
		event.stopPropagation();
	}
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	onclick={() => {
		close();
	}}
	class="fixed inset-0 z-50 bg-black/40"
>
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={stopOverlayClick}
		class="fixed inset-5 m-auto flex max-h-160 max-w-3xl flex-col rounded bg-neutral-800 md:inset-10"
	>
		<button
			type="button"
			onclick={() => {
				close();
			}}
			class="absolute top-3 right-3 z-10 flex h-9 w-9 items-center justify-center rounded-full bg-neutral-700 text-xl leading-none text-neutral-100 transition-colors hover:bg-neutral-600"
			aria-label="Close overlay"
		>
			&times;
		</button>

		<div class="overflow-y-auto p-4">
			{@render children()}
		</div>
	</div>
</div>
