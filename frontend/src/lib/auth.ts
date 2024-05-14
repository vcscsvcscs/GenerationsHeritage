import { UserManager, WebStorageStateStore, type User } from 'oidc-client-ts';
import { isAuthenticated, user } from './stores';
import { PUBLIC_ZITADEL_CLIENT_ID, PUBLIC_ISSUER } from '$env/static/public';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

let userManager: UserManager;
if (browser) {
	const config = {
		authority: PUBLIC_ISSUER, // At Zitadel Project Console > [Your project] > [Your application] > URLs - Issuer
		client_id: PUBLIC_ZITADEL_CLIENT_ID, // At Zitadel Project Console > [Your project] > [Your application] > Configuration - Client ID
		redirect_uri: window.location.hostname + '/callback', // At Zitadel Project Console > [Your project] > [Your application] > URLs - Login Redirect URI
		response_type: 'code',
		scope: 'openid profile email',
		post_logout_redirect_uri: window.location.hostname,
		userStore: new WebStorageStateStore({ store: window.localStorage }),
		automaticSilentRenew: true,
		silent_redirect_uri: window.location.hostname + '/silent-refresh'
	};

	userManager = new UserManager(config);

	userManager.events.addUserLoaded((loadedUser: User) => {
		console.log('userManager.events.addUserLoaded');
		user.set(loadedUser);
		isAuthenticated.set(true);
	});

	userManager.events.addUserUnloaded(() => {
		console.log('userManager.events.addUserUnloaded');
		user.set(null);
		isAuthenticated.set(false);
	});
}

async function login(): Promise<void> {
	console.log('UserManager.login()');
	if (browser) {
		await userManager.signinRedirect();
	}
}

async function logout(): Promise<void> {
	if (browser) {
		await userManager.signoutRedirect();
	}
}

async function handleCallback(): Promise<void> {
	if (browser) {
		await userManager.signinRedirectCallback();
		goto('/');
	}
}

async function handleSilentCallback(): Promise<void> {
	if (browser) {
		await userManager.signinSilentCallback();
		goto('/');
	}
}

export { login, logout, handleCallback, handleSilentCallback };
