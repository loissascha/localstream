<script lang="ts">
	import Overlay from './Overlay.svelte';

	interface Props {
		close: () => void;
	}
	let { close }: Props = $props();
</script>

<Overlay {close}>
	<h1 class="text-2xl font-bold tracking-wide">Create Collection</h1>
	<form class="mt-8 mb-4">
		<label for="name">Collection Name</label>
		<input
			id="name"
			type="text"
			class="w-full rounded bg-neutral-700 px-4 py-2"
			placeholder="Collection Name"
		/>
		<button
			class="mt-4 w-full cursor-pointer rounded-full bg-brand/80 px-4 py-2 font-semibold hover:bg-brand"
			>&plus; Create</button
		>
	</form>
</Overlay>
