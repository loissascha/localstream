<script lang="ts">
	import { updateCollection } from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import type { CollectionInfo } from '$lib/types/export_types';
	import Overlay from './Overlay.svelte';

	interface Props {
		close: () => void;
		collection: CollectionInfo;
	}
	let { close, collection }: Props = $props();

	$effect(() => {
		collection;
		name = collection.name;
	});

	let name = $state('');
	let error_message = $state('');

	async function submitForm() {
		try {
			if (!auth.token) return;
			if (name.trim() == '') {
				error_message = 'Name must not be empty!';
				return;
			}

			await updateCollection(auth.token, collection.id, {
				name: name
			});

			collection.name = name;
			close();
		} catch (e) {
			const m = (e as Error).message;
			error_message = m;
		}
	}
</script>

<Overlay {close}>
	<h1 class="text-2xl font-bold tracking-wide">Rename Collection</h1>
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
			>Rename</button
		>
	</form>
</Overlay>
