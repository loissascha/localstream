<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Home' }: Props = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="1.8"
	stroke-linecap="round"
	stroke-linejoin="round"
	width={size}
	height={size}
	class={className}
	role="img"
	aria-hidden={title ? undefined : 'true'}
>
	{#if title}
		<title>{title}</title>
	{/if}
	<path d="M3.75 9.75 12 3l8.25 6.75" />
	<path d="M6.75 8.85V20.25h10.5V8.85" />
	<path d="M10.5 20.25v-5.25h3v5.25" />
</svg>
