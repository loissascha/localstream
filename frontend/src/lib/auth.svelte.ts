import { checkIfUserIsAdmin } from './api/admin/auth';
import { getCookie } from './cookies';

type AuthState = {
	initialized: boolean;
	loggedIn: boolean;
	token: string | null;
	isAdmin: boolean;
};

export const auth = $state<AuthState>({
	initialized: false,
	loggedIn: false,
	token: null,
	isAdmin: false
});

export async function loadAuthFromCookies() {
	const token = getCookie('bearer');

	auth.isAdmin = false;
	if (token != null) {
		var response = await checkIfUserIsAdmin(token);
		auth.isAdmin = response.is_admin;
	}

	auth.token = token;
	auth.loggedIn = !!token;
	auth.initialized = true;
}

export function clearAuth() {
	auth.token = null;
	auth.loggedIn = false;
}
