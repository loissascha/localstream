<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Notifications' }: Props = $props();
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
		d="M12 3.5A5.5 5.5 0 0 0 6.5 9v3.34c0 .73-.25 1.43-.7 1.99l-1.14 1.42A1.25 1.25 0 0 0 5.64 18h12.72a1.25 1.25 0 0 0 .98-2.25l-1.14-1.42a3.2 3.2 0 0 1-.7-2V9A5.5 5.5 0 0 0 12 3.5Z"
	/>
	<path d="M9.75 19a2.25 2.25 0 0 0 4.5 0h-4.5Z" />
</svg>
