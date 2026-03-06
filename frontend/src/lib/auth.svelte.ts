import { getCookie } from './cookies';

type AuthState = {
    initialized: boolean;
    loggedIn: boolean;
    token: string | null;
};

export const auth = $state<AuthState>({
    initialized: false,
    loggedIn: false,
    token: null
});

export function loadAuthFromCookies() {
    const token = getCookie('bearer');

    auth.token = token;
    auth.loggedIn = !!token;
    auth.initialized = true;
}

export function clearAuth() {
    auth.token = null;
    auth.loggedIn = false;
}
