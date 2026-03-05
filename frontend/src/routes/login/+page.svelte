<script lang="ts">
	let { form } = $props();
</script>

<main
	class="grid min-h-dvh place-items-center bg-[radial-gradient(circle_at_20%_10%,#d7e8ff_0%,transparent_40%),radial-gradient(circle_at_80%_0%,#c8f5e9_0%,transparent_32%),linear-gradient(180deg,#eef2f7_0%,#dce6f2_100%)] px-4 py-6"
>
	<section
		class="w-full max-w-sm rounded-2xl border border-slate-900/10 bg-white/85 p-6 shadow-lg shadow-slate-900/10"
	>
		<h1 class="m-0 text-2xl font-semibold text-slate-900">Sign in</h1>
		<p class="mt-2 text-sm text-slate-600">Enter any username to create a session.</p>

		<form method="POST" class="mt-5 grid gap-3">
			<label class="grid gap-1.5 text-sm text-slate-700" for="name">
				Username
				<input
					id="name"
					name="name"
					type="text"
					required
					autocomplete="username"
					class="rounded-md border border-slate-300 bg-white px-3 py-2 text-slate-900"
				/>
			</label>

			{#if form?.error}
				<p class="text-sm text-red-700">{form.error}</p>
			{/if}

			<button
				type="submit"
				class="cursor-pointer rounded-md bg-slate-900 px-3 py-2 text-sm font-medium text-white"
			>
				Continue
			</button>
		</form>
	</section>
</main>
