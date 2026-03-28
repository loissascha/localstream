<script lang="ts">
	type FormData = {
		name: string;
		type: string;
		path: string;
	};

	let form = $state<FormData>({
		name: '',
		type: 'none',
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
>
	<div>
		<label for="name">Name:</label>
		<input type="text" id="name" name="name" bind:value={form.name} />
	</div>
	<div>
		<label for="type">Type:</label>
		<select id="type" name="type" bind:value={form.type}>
			<option value="none"> Please select </option>
			<option value="shows"> TV Shows </option>
			<option value="movies"> Movies </option>
		</select>
	</div>
	<div>
		<label for="path">Path:</label>
		<input type="text" id="path" name="path" bind:value={form.path} />
	</div>
	<div>
		<button>Submit</button>
	</div>
</form>
