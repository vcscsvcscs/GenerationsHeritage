import { error_generic, error_not_found, error_unauthorized } from '$lib/paraglide/messages';

export const errorMessage = (status: number): string =>
	status === 401 ? error_unauthorized() : status === 404 ? error_not_found() : error_generic();
