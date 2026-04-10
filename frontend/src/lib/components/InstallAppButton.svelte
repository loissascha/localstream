<script lang="ts">
	import { onMount } from 'svelte';

	type BeforeInstallPromptEvent = Event & {
		prompt: () => Promise<void>;
		userChoice: Promise<{ outcome: 'accepted' | 'dismissed'; platform: string }>;
	};

	let deferredPrompt = $state<BeforeInstallPromptEvent | null>(null);
	let isStandalone = $state(false);

	const detectStandalone = () => {
		return (
			window.matchMedia('(display-mode: standalone)').matches ||
			(window.navigator as Navigator & { standalone?: boolean }).standalone === true
		);
	};

	const installApp = async () => {
		if (!deferredPrompt) return;

		await deferredPrompt.prompt();
		await deferredPrompt.userChoice;
		deferredPrompt = null;
		isStandalone = detectStandalone();
	};

	onMount(() => {
		isStandalone = detectStandalone();

		const onBeforeInstallPrompt = (event: Event) => {
			event.preventDefault();
			deferredPrompt = event as BeforeInstallPromptEvent;
		};

		const onAppInstalled = () => {
			deferredPrompt = null;
			isStandalone = true;
		};

		window.addEventListener('beforeinstallprompt', onBeforeInstallPrompt);
		window.addEventListener('appinstalled', onAppInstalled);

		return () => {
			window.removeEventListener('beforeinstallprompt', onBeforeInstallPrompt);
			window.removeEventListener('appinstalled', onAppInstalled);
		};
	});
</script>

{#if deferredPrompt}
	<button
		type="button"
		onclick={installApp}
		class="cursor-pointer rounded-md border border-green-500 bg-green-700 px-3 py-1.5 text-sm font-medium text-white hover:bg-green-600"
	>
		Install App
	</button>
{/if}
