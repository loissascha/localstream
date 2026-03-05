import { dev } from '$app/environment';
import { fail, redirect } from '@sveltejs/kit';

export const load = async ({ locals, url }) => {
	if (locals.user) {
		throw redirect(302, url.searchParams.get('next') || '/');
	}
};

export const actions = {
	default: async ({ cookies, request, url }) => {
		const data = await request.formData();
		const name = data.get('name');

		if (typeof name !== 'string' || name.trim().length === 0) {
			return fail(400, {
				error: 'Please enter a username.'
			});
		}

		cookies.set('session', name.trim(), {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			secure: !dev,
			maxAge: 60 * 60 * 24 * 7
		});

		throw redirect(302, url.searchParams.get('next') || '/');
	}
};
