import { getCookie } from './cookies';

type AuthState = {
    loggedIn: boolean;
    token: string | null;
};

export const auth = $state<AuthState>({
    loggedIn: false,
    token: null
});

export function loadAuthFromCookies() {
    const token = getCookie('bearer');

    auth.token = token;
    auth.loggedIn = !!token;
}

export function clearAuth() {
    auth.token = null;
    auth.loggedIn = false;
}
