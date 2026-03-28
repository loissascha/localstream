<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Close' }: Props = $props();
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
		d="M6.47 5.41a.75.75 0 0 1 1.06 0L12 9.88l4.47-4.47a.75.75 0 1 1 1.06 1.06L13.06 10.94l4.47 4.47a.75.75 0 0 1-1.06 1.06L12 12l-4.47 4.47a.75.75 0 1 1-1.06-1.06l4.47-4.47-4.47-4.47a.75.75 0 0 1 0-1.06Z"
	/>
</svg>
