<script lang="ts">
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
		console.log('submit form data: ', form);
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
			class="rounded bg-neutral-800 px-2 py-1 focus:outline-none focus:ring focus:ring-green-400"
			bind:value={form.name}
		/>
	</div>
	<div class="flex flex-col">
		<label for="type">Type:</label>
		<select
			id="type"
			name="type"
			class="rounded bg-neutral-800 px-2 py-1 focus:outline-none focus:ring focus:ring-green-400"
			bind:value={form.type}
		>
			<option value="shows"> TV Shows </option>
			<option value="movies"> Movies </option>
		</select>
	</div>
	<div class="flex flex-col">
		<label for="path">Path:</label>
		<input
			type="text"
			id="path"
			name="path"
			class="rounded bg-neutral-800 px-2 py-1 focus:outline-none focus:ring focus:ring-green-400"
			bind:value={form.path}
		/>
	</div>
	<div>
		<button class="rounded bg-neutral-800 px-4 py-2 hover:bg-neutral-700">Submit</button>
	</div>
</form>
