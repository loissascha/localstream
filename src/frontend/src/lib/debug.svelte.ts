type DebugState = {
	showDebug: boolean;
};

export const debug = $state<DebugState>({
	showDebug: false
});

export function showDebug(toggle: boolean) {
	debug.showDebug = toggle;
}
