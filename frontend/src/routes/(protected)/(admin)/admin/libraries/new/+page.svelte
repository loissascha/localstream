<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { createLibrary } from '$lib/api/admin/libraries';
	import { auth } from '$lib/auth.svelte';
	import { LibraryType } from '$lib/types/enums';

	type FormData = {
		name: string;
		type: LibraryType;
		path: string;
	};

	let form = $state<FormData>({
		name: '',
		type: LibraryType.Shows,
		path: ''
	});

	async function handleSubmit() {
		try {
			if (!auth.token) {
				throw new Error('Auth token not set');
			}
			await createLibrary(auth.token, {
				name: form.name,
				type: form.type,
				path: form.path
			});
			goto(resolve('/(protected)/(admin)/admin/libraries'));
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}
</script>

<h1 class="mb-4 text-2xl tracking-wider">New Library</h1>

<form
	onsubmit={(e) => {
		e.preventDefault();
		handleSubmit();
	}}
	class="flex flex-col gap-2"
>
	<div class="flex flex-col">
		<label for="name">Name:</label>
		<input
			type="text"
			id="name"
			name="name"
			class="rounded bg-neutral-800 px-2 py-1 focus:ring focus:ring-brand focus:outline-none"
			bind:value={form.name}
		/>
	</div>
	<div class="flex flex-col">
		<label for="type">Type:</label>
		<div class="flex gap-2">
			<label
				for="type_shows"
				class={`cursor-pointer rounded-lg px-4 py-2 select-none ${form.type == LibraryType.Shows ? 'bg-brand/20' : 'bg-neutral-800'}`}
			>
				<input
					type="radio"
					name="type"
					value={LibraryType.Shows}
					bind:group={form.type}
					id="type_shows"
				/>
				Shows
			</label>
			<label
				for="type_movies"
				class={`cursor-pointer rounded-lg px-4 py-2 select-none ${form.type == LibraryType.Movies ? 'bg-brand/20' : 'bg-neutral-800'}`}
			>
				<input
					type="radio"
					name="type"
					value={LibraryType.Movies}
					id="type_movies"
					bind:group={form.type}
				/>
				Movies
			</label>
		</div>
	</div>
	<div class="flex flex-col">
		<label for="path">Path:</label>
		<input
			type="text"
			id="path"
			name="path"
			class="rounded bg-neutral-800 px-2 py-1 focus:ring focus:ring-brand focus:outline-none"
			bind:value={form.path}
		/>
	</div>
	<div>
		<button class="cursor-pointer rounded bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
			>Submit</button
		>
	</div>
</form>
