<script lang="ts">
	import Overlay from './Overlay.svelte';

	interface Props {
		close: () => void;
	}
	let { close }: Props = $props();

	let name = $state('');
	let error_message = $state('');

	function submitForm() {
		if (name.trim() == '') {
			error_message = 'Name must not be empty!';
			return;
		}
		name = '';
	}
</script>

<Overlay {close}>
	<h1 class="text-2xl font-bold tracking-wide">Create Collection</h1>
	{#if error_message != ''}
		<div class="mt-4 text-red-500">
			{error_message}
		</div>
	{/if}
	<form
		class="mt-4 mb-4"
		onsubmit={(e) => {
			e.preventDefault();
			submitForm();
		}}
	>
		<label for="name">Collection Name</label>
		<input
			bind:value={name}
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
