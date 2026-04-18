<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		close: () => void;
	}
	let { children, close }: Props = $props();
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
		onclick={(e) => {
			e.stopPropagation();
		}}
		class="fixed inset-5 m-auto max-h-160 max-w-3xl overflow-y-auto rounded bg-neutral-800 p-4 md:inset-10"
	>
		<div class="absolute"></div>
		{@render children()}
	</div>
</div>
