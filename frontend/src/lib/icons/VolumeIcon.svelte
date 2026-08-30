<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Volume' }: Props = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	viewBox="0 0 24 24"
	fill="currentColor"
	width={size}
	height={size}
	class={className}
	role="img"
	aria-hidden={title ? undefined : 'true'}
>
	{#if title}
		<title>{title}</title>
	{/if}
	<path
		d="M10.6 5.47A1.5 1.5 0 0 1 13 6.67v10.66a1.5 1.5 0 0 1-2.4 1.2l-3.5-2.66H5a2 2 0 0 1-2-2v-3.74a2 2 0 0 1 2-2h2.1l3.5-2.66Z"
	/>
	<path
		d="M15.75 8.5a.75.75 0 0 1 1.06 0A5 5 0 0 1 18.3 12a5 5 0 0 1-1.49 3.5.75.75 0 1 1-1.06-1.06A3.5 3.5 0 0 0 16.8 12a3.5 3.5 0 0 0-1.05-2.44.75.75 0 0 1 0-1.06Z"
	/>
	<path
		d="M17.9 6.35a.75.75 0 0 1 1.06 0A8 8 0 0 1 21.3 12a8 8 0 0 1-2.34 5.65.75.75 0 1 1-1.06-1.06A6.5 6.5 0 0 0 19.8 12a6.5 6.5 0 0 0-1.9-4.59.75.75 0 0 1 0-1.06Z"
	/>
</svg>
