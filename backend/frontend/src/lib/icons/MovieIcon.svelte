<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Movies' }: Props = $props();
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
		d="M6 3.75a2.25 2.25 0 0 0-2.25 2.25v12A2.25 2.25 0 0 0 6 20.25h12A2.25 2.25 0 0 0 20.25 18V6A2.25 2.25 0 0 0 18 3.75H6ZM8 5.75v2h2v-2h2v2h2v-2h2v2h2v8.5h-2v2h-2v-2h-2v2h-2v-2H8v2H6v-2H5.75v-2H8v-2H5.75v-2H8v-2H5.75v-2H8v-2h2Z"
	/>
	<path d="M10.25 10.25c0-1.03 1.12-1.65 2-1.12l2.75 1.62a1.3 1.3 0 0 1 0 2.24l-2.75 1.62a1.3 1.3 0 0 1-2-1.12v-3.24Z" />
</svg>
