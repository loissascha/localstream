import type { AuthUserIsAdminResponse } from '$lib/types/export_types';

export async function checkIfUserIsAdmin(bearerToken: string): Promise<AuthUserIsAdminResponse> {
	const response = await fetch('/api/auth/user/admin', {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
    const result = (await response.json()) as AuthUserIsAdminResponse;
    return result;
}
