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
			class="rounded bg-neutral-800 px-2 py-1 focus:ring focus:ring-green-400 focus:outline-none"
			bind:value={form.name}
		/>
	</div>
	<div class="flex flex-col">
		<label for="type">Type:</label>
		<div class="flex gap-2">
			<div>
				<input
					type="radio"
					name="type"
					value={LibraryType.Shows}
					bind:group={form.type}
					id="type_shows"
				/>
				<label for="type_shows">Shows</label>
			</div>
			<div>
				<input
					type="radio"
					name="type"
					value={LibraryType.Movies}
					id="type_movies"
					bind:group={form.type}
				/>
				<label for="type_movies">Movies</label>
			</div>
		</div>
	</div>
	<div class="flex flex-col">
		<label for="path">Path:</label>
		<input
			type="text"
			id="path"
			name="path"
			class="rounded bg-neutral-800 px-2 py-1 focus:ring focus:ring-green-400 focus:outline-none"
			bind:value={form.path}
		/>
	</div>
	<div>
		<button class="cursor-pointer rounded bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
			>Submit</button
		>
	</div>
</form>
