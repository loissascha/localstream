<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Mute' }: Props = $props();
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
		d="M5.47 4.41a.75.75 0 0 1 1.06 0L12 9.88l5.47-5.47a.75.75 0 1 1 1.06 1.06L13.06 10.94l5.47 5.47a.75.75 0 0 1-1.06 1.06L12 12l-5.47 5.47a.75.75 0 0 1-1.06-1.06l5.47-5.47-5.47-5.47a.75.75 0 0 1 0-1.06Z"
	/>
</svg>
